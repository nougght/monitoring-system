package agentregistry

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nougght/monitoring-system/server/internal/config"
	"github.com/nougght/monitoring-system/server/internal/model"
	"github.com/nougght/monitoring-system/server/internal/model/agent"
	agent_model "github.com/nougght/monitoring-system/server/internal/model/agent"
	"github.com/nougght/monitoring-system/server/internal/storage/timescale/repository"
	"github.com/nougght/monitoring-system/server/internal/util"
)

type AgentRegistryService struct {
	cfg                *config.Config
	agentRepo          *repository.AgentRepository
	enrollmentKeysRepo *repository.EnrollmentKeysRepository
	transactor         model.Transactor
}

func NewAgentRegistryService(cfg *config.Config, agentRepo *repository.AgentRepository,
	enrollmentKeysRepo *repository.EnrollmentKeysRepository, transactor model.Transactor) (*AgentRegistryService, error) {
	if cfg == nil || agentRepo == nil || enrollmentKeysRepo == nil || transactor == nil {
		return nil, fmt.Errorf("params required")
	}
	return &AgentRegistryService{
		cfg:                cfg,
		agentRepo:          agentRepo,
		enrollmentKeysRepo: enrollmentKeysRepo,
		transactor:         transactor,
	}, nil
}

func (s *AgentRegistryService) genEnrollmentKey() *string {
	rawKey := make([]byte, s.cfg.SettingsConfig.AgentEnrollmentKeyLength)
	_, err := rand.Read(rawKey)
	if err != nil {
		return nil
	}
	keyString := base64.RawStdEncoding.EncodeToString(rawKey)
	return &keyString
}

func (s *AgentRegistryService) CreateAgent(ctx context.Context, name string, description *string) (*agent_model.CreateAgentResult, error) {
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("agent name can't be blank: %w", model.ErrBadRequest)
	}

	tx, err := s.transactor.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed begin transaction: %w", err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Printf("rollback failed: %s", err.Error())
		}
	}()
	ctx = context.WithValue(ctx, model.TxKey, tx)

	agent, err := s.agentRepo.CreateAgent(ctx, &agent_model.Agent{
		Name:        name,
		Description: description,
		Status:      string(agent_model.AgentStatusNotEnrolled),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	enrollmentKey := s.genEnrollmentKey()
	keyHash, err := util.Hash(*enrollmentKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate enrollment key: %w", err)
	}

	_, err = s.enrollmentKeysRepo.CreateKey(ctx, &agent_model.EnrollmentKey{
		AgentID:    agent.ID,
		HashString: keyHash,
		ExpiresAt:  time.Now().Add(time.Minute * 20),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create enrollment key: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("commit transaction error: %w", err)
	}
	return &agent_model.CreateAgentResult{
		Agent:         *agent,
		EnrollmentKey: *enrollmentKey,
	}, nil
}

func (s *AgentRegistryService) GenerateAgentSetupConfig(ctx context.Context, agentID uuid.UUID, enrollmentKey string) *agent.AgentSetupConfig {
	return &agent_model.AgentSetupConfig{
		EnrollmentKey:     enrollmentKey,
		EnrollmentAddress: fmt.Sprintf("%s:%d", s.cfg.SettingsConfig.Address, s.cfg.GRPC.EnrollmentPort),
		ServerAddress:     fmt.Sprintf("%s:%d", s.cfg.SettingsConfig.Address, s.cfg.GRPC.MainPort),
	}
}

// TODO: return agent token
func (s *AgentRegistryService) EnrollAgent(ctx context.Context, agentID uuid.UUID, enrollmentKey string) error {
	tx, err := s.transactor.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed begin transaction: %w", err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Printf("rollback failed: %s", err.Error())
		}
	}()
	ctx = context.WithValue(ctx, model.TxKey, tx)

	key, err := s.enrollmentKeysRepo.GetKeyByAgentId(ctx, agentID)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("key with id '%s' not found: %w", agentID.String(), model.ErrNotFound)
	}
	if err != nil {
		return fmt.Errorf("failed to get enrollment key: %w", err)
	}

	if key.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("enrollment key expired: %w", model.ErrBadRequest)
	}

	err = s.enrollmentKeysRepo.SetUsed(ctx, agentID, time.Now())
	if errors.Is(err, repository.ErrNoAffectedRows) {
		return fmt.Errorf("enrollment key already used: %w", model.ErrBadRequest)
	}

	err = s.agentRepo.UpdateStatus(ctx, agentID, agent_model.AgentStatusActive)
	if err != nil {
		return fmt.Errorf("failed update agent status: %w", err)
	}

	return nil
}

func (s *AgentRegistryService) GetAllAgents(ctx context.Context) ([]*agent_model.Agent, error) {
	agents, err := s.agentRepo.GetAllAgents(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get agents: %w", err)
	}

	return agents, nil
}
