package enrollment_client

import (
	"agent/internal/model"

	pb "github.com/nougght/monitoring-system/shared/go/proto/gen/enrollment/v1"
)

func convertEnrollParamsToProto(p *model.EnrollParams) *pb.EnrollRequest {
	if p == nil {
		return nil
	}
	return &pb.EnrollRequest{
		EnrollmentKey: p.EnrollmentKey,
		CsrDer:        p.CsrDer,
	}
}

func convertEnrollResultFromProto(p *pb.EnrollResponse) *model.EnrollResult {
	if p == nil {
		return nil
	}
	return &model.EnrollResult{
		CertDer:    p.CertificateDer,
		CAChainDer: p.CaChainDer,
		NotAfter:   p.NotAfter.AsTime(),
	}
}
