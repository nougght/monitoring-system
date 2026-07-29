package agentregistry

import (
	"fmt"

	"github.com/nougght/monitoring-system/server/internal/config"
	"github.com/nougght/monitoring-system/server/internal/storage/timescale/repository"
)

type AgentRegistryService struct {
	cfg                *config.Config
	agentRepo          *repository.AgentRepository
	enrollmentKeysRepo *repository.EnrollmentKeysRepository
}

func NewAgentRegistryService(cfg *config.Config, agentRepo *repository.AgentRepository,
	enrollmentKeysRepo *repository.EnrollmentKeysRepository) (*AgentRegistryService, error) {
	if cfg == nil || agentRepo == nil || enrollmentKeysRepo == nil {
		return nil, fmt.Errorf("params required")
	}
	return &AgentRegistryService{
		cfg:                cfg,
		agentRepo:          agentRepo,
		enrollmentKeysRepo: enrollmentKeysRepo,
	}, nil
}

func (s *AgentRegistryService) genEnrollmentKey() {

}
