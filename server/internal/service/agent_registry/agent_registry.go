package agentregistry

import (
	"context"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nougght/monitoring-system/server/internal/config"
	"github.com/nougght/monitoring-system/server/internal/model"
	agent_model "github.com/nougght/monitoring-system/server/internal/model/agent"
	"github.com/nougght/monitoring-system/server/internal/storage/timescale/repository"
	"github.com/nougght/monitoring-system/server/internal/util"
	utilShared "github.com/nougght/monitoring-system/shared/go/util"
)

type AgentRegistryService struct {
	cfg                *config.Config
	agentRepo          *repository.AgentRepository
	enrollmentKeysRepo *repository.EnrollmentKeysRepository
	specsRepo          *repository.SpecsRepository
	transactor         model.Transactor
	cert               *model.Certs
	sessions           map[uuid.UUID]*agent_model.AgentSession
	mu                 sync.RWMutex
}

func NewAgentRegistryService(cfg *config.Config, agentRepo *repository.AgentRepository,
	enrollmentKeysRepo *repository.EnrollmentKeysRepository,
	specsRepo *repository.SpecsRepository,
	transactor model.Transactor, cert *model.Certs) (*AgentRegistryService, error) {
	if cfg == nil || agentRepo == nil || enrollmentKeysRepo == nil || transactor == nil {
		return nil, fmt.Errorf("params required")
	}
	return &AgentRegistryService{
		cfg:                cfg,
		agentRepo:          agentRepo,
		enrollmentKeysRepo: enrollmentKeysRepo,
		specsRepo:          specsRepo,
		transactor:         transactor,
		cert:               cert,
		sessions:           make(map[uuid.UUID]*agent_model.AgentSession, 10),
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

func (s *AgentRegistryService) CreateSession(agentID uuid.UUID) *agent_model.AgentSession {
	s.mu.Lock()
	defer s.mu.Unlock()
	session := agent_model.NewAgentSession(agentID)
	s.sessions[agentID] = session
	return session
}

func (s *AgentRegistryService) GetSession(agentID uuid.UUID) (*agent_model.AgentSession, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, ok := s.sessions[agentID]
	return session, ok
}

func (s *AgentRegistryService) RemoveSession(agentID uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, agentID)
	log.Printf("session removed: %s", agentID)
}

func (s *AgentRegistryService) OnlineList() []uuid.UUID {
	s.mu.RLock()
	defer s.mu.RUnlock()
	onlineAgents := make([]uuid.UUID, 0, len(s.sessions))
	for agentID := range s.sessions {
		onlineAgents = append(onlineAgents, agentID)
	}
	return onlineAgents
}

func (s *AgentRegistryService) IsOnline(agentID uuid.UUID) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.sessions[agentID]
	return ok
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
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Printf("rollback failed: %s", err.Error())
		}
	}()
	ctx = context.WithValue(ctx, model.ContextKeyTx, tx)

	agent, err := s.agentRepo.CreateAgent(ctx, &agent_model.Agent{
		Name:        name,
		Description: description,
		Status:      utilShared.Ptr(string(agent_model.AgentStatusNotEnrolled)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	enrollmentKey := s.genEnrollmentKey()
	keyHash, err := util.Hash(*enrollmentKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate enrollment key: %w", err)
	}

	selector := uuid.New().String()
	key, err := s.enrollmentKeysRepo.CreateKey(ctx, &agent_model.EnrollmentKey{
		AgentID:    agent.ID,
		HashString: keyHash,
		Selector:   selector,
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
		EnrollmentKey: fmt.Sprintf("%s.%s", key.Selector, *enrollmentKey),
	}, nil
}

func (s *AgentRegistryService) GenerateAgentSetupConfig(ctx context.Context, agentID uuid.UUID, enrollmentKey string) *agent_model.AgentSetupConfig {
	return &agent_model.AgentSetupConfig{
		EnrollmentKey:     enrollmentKey,
		EnrollmentAddress: fmt.Sprintf("%s:%d", s.cfg.SettingsConfig.Address, s.cfg.GRPC.EnrollmentPort),
		ServerAddress:     fmt.Sprintf("%s:%d", s.cfg.SettingsConfig.Address, s.cfg.GRPC.MainPort),
	}
}

// TODO: return agent token
func (s *AgentRegistryService) Enroll(ctx context.Context, params *agent_model.EnrollParams) (*agent_model.EnrollResult, error) {
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
	ctx = context.WithValue(ctx, model.ContextKeyTx, tx)

	agentID, pubKey, err := s.validateAgentEnrollment(ctx, params.EnrollmentKey, params.CsrDer)
	if err != nil {
		log.Printf("failed to validate agent enrollment:%s", err.Error())
		return nil, fmt.Errorf("failed to validate agent enrollment")
	}
	log.Printf("agent id: %s", agentID.String())
	err = s.enrollmentKeysRepo.SetUsed(ctx, agentID, time.Now())
	if errors.Is(err, repository.ErrNoAffectedRows) {
		return nil, fmt.Errorf("enrollment key already used: %w", model.ErrBadRequest)
	}

	err = s.agentRepo.UpdateStatus(ctx, agentID, agent_model.AgentStatusActive)
	if err != nil {
		return nil, fmt.Errorf("failed update agent status: %w", err)
	}

	certNotAfter := time.Now().Add(agent_model.DefaultAgentCertificateDuration)
	agentCert, err := s.issueCertificate(agentID, pubKey, certNotAfter)
	if err != nil {
		log.Printf("failed issue agent certificate: %s", err.Error())
		return nil, fmt.Errorf("failed issue agent certificate")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("commit transaction error: %w", err)
	}
	return &agent_model.EnrollResult{
		CertDer:    agentCert.Raw,
		CAChainDer: [][]byte{s.cert.CA.Raw},
		NotAfter:   certNotAfter,
	}, nil
}

func (s *AgentRegistryService) validateAgentEnrollment(ctx context.Context, enrollmentKey string, csrDER []byte) (agentID uuid.UUID, pubKey any, err error) {
	csr, err := x509.ParseCertificateRequest(csrDER)
	if err != nil {
		return uuid.Nil, nil, err
	}
	if err := csr.CheckSignature(); err != nil {
		return uuid.Nil, nil, fmt.Errorf("CSR signature invalid: %w", err) // доказательство владения ключом провалено
	}

	// agentID, err = uuid.Parse(csr.Subject.CommonName)
	// log.Printf("agent common name: %s", csr.Subject.CommonName)
	// if err != nil {
	// 	return uuid.Nil, nil, err
	// }
	parts := strings.SplitN(enrollmentKey, ".", 2)
	if len(parts) != 2 {
		return uuid.Nil, nil, fmt.Errorf("invalid enrollment key format: %w", model.ErrBadRequest)
	}
	selector, keyVerifier := parts[0], parts[1]

	key, err := s.enrollmentKeysRepo.GetKeyBySelector(ctx, selector)
	if errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, nil, fmt.Errorf("key with selector '%s' not found: %w", selector, model.ErrBadRequest)
	}
	if err != nil {
		log.Printf("failed to get enrollment key by selector from bd: %s", err.Error())
		return uuid.Nil, nil, fmt.Errorf("failed validate enrollment key")
	}

	// сверяем с хэшированным значением из бд
	if ok, err := util.CompareHash(key.HashString, keyVerifier); !ok || err != nil {
		log.Printf("failed to compare hash enrollment key: %s", err.Error())
		return uuid.Nil, nil, fmt.Errorf("failed to validate enrollment key")
	}

	if key.ExpiresAt.Before(time.Now()) {
		return uuid.Nil, nil, fmt.Errorf("enrollment key expired: %w", model.ErrBadRequest)
	}

	return key.AgentID, csr.PublicKey, nil
}

func (s *AgentRegistryService) issueCertificate(agentID uuid.UUID, pubKey any, notAfter time.Time) (*x509.Certificate, error) {
	sn, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return nil, fmt.Errorf("failed to gen sn")
	}
	template := &x509.Certificate{
		SerialNumber: sn,
		Subject:      pkix.Name{CommonName: agentID.String()},
		NotBefore:    time.Now(),
		NotAfter:     notAfter,
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, template, s.cert.CA, pubKey, s.cert.Key)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(derBytes)
}

func (s *AgentRegistryService) GetAgentByID(ctx context.Context, id uuid.UUID) (*agent_model.Agent, error) {
	agent, err := s.agentRepo.GetAgentByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get agent by ID: %w", err)
	}

	agent.IsOnline = utilShared.Ptr(s.IsOnline(agent.ID))
	return agent, nil
}
func (s *AgentRegistryService) GetAllAgents(ctx context.Context) ([]*agent_model.Agent, error) {
	agents, err := s.agentRepo.GetAllAgents(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get agents: %w", err)
	}

	for i, agent := range agents {
		agents[i].IsOnline = utilShared.Ptr(s.IsOnline(agent.ID))
	}
	return agents, nil
}

// specifications

func (s *AgentRegistryService) GetSpecifications(ctx context.Context, agentID uuid.UUID) (*agent_model.Specs, error) {
	specs, err := s.specsRepo.GetCurrentSpecs(ctx, agentID)
	if errors.Is(err, repository.ErrNotFound) {
		err = model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get specifications: %w", err)
	}
	return specs, nil
}

func (s *AgentRegistryService) UpdateSpecifications(ctx context.Context, agentID uuid.UUID, specifications *agent_model.Specs) error {
	_, err := s.specsRepo.CreateOrUpdateSpecs(ctx, specifications)
	if err != nil {
		return fmt.Errorf("failed to update specifications: %w", err)
	}
	return nil
}
