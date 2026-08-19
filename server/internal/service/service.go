package service

import (
	"log"

	"github.com/nougght/monitoring-system/server/internal/config"
	"github.com/nougght/monitoring-system/server/internal/model"
	agent "github.com/nougght/monitoring-system/server/internal/service/agent_interaction"
	agentregistry "github.com/nougght/monitoring-system/server/internal/service/agent_registry"
	"github.com/nougght/monitoring-system/server/internal/storage/timescale/repository"
)

type Services struct {
	agentRegistry *agentregistry.AgentRegistryService
	agent         *agent.AgentInteractionService
}

//nolint:unused
type ServicesOptions struct {
	Config       *config.Config
	Repositories *repository.Repositories
	Transactor   model.Transactor
	Cert         *model.Certs
}

func New(opts ServicesOptions) *Services {
	agentRegistry, err := agentregistry.NewAgentRegistryService(
		opts.Config,
		opts.Repositories.AgentRepository(),
		opts.Repositories.EnrollmentKeysRepository(),
		opts.Repositories.SpecsRepository(),
		opts.Transactor,
		opts.Cert,
	)
	if err != nil {
		log.Panicf("failed initialize agent registry: %s", err.Error())
	}

	agentInteraction, err := agent.NewAgentInteractionService(
		opts.Config,
		agentRegistry,
	)
	return &Services{
		agentRegistry: agentRegistry,
		agent:         agentInteraction,
	}
}

func (s *Services) AgentRegistry() *agentregistry.AgentRegistryService {
	return s.agentRegistry
}

func (s *Services) AgentInteractionService() *agent.AgentInteractionService {
	return s.agent
}
