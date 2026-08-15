package mapper

import (
	"log"

	agent_model "github.com/nougght/monitoring-system/server/internal/model/agent"
	dto "github.com/nougght/monitoring-system/server/internal/transport/dto/types"
)

// -- agents --

func AgentToDTO(domain *agent_model.Agent) (res *dto.AgentDTO) {
	if domain == nil {
		return
	}
	isOnline := false
	if domain.IsOnline != nil {
		isOnline = *domain.IsOnline
	} else {
		log.Println("Warning: agent IsOnline field is nil, defaulting to false in DTO")
	}
	res = &dto.AgentDTO{
		ID:          domain.ID,
		Name:        domain.Name,
		Description: domain.Description,
		CreatedAt:   domain.CreatedAt,
		LastSeenAt:  domain.LastSeenAt,
		Status:      domain.Status,
		IsOnline:    isOnline,
	}
	return
}

func CreateAgentResultToDTO(domain *agent_model.CreateAgentResult) (res *dto.CreateAgentResponse) {
	if domain == nil {
		return
	}
	res = &dto.CreateAgentResponse{
		AgentDTO:      *AgentToDTO(&domain.Agent),
		EnrollmentKey: domain.EnrollmentKey,
	}
	return
}

func AgentConfigToDTO(domain *agent_model.AgentSetupConfig) (res *dto.AgentConfigResponse) {
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
