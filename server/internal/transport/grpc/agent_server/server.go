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

	messagesByAgentID map[uuid.UUID]*messages
	mu                sync.RWMutex
	wg                sync.WaitGroup
}

type messages struct {
	toSend          chan *pb.ServerMessage
	pendingCommands map[string]chan *pb.CommandResult
	mu              sync.Mutex
}

func NewAgentService(interactionService *agent.AgentInteractionService) *AgentService {
	return &AgentService{
		agentInteractionService: interactionService,
		messagesByAgentID:       make(map[uuid.UUID]*messages, 10),
		wg:                      sync.WaitGroup{},
	}
}

func (s *AgentService) Register(server *grpc.Server) {
	pb.RegisterAgentServiceServer(server, s)
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

	messages := &messages{
		toSend:          make(chan *pb.ServerMessage, 100),
		pendingCommands: make(map[string]chan *pb.CommandResult),
	}

	s.mu.Lock()
	s.messagesByAgentID[idFromTLS] = messages
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.messagesByAgentID, idFromTLS)
		for _, c := range messages.pendingCommands {
			close(c)
		}
		s.mu.Unlock()
	}()

	s.agentInteractionService.HandleConnection(idFromTLS)

	ctx, cancel := context.WithCancel(stream.Context())
	s.runReader(
		stream,
		idFromTLS,
		func() {
			cancel()
		},
	)
	s.runWriter(ctx, messages.toSend, stream)

	s.wg.Wait()

	s.agentInteractionService.HandleDisconnection(idFromTLS)

	return nil
}

func (s *AgentService) runReader(stream pb.AgentService_ConnectServer, agentID uuid.UUID, onClose func()) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			msg, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					log.Println("grpc client disconnected")
					return
				}
				log.Println("error receiving message from grpc client:", err)
				onClose()
				return
			}

			switch msg.Payload.(type) {
			case *pb.AgentMessage_Metrics:
				metrics := msg.GetMetrics()
				log.Println("Metrics received:", metrics)
				// TODO: check context
				s.agentInteractionService.HandleMetricsBatch(context.Background(),
					convertMetricsBatchWithAgentIDFromProto(metrics, agentID),
				)

			case *pb.AgentMessage_CommandResult:
				commandResult := msg.GetCommandResult()
				s.handleCommandResult(agentID, commandResult)

			default:
				log.Println("unknown message received:", msg)
			}
		}
	}()
}

func (s *AgentService) runWriter(ctx context.Context, toSend chan *pb.ServerMessage, stream pb.AgentService_ConnectServer) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
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

func (s *AgentService) sendCommand(ctx context.Context, agentID uuid.UUID, command *pb.Command) (*pb.CommandResult, error) {
	s.mu.RLock()
	commands := s.messagesByAgentID[agentID]
	s.mu.RUnlock()

	// add chan to pending map to receive result
	commandID := uuid.NewString()
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
	agentCommands := s.messagesByAgentID[agentID]
	s.mu.Unlock()

	agentCommands.mu.Lock()
	resChan, ok := agentCommands.pendingCommands[result.CommandUuid]
	if !ok {
		log.Printf("CommandResult commandUUID not found in pending commands: %#v", result)
		return
	}

	resChan <- result
}

func (s *AgentService) RequestSpecifications() {
	command := &pb.Command{
		Payload: &pb.Command_SpecificationsRequest{
			SpecificationsRequest: &pb.SpecificationsRequest{},
		},
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pendingCommands[uuid.New().String()] = command
	s.toSend <- &pb.ServerMessage{
		Payload: &pb.ServerMessage_Command{
			Command: command,
		},
	}
}
