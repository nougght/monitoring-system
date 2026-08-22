package agent

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/nougght/monitoring-system/server/internal/config"
	agent_model "github.com/nougght/monitoring-system/server/internal/model/agent"
	"github.com/nougght/monitoring-system/server/internal/model/metrics"
	agentregistry "github.com/nougght/monitoring-system/server/internal/service/agent_registry"
)

// TEMP: replace with single SendCommand method
type Requester interface {
	RequestSpecifications(ctx context.Context, agentID uuid.UUID) (*agent_model.Specs, error)
}

type AgentInteractionService struct {
	cfg       *config.Config
	registry  *agentregistry.AgentRegistryService
	requester Requester

	// lastFrames  map[uuid.UUID][]byte
	viewerChans map[uuid.UUID]map[uuid.UUID]chan []byte
	mu          sync.RWMutex
}

func NewAgentInteractionService(cfg *config.Config, registry *agentregistry.AgentRegistryService) (*AgentInteractionService, error) {
	if cfg == nil || registry == nil {
		return nil, fmt.Errorf("params required")
	}
	return &AgentInteractionService{
		cfg:         cfg,
		registry:    registry,
		viewerChans: make(map[uuid.UUID]map[uuid.UUID]chan []byte),
	}, nil
}

func (s *AgentInteractionService) SetRequester(requester Requester) {
	s.requester = requester
}

func (s *AgentInteractionService) HandleConnection(agentID uuid.UUID) {
	s.registry.CreateSession(agentID)
	specs, err := s.RequestSpecifications(context.Background(), agentID)
	if err != nil {
		log.Printf("failed to request specifications: %s", err.Error())
		return
	}
	specs.AgentID = agentID
	err = s.registry.UpdateSpecifications(context.Background(), agentID, specs)
	if err != nil {
		log.Printf("failed to update specifications: %s", err.Error())
	}

}

func (s *AgentInteractionService) HandleDisconnection(agentID uuid.UUID) {
	s.registry.RemoveSession(agentID)
}

func (s *AgentInteractionService) Enroll(ctx context.Context, params *agent_model.EnrollParams) (*agent_model.EnrollResult, error) {
	return s.registry.Enroll(ctx, params)
}

func (s *AgentInteractionService) RequestSpecifications(ctx context.Context, agentID uuid.UUID) (*agent_model.Specs, error) {
	if s.requester == nil {
		return nil, fmt.Errorf("requester is nil")
	}
	return s.requester.RequestSpecifications(ctx, agentID)
}

func (s *AgentInteractionService) HandleSpecifications(ctx context.Context, agentID uuid.UUID, specifications *agent_model.Specs) error {
	return s.registry.UpdateSpecifications(ctx, agentID, specifications)
}

func (s *AgentInteractionService) HandleMetricsBatch(ctx context.Context, batch *metrics.MetricsBatch) {

}

func (s *AgentInteractionService) SubStreaming(agentID, viewerID uuid.UUID) (<-chan []byte, error) {
	_, ok := s.registry.GetSession(agentID)
	if !ok {
		// return nil, fmt.Errorf("agent is offline: %w", model.ErrServiceUnavailable)
	}
	// log.Printf("sub stream %s", agentID.String())
	framesChan := make(chan []byte, 1)
	s.mu.Lock()
	defer s.mu.Unlock()
	subs := s.viewerChans[agentID]
	if subs == nil {
		subs = make(map[uuid.UUID]chan []byte, 1)
		s.viewerChans[agentID] = subs
	}
	subs[viewerID] = framesChan
	return framesChan, nil
}

func (s *AgentInteractionService) UnsubAllStreaming(viewerID uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, subs := range s.viewerChans {
		_, ok := subs[viewerID]
		if ok {
			delete(subs, viewerID)
		}
	}
}

func (s *AgentInteractionService) HandleFrame(frame []byte, agentID uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	subs, ok := s.viewerChans[agentID]
	if !ok {
		log.Printf("no subs for agent: %s", agentID.String())
		return
	}
	if len(subs) == 0 {
		log.Println("not found stream subs")
	}
	for _, ch := range subs {
		select {
		case ch <- frame:
		default: // skip old frame
			select {
			case <-ch:
			default:
			}
			select { // send new frame
			case ch <- frame:
			default:
			}
		}
	}
}
