package grpc_handler

import (
	"errors"
	"io"
	"log"
	"sync"

	"github.com/google/uuid"
	pb "github.com/nougght/monitoring-system/shared/go/proto/gen/agent/v1"
	"google.golang.org/grpc"
)

type AgentService struct {
	pb.UnimplementedAgentServiceServer
	// metricsService metrics.MetricsService
	toSend          chan *pb.ServerMessage
	commandResults  chan *pb.CommandResult
	pendingCommands map[string]*pb.Command
	mu              sync.Mutex
	wg              sync.WaitGroup
}

func NewAgentService() *AgentService {
	return &AgentService{
		toSend:          make(chan *pb.ServerMessage),
		commandResults:  make(chan *pb.CommandResult),
		pendingCommands: make(map[string]*pb.Command),
		wg:              sync.WaitGroup{},
	}
}

func (s *AgentService) Register(server *grpc.Server) {
	pb.RegisterAgentServiceServer(server, s)
}

func (s *AgentService) Connect(stream pb.AgentService_ConnectServer) error {
	log.Println("grpc client connected")
	s.runReader(stream)
	s.runWriter(stream)
	s.RequestSpecifications()
	s.wg.Wait()
	return nil
}

func (s *AgentService) runReader(stream pb.AgentService_ConnectServer) {
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

func (s *AgentService) runWriter(stream pb.AgentService_ConnectServer) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			msg := <-s.toSend
			err := stream.Send(msg)
			if err != nil {
				log.Println("error sending message to grpc client:", err)
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
