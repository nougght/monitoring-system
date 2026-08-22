package agent_client

import (
	"agent/internal/config"
	"agent/internal/model"
	"agent/internal/service/streaming"
	"bytes"
	"context"
	"errors"
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"sync"
	"time"

	pb "github.com/nougght/monitoring-system/shared/go/proto/gen/agent/v1"
	"google.golang.org/grpc"
)

// TODO: add writer channels
type MetricsProvider interface {
	GetMetrics() *model.Metrics
	GetSpecs(ctx context.Context) (*model.Specs, error)
}

type AgentClient struct {
	conn            *grpc.ClientConn
	config          *config.Config
	state           *model.AgentState
	grpcClient      pb.AgentServiceClient
	metricsProvider MetricsProvider
}

func NewAgentClient(conn *grpc.ClientConn, config *config.Config, state *model.AgentState, metricsProvider MetricsProvider) *AgentClient {
	return &AgentClient{
		conn:            conn,
		config:          config,
		state:           state,
		grpcClient:      pb.NewAgentServiceClient(conn),
		metricsProvider: metricsProvider,
	}
}

func (c *AgentClient) Connect(ctx context.Context) error {
	c.state.SetPending(true)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	stream, err := c.grpcClient.Connect(ctx)
	if err != nil {
		return err
	}

	err = stream.Send(
		&pb.AgentMessage{
			Payload: &pb.AgentMessage_Handshake{
				Handshake: &pb.Handshake{
					AgentUuid: c.state.AgentID().String(),
				},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to send handshake message: %w", err)
	}

	msg, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("failed to recv handshake status message: %w", err)
	}
	switch msg.Payload.(type) {
	case *pb.ServerMessage_Status:
		connected := msg.GetStatus().Connected
		log.Printf("handshake status received: %v", connected)
		if !connected {
			return fmt.Errorf("status not connected")
		}
	default:
		return fmt.Errorf("expectede server handshake status message")
	}
	c.state.SetConnected(true)
	wg := sync.WaitGroup{}
	c.runReader(stream, &wg, func() { cancel() })
	err = c.runWriter(ctx, &wg, stream)
	if err != nil {
		log.Println(err.Error())
	}
	wg.Wait()
	c.state.SetConnected(false)
	return nil
}

func (c *AgentClient) runReader(stream pb.AgentService_ConnectClient, wg *sync.WaitGroup, onClose func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			msg, err := stream.Recv()
			if err != nil {
				onClose()
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
				log.Println("command received:", command.CommandUuid)
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
								CommandUuid: command.GetCommandUuid(),
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
			case *pb.ServerMessage_Status:
				connected := msg.GetStatus().Connected
				log.Printf("status msg received: %v", connected)
				if !connected {
					onClose()
				}
			default:
				log.Println("unknow msg type received")
			}
		}
	}()
}

func (c *AgentClient) runWriter(ctx context.Context, wg *sync.WaitGroup, stream pb.AgentService_ConnectClient) error {
	ticker := time.NewTicker(c.config.MetricsSendingInterval)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ticker.C:
				metrics := c.metricsProvider.GetMetrics()
				err := stream.Send(&pb.AgentMessage{
					Payload: &pb.AgentMessage_Metrics{
						Metrics: convertMetricsToProto(metrics),
					},
				})
				if err != nil {
					log.Println("error sending metrics:", err)
					return
				}
			case <-ctx.Done():
				log.Println("writer ctx closed")
				if err := stream.CloseSend(); err != nil {
					log.Printf("main grpc close send failed: %s", err.Error())
				}
				log.Println("main grpc close send")
				return
			}
		}
	}()
	return nil
}

func (c *AgentClient) StartStreamMJPEG(ctx context.Context) error {
	// connect only if main grpc stream connected
	connectedChan := make(chan bool)
	log.Println("check for main stream connected")
	if !c.state.Connected() {
		go func() {
			ticker := time.NewTicker(time.Millisecond * 200)
			for {
				select {
				case <-ctx.Done():
					connectedChan <- false
					return
				case <-ticker.C:
					if !c.state.Connected() && !c.state.Pending() {
						connectedChan <- false
						return
					}
					if c.state.Connected() {
						connectedChan <- true
						return
					}
				}
			}
		}()
		connected := <-connectedChan
		if !connected {
			return fmt.Errorf("main grpc stream is not connected")
		}
	}
	stream, err := c.grpcClient.StartStreamMJPEG(ctx)
	if err != nil {
		log.Printf("failed connect streaming: %s", err.Error())
		log.Println("start retrying")
		t := time.NewTicker(time.Millisecond * 200)
		cnt := 0
		for {
			select {
			case <-t.C:
				stream, err = c.grpcClient.StartStreamMJPEG(ctx)
				cnt += 1
				if err != nil {
					log.Printf("retry %d failed: %s", cnt, err.Error())
				}
				if cnt > 10 {
					log.Println("connect streaming retrying failed")
				}
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return err
	}
	log.Println("streaming grpc started")

	wg := sync.WaitGroup{}
	// c.runStreamingReader(stream)
	err = c.runStreamingWriter(ctx, &wg, stream, time.Millisecond*60)
	if err != nil {
		log.Println(err.Error())
	}
	wg.Wait()
	err = stream.CloseSend()
	if err != nil {
		return fmt.Errorf("streaming grpc close send failed: %w", err)
	}
	log.Println("streaming close send")
	return nil
}

func (c *AgentClient) runStreamingWriter(ctx context.Context, wg *sync.WaitGroup, stream pb.AgentService_StartStreamMJPEGClient, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ticker.C:
				image, err := streaming.TakeScreenshot()
				if err != nil {
					log.Println(err.Error())
				}
				var raw bytes.Buffer
				err = jpeg.Encode(&raw, image, &jpeg.Options{Quality: 60})
				if err != nil {
					log.Printf("failed to encode frame: %w", err)
				}
				err = stream.Send(&pb.StreamingAgentMessage{
					Payload: &pb.StreamingAgentMessage_Frame{
						Frame: &pb.VideoFrame{
							Data: raw.Bytes(),
						},
					},
				})
				if err != nil {
					log.Println("error sending frame:", err)
					return
				}
				log.Println("frame send")
			case <-ctx.Done():
				log.Println("streaming writer ctx canceled")
				return
			}
		}
	}()
	return nil
}
