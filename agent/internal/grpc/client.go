package grpc_client

import (
	"agent/internal/model"
	"context"
	"errors"
	"io"
	"log"

	pb "github.com/nougght/monitoring-system/shared/go/proto/gen/agent/v1"
	"google.golang.org/grpc"
)

type MetricsProvider interface {
	GetMetrics() *model.Metrics
	GetSpecs(ctx context.Context) (*model.Specs, error)
}

type AgentClient struct {
	conn            *grpc.ClientConn
	grpcClient      pb.AgentServiceClient
	metricsProvider MetricsProvider
}

func NewAgentClient(conn *grpc.ClientConn, metricsProvider MetricsProvider) *AgentClient {
	return &AgentClient{
		conn:            conn,
		grpcClient:      pb.NewAgentServiceClient(conn),
		metricsProvider: metricsProvider,
	}
}

func (c *AgentClient) Connect(ctx context.Context) error {
	stream, err := c.grpcClient.Connect(ctx)
	if err != nil {
		return err
	}
	c.runReader(stream)
	// c.runWriter(stream)
	return nil
}

func (c *AgentClient) runReader(stream pb.AgentService_ConnectClient) {
	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					log.Println("grpc stream closed")
					return
				}
				log.Println("error receiving message from grpc stream:", err)
				return
			}
			log.Println("message received:", msg)

			switch msg.Payload.(type) {
			case *pb.ServerMessage_Command:
				command := msg.GetCommand()
				log.Println("command received:", command)
				switch command.Payload.(type) {
				case *pb.Command_SpecificationsRequest:
					specs, err := c.metricsProvider.GetSpecs(stream.Context())
					if err != nil {
						log.Println("error getting specs:", err)
						continue
					}
					err = stream.Send(&pb.AgentMessage{
						Payload: &pb.AgentMessage_CommandResult{
							CommandResult: &pb.CommandResult{
								Payload: &pb.CommandResult_SpecificationsResponse{
									SpecificationsResponse: &pb.SpecificationsResponse{
										Specs: convertSpecsToProto(specs),
									}}}},
					})
					if err != nil {
						log.Println("error sending specifications response:", err)
						return
					}

				}
			}
		}
	}()
}

// func (c *AgentClient) runWriter(stream pb.AgentService_ConnectClient) error {
// 	ticker := time.NewTicker(1 * time.Second)
// 	go func() {
// 		for {
// 			metrics, err := c.metricsProvider.GetMetrics()
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}()
// 	return nil
// }
