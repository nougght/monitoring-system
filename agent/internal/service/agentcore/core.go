package agentcore

import (
	"agent/internal/config"
	"agent/internal/model"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"sync/atomic"

	"github.com/google/uuid"
)

type CoreService struct {
	setupCfg  *config.SetupConfig
	state     *model.AgentState
	certStore model.CertStore
	cert      atomic.Pointer[tls.Certificate]
}

// func New(cfg Config, store CertStore) *Core {

// }

func NewCore(setupCfg *config.SetupConfig, certStore model.CertStore) (*CoreService, error) {
	cert, err := certStore.LoadCertificate()
	if err != nil {
		return nil, err
	}
	agentID, err := getAgentIDFromCert(*cert)
	if err != nil {
		return nil, err
	}

	core := &CoreService{
		setupCfg:  setupCfg,
		state:     model.NewAgentState(agentID),
		certStore: certStore,
		cert:      atomic.Pointer[tls.Certificate]{},
	}
	core.cert.Store(cert)
	return core, nil
}

func (c *CoreService) State() *model.AgentState {
	return c.state
}

func getAgentIDFromCert(cert tls.Certificate) (uuid.UUID, error) {
	cn, err := getCN(cert)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get agent id from CN")
	}

	if err := uuid.Validate(cn); err != nil {
		return uuid.Nil, fmt.Errorf("invalid agent uuid")
	}

	id, err := uuid.FromBytes([]byte(cn))
	if err != nil {
		return uuid.Nil, fmt.Errorf("faile get uuid from cn")
	}
	return id, nil
}

func getCN(cert tls.Certificate) (string, error) {
	cn := tryGetCNFromCert(cert)
	if cn != nil {
		return *cn, nil
	}

	parsed, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return "", err
	}
	return parsed.Subject.CommonName, nil
}

func tryGetCNFromCert(cert tls.Certificate) *string {
	if cert.Leaf == nil {
		return nil
	}
	return &cert.Leaf.Subject.CommonName
}

// func (c *CoreService) Enroll(ctx context.Context) error {

// }
// func (c *CoreService) EnsureValidCertificate(ctx context.Context) error {

// }
