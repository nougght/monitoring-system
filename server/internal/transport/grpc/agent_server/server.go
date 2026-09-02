package agent_server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	agent_model "github.com/nougght/monitoring-system/server/internal/model/agent"
	agent "github.com/nougght/monitoring-system/server/internal/service/agent_interaction"
	"github.com/nougght/monitoring-system/server/internal/util"
	pb "github.com/nougght/monitoring-system/shared/go/proto/gen/agent/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AgentService struct {
	pb.UnimplementedAgentServiceServer

	agentInteractionService *agent.AgentInteractionService

	connByAgentID map[uuid.UUID]*conn
	mu            sync.RWMutex
}

type conn struct {
	toSend          chan *pb.ServerMessage
	pendingCommands map[string]chan *pb.CommandResult
	mu              sync.Mutex
}

func NewAgentService(interactionService *agent.AgentInteractionService) *AgentService {
	return &AgentService{
		agentInteractionService: interactionService,
		connByAgentID:           make(map[uuid.UUID]*conn, 10),
	}
}

func (s *AgentService) Register(server *grpc.Server) {
	pb.RegisterAgentServiceServer(server, s)
}

func (s *AgentService) BroadcastShutdown() {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.connByAgentID {
		select {
		case c.toSend <- &pb.ServerMessage{
			Payload: &pb.ServerMessage_Status{
				Status: &pb.HandshakeStatus{Connected: false},
			},
		}:
		default:
		}
	}
}

func (s *AgentService) StartStreamMJPEG(stream pb.AgentService_StartStreamMJPEGServer) error {
	log.Println("grpc streaming client connected")
	// TODO: sync auth with main grpc service
	idFromTLS, err := util.GetAgentIDFromContext(stream.Context())
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "invalid agent ID: %s", err.Error())
	}
	if idFromTLS == uuid.Nil {
		log.Println("parsin id from context")
		idFromTLS, err = ParseIDFromContext(stream.Context())
		if err != nil {
			return status.Errorf(codes.Unauthenticated, "invalid agent ID: %s", err.Error())
		}
	}
	if _, ok := s.connByAgentID[idFromTLS]; !ok {
		log.Println("client with agentID is not connected to main stream")
		return status.Errorf(codes.Unauthenticated, "agent is not connected")
	}
	log.Println("run streaming reader")
	log.Printf("id: %s", idFromTLS.String())
	wg := sync.WaitGroup{}
	err = s.runStreamingReader(stream, &wg, idFromTLS)
	wg.Wait()

	log.Println("streaming grpc finished")
	return err
}

func (s *AgentService) runStreamingReader(stream pb.AgentService_StartStreamMJPEGServer, wg *sync.WaitGroup, agentID uuid.UUID) error {

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()
		for {
			msg, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					log.Println("grpc streaming client disconnected")
					return
				}
				log.Println("error receiving message from streaming grpc client:", err)
				// onClose()
				return
			}

			switch msg.Payload.(type) {
			case *pb.StreamingAgentMessage_Frame:
				s.agentInteractionService.HandleFrame(msg.GetFrame().Data, agentID)

			case *pb.StreamingAgentMessage_Info:

			default:
				log.Println("unknown streaming message received:", msg)
			}
		}
	}()
	return nil
}

func (s *AgentService) Connect(stream pb.AgentService_ConnectServer) error {
	log.Println("grpc client connected")

	handshakeChan := make(chan *pb.AgentMessage, 1)
	go func() {
		msg, err := stream.Recv()
		if err != nil {
			log.Println("error receiving handshake message:", err)
			return
		}
		handshakeChan <- msg
	}()

	t := time.NewTimer(time.Minute * 5)

	var handshakeMsg *pb.Handshake
	select {
	case <-t.C:
		return status.Error(codes.DeadlineExceeded, "handshake timeout")
	case msg := <-handshakeChan:
		if msg.GetHandshake() == nil {
			return status.Error(codes.InvalidArgument, "expected handshake message")
		}
		handshakeMsg = msg.GetHandshake()
	}

	idFromTLS, err := util.GetAgentIDFromContext(stream.Context())
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "invalid agent ID: %s", err.Error())
	}

	if idFromTLS.String() != handshakeMsg.GetAgentUuid() {
		return status.Error(codes.Unauthenticated, "agent ID mismatch")
	}

	conn := &conn{
		toSend:          make(chan *pb.ServerMessage, 100),
		pendingCommands: make(map[string]chan *pb.CommandResult),
	}

	s.mu.Lock()
	s.connByAgentID[idFromTLS] = conn
	s.mu.Unlock()

	err = stream.Send(
		&pb.ServerMessage{
			Payload: &pb.ServerMessage_Status{
				Status: &pb.HandshakeStatus{
					Connected: true,
				},
			},
		},
	)

	defer func() {
		s.mu.Lock()
		delete(s.connByAgentID, idFromTLS)
		for _, c := range conn.pendingCommands {
			close(c)
		}
		s.mu.Unlock()
	}()

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(stream.Context())
	s.runReader(
		stream,
		&wg,
		idFromTLS,
		cancel,
	)
	s.runWriter(ctx, &wg, idFromTLS, conn.toSend, stream)

	s.agentInteractionService.HandleConnection(idFromTLS)

	wg.Wait()
	log.Println("main grpc finished")
	s.agentInteractionService.HandleDisconnection(idFromTLS)

	return nil
}

func (s *AgentService) runReader(stream pb.AgentService_ConnectServer, wg *sync.WaitGroup, agentID uuid.UUID, onClose func()) {
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()
		for {
			msg, err := stream.Recv()
			if err != nil {
				onClose()
				if errors.Is(err, io.EOF) {
					log.Println("grpc client disconnected")
					return
				}
				log.Println("error receiving message from grpc client:", err)
				return
			}

			switch msg.Payload.(type) {
			case *pb.AgentMessage_Metrics:
				metrics := msg.GetMetrics()
				log.Printf("Metrics received: %#v", metrics)
				// TODO: check context
				err = s.agentInteractionService.HandleMetricsBatch(context.Background(),
					convertMetricsBatchWithAgentIDFromProto(metrics, agentID),
				)
				if err != nil {
					// TODO: add response to agent
					log.Println(err.Error())
				}

			case *pb.AgentMessage_CommandResult:
				commandResult := msg.GetCommandResult()
				s.handleCommandResult(agentID, commandResult)

			default:
				log.Println("unknown message received:", msg)
			}
		}
	}()
}

func (s *AgentService) runWriter(ctx context.Context, wg *sync.WaitGroup, agentID uuid.UUID, toSend chan *pb.ServerMessage, stream pb.AgentService_ConnectServer) {
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()
		for {
			select {
			case msg := <-toSend:
				err := stream.Send(msg)
				if err != nil {
					log.Println("error sending message to grpc client:", err)
					return
				}
			case <-ctx.Done():
				log.Println("writer context closed")
				return
			}

		}
	}()
}

func (s *AgentService) sendCommandWithTimeout(ctx context.Context, agentID uuid.UUID, command *pb.Command) (*pb.CommandResult, error) {
	commandCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	return s.sendCommand(commandCtx, agentID, command)
}

func (s *AgentService) sendCommand(ctx context.Context, agentID uuid.UUID, command *pb.Command) (*pb.CommandResult, error) {
	s.mu.RLock()
	commands := s.connByAgentID[agentID]
	s.mu.RUnlock()

	// add chan to pending map to receive result
	commandID := uuid.NewString()
	command.CommandUuid = commandID
	resultChan := make(chan *pb.CommandResult, 1)
	commands.mu.Lock()
	commands.pendingCommands[commandID] = resultChan
	select {
	case commands.toSend <- &pb.ServerMessage{
		Payload: &pb.ServerMessage_Command{
			Command: command,
		},
	}:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	commands.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(commands.pendingCommands, commandID)
		s.mu.Unlock()
	}()

	select {
	case result, ok := <-resultChan:
		if !ok {
			return nil, fmt.Errorf("result chan closed")
		}
		return result, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *AgentService) handleCommandResult(agentID uuid.UUID, result *pb.CommandResult) {
	s.mu.Lock()
	agentCommands := s.connByAgentID[agentID]
	s.mu.Unlock()

	agentCommands.mu.Lock()
	defer agentCommands.mu.Unlock()
	resChan, ok := agentCommands.pendingCommands[result.CommandUuid]
	if !ok {
		log.Printf("CommandResult commandUUID not found in pending commands: %#v", result)
		return
	}

	resChan <- result
}

func (s *AgentService) RequestSpecifications(ctx context.Context, agentID uuid.UUID) (*agent_model.Specs, error) {
	res, err := s.sendCommandWithTimeout(ctx, agentID, &pb.Command{
		Payload: &pb.Command_SpecificationsRequest{
			SpecificationsRequest: &pb.SpecificationsRequest{},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send specifications command: %w", err)
	}

	specs, ok := res.Payload.(*pb.CommandResult_SpecificationsResponse)
	if !ok {
		return nil, fmt.Errorf("incorrect command result type")
	}

	return convertSpecsFromProto(specs.SpecificationsResponse.Specs), nil
}
