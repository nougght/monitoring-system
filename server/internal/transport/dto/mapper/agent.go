package mapper

import (
	"github.com/nougght/monitoring-system/server/internal/model/agent"
	dto "github.com/nougght/monitoring-system/server/internal/transport/dto/types"
)

// -- agents --

func AgentToDTO(domain *agent.Agent) (res *dto.AgentDTO) {
	if domain == nil {
		return
	}
	res = &dto.AgentDTO{
		ID:          domain.ID,
		Name:        domain.Name,
		Description: domain.Description,
		CreatedAt:   domain.CreatedAt,
		LastSeenAt:  domain.LastSeenAt,
		IsOnline:    domain.IsOnline,
	}
	return
}

func CreateAgentResultToDTO(domain *agent.CreateAgentResult) (res *dto.CreateAgentResponse) {
	if domain == nil {
		return
	}
	res = &dto.CreateAgentResponse{
		AgentDTO:      *AgentToDTO(&domain.Agent),
		EnrollmentKey: domain.EnrollmentKey,
	}
	return
}

func AgentConfigToDTO(domain *agent.AgentSetupConfig) (res *dto.AgentConfigResponse) {
	if domain == nil {
		return
	}
	res = &dto.AgentConfigResponse{
		EnrollmentKey:     domain.EnrollmentKey,
		EnrollmentAddress: domain.EnrollmentAddress,
		ServerAddress:     domain.ServerAddress,
	}
	return
}
