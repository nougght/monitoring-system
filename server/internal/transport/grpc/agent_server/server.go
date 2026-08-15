package agent_server

import (
	"context"
	"errors"
	"io"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	agentregistry "github.com/nougght/monitoring-system/server/internal/service/agent_registry"
	"github.com/nougght/monitoring-system/server/internal/util"
	pb "github.com/nougght/monitoring-system/shared/go/proto/gen/agent/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AgentService struct {
	pb.UnimplementedAgentServiceServer

	registryService *agentregistry.AgentRegistryService
	// metricsService metrics.MetricsService
	toSend          chan *pb.ServerMessage
	commandResults  chan *pb.CommandResult
	pendingCommands map[string]*pb.Command
	mu              sync.Mutex
	wg              sync.WaitGroup
}

func NewAgentService(registryService *agentregistry.AgentRegistryService) *AgentService {
	return &AgentService{
		registryService: registryService,
		toSend:          make(chan *pb.ServerMessage, 100),
		commandResults:  make(chan *pb.CommandResult, 100),
		pendingCommands: make(map[string]*pb.Command),
		wg:              sync.WaitGroup{},
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

	s.registryService.CreateSession(idFromTLS)

	ctx, cancel := context.WithCancel(stream.Context())
	s.runReader(
		stream,
		func() {
			cancel()
		},
	)
	s.runWriter(ctx, stream)

	s.RequestSpecifications()
	s.wg.Wait()

	s.registryService.RemoveSession(idFromTLS)

	return nil
}

func (s *AgentService) runReader(stream pb.AgentService_ConnectServer, onClose func()) {
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

			case *pb.AgentMessage_CommandResult:
				commandResult := msg.GetCommandResult()
				log.Println("Command result received:", commandResult)
				s.commandResults <- commandResult

			default:
				log.Println("unknown message received:", msg)
			}
		}
	}()
}

func (s *AgentService) runWriter(ctx context.Context, stream pb.AgentService_ConnectServer) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case msg := <-s.toSend:
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
