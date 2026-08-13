package agent

import (
	"fmt"

	"github.com/nougght/monitoring-system/server/internal/config"
	"github.com/nougght/monitoring-system/server/internal/storage/timescale/repository"
)

type AgentInteractionService struct {
	cfg       *config.Config
	agentRepo *repository.AgentRepository
}

func NewAgentService(cfg *config.Config, agentRepo *repository.AgentRepository) (*AgentInteractionService, error) {
	if cfg == nil || agentRepo == nil {
		return nil, fmt.Errorf("params required")
	}
	return &AgentInteractionService{
		cfg:       cfg,
		agentRepo: agentRepo,
	}, nil
}
