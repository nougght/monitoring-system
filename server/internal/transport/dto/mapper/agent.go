package mapper

import (
	"github.com/nougght/monitoring-system/server/internal/model/agent"
	dto "github.com/nougght/monitoring-system/server/internal/transport/dto/types"
)

// -- agents --

func AgentToDTO(domain *agent.Agent, isOnline bool) (res *dto.AgentDTO, err error) {
	if domain == nil {
		return res, nil
	}
	return &dto.AgentDTO{
		ID:          domain.ID,
		Name:        domain.Name,
		Description: domain.Description,
		CreatedAt:   domain.CreatedAt,
		LastSeenAt:  domain.LastSeenAt,
		IsOnline:    isOnline,
	}, nil
}
