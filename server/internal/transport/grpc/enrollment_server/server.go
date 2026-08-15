package enrollment_server

import (
	"context"

	agent_model "github.com/nougght/monitoring-system/server/internal/model/agent"
	pb "github.com/nougght/monitoring-system/shared/go/proto/gen/enrollment/v1"
	"google.golang.org/grpc"
)

type Enroller interface {
	Enroll(ctx context.Context, params *agent_model.EnrollParams) (*agent_model.EnrollResult, error)
}

type EnrollmentService struct {
	pb.UnimplementedAgentEnrollmentServiceServer
	enroller Enroller
}

func NewEnrollmentService(enroller Enroller) *EnrollmentService {
	return &EnrollmentService{
		enroller: enroller,
	}
}

func (s *EnrollmentService) Register(server *grpc.Server) {
	pb.RegisterAgentEnrollmentServiceServer(server, s)
}

func (s *EnrollmentService) Enroll(ctx context.Context, req *pb.EnrollRequest) (*pb.EnrollResponse, error) {
	params := convertEnrollParamsFromProto(req)
	result, err := s.enroller.Enroll(ctx, params)
	if err != nil {
		return nil, err
	}
	return convertEnrollResultToProto(result), nil
}
