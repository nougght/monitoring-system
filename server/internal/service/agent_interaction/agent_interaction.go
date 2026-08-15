package agent

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nougght/monitoring-system/server/internal/config"
	agent_model "github.com/nougght/monitoring-system/server/internal/model/agent"
	agentregistry "github.com/nougght/monitoring-system/server/internal/service/agent_registry"
)

type AgentInteractionService struct {
	cfg      *config.Config
	registry *agentregistry.AgentRegistryService
}

func NewAgentInteractionService(cfg *config.Config, registry *agentregistry.AgentRegistryService) (*AgentInteractionService, error) {
	if cfg == nil || registry == nil {
		return nil, fmt.Errorf("params required")
	}
	return &AgentInteractionService{
		cfg:      cfg,
		registry: registry,
	}, nil
}

func (s *AgentInteractionService) HandleConnection(agentID uuid.UUID) {
	s.registry.CreateSession(agentID)
}

func (s *AgentInteractionService) HandleDisconnection(agentID uuid.UUID) {
	s.registry.RemoveSession(agentID)
}

func (s *AgentInteractionService) Enroll(ctx context.Context, params *agent_model.EnrollParams) (*agent_model.EnrollResult, error) {
	return s.registry.Enroll(ctx, params)
}
