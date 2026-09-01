package service

import (
	"log"

	"github.com/nougght/monitoring-system/server/internal/config"
	"github.com/nougght/monitoring-system/server/internal/model"
	agent "github.com/nougght/monitoring-system/server/internal/service/agent_interaction"
	agentregistry "github.com/nougght/monitoring-system/server/internal/service/agent_registry"
	"github.com/nougght/monitoring-system/server/internal/service/metrics"
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
	metrics, err := metrics.NewMetricsService(
		opts.Config,
		opts.Transactor,
		opts.Repositories.MetricsRepository,
	)
	agentInteraction, err := agent.NewAgentInteractionService(
		opts.Config,
		agentRegistry,
		metrics,
	)
	if err != nil {
		log.Panicf("failed initialize agent interaction service: %s", err.Error())
	}
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
