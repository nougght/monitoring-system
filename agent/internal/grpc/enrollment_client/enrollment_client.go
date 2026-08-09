package enrollment_client

import (
	"agent/internal/config"
	"agent/internal/model"
	"context"
	"fmt"

	pb "github.com/nougght/monitoring-system/shared/go/proto/gen/enrollment/v1"
	"google.golang.org/grpc"
)

type EnrollmentClient struct {
	conn       *grpc.ClientConn
	config     *config.SetupConfig
	grpcClient pb.AgentEnrollmentServiceClient
}

func NewEnrollmentClient(conn *grpc.ClientConn, setupCfg *config.SetupConfig) (*EnrollmentClient, error) {
	grpcClient := pb.NewAgentEnrollmentServiceClient(conn)
	return &EnrollmentClient{
		conn:       conn,
		config:     setupCfg,
		grpcClient: grpcClient,
	}, nil
}

func (c *EnrollmentClient) Enroll(ctx context.Context, params *model.EnrollParams) (*model.EnrollResult, error) {
	resp, err := c.grpcClient.Enroll(ctx, convertEnrollParamsToProto(params))
	if err != nil {
		return nil, fmt.Errorf("enroll agent request failed")
	}
	return convertEnrollResultFromProto(resp), nil
}
