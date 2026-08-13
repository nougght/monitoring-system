package agentcore

import (
	"agent/internal/config"
	"agent/internal/grpc/enrollment_client"
	"agent/internal/model"
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type CoreService struct {
	setupCfg  *config.SetupConfig
	state     *model.AgentState
	certStore model.CertStore
	cert      atomic.Pointer[tls.Certificate]
	ca        atomic.Pointer[x509.CertPool]
}

// func New(cfg Config, store CertStore) *Core {

// }

func NewCore(setupCfg *config.SetupConfig, certStore model.CertStore) (*CoreService, error) {
	cert, err := certStore.LoadCertificate()
	if err != nil {
		return nil, err
	}
	ca, err := certStore.LoadCA()

	agentID, err := getAgentIDFromCert(*cert)
	if err != nil {
		return nil, err
	}

	core := &CoreService{
		setupCfg:  setupCfg,
		state:     model.NewAgentState(agentID),
		certStore: certStore,
		cert:      atomic.Pointer[tls.Certificate]{},
		ca:        atomic.Pointer[x509.CertPool]{},
	}
	core.cert.Store(cert)
	core.ca.Store(ca)
	return core, nil
}

func (c *CoreService) State() *model.AgentState {
	return c.state
}

func (c *CoreService) TLSConfigForGRPC() *tls.Config {
	return &tls.Config{
		Certificates: []tls.Certificate{*c.cert.Load()},
		RootCAs:      c.ca.Load(),
		MinVersion:   tls.VersionTLS12,
	}
}

func EnrollAgent(ctx context.Context, setupCfg *config.SetupConfig, certStore model.CertStore) error {
	if setupCfg == nil {
		return fmt.Errorf(`setup config is required`)
	}

	if setupCfg.EnrollmentKey == model.EnrollmentKeyUsed {
		return fmt.Errorf(`agent is already enrolled, use 'agent run --setupconfig="/path/to/setup/config"`)
	}

	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	template := x509.CertificateRequest{
		Subject:            pkix.Name{},
		SignatureAlgorithm: x509.ECDSAWithSHA256,
	}

	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &template, privateKey)
	if err != nil {
		return fmt.Errorf(`failed to generate CSR: %w`, err)
	}

	caPool, err := certStore.LoadCA()
	if err != nil {
		return fmt.Errorf("failed to load CA")
	}

	client, err := setupEnrollmentClient(setupCfg, caPool)
	if err != nil {
		return fmt.Errorf("failed to setup enrollment client")
	}

	res, err := client.Enroll(ctx, &model.EnrollParams{
		EnrollmentKey: setupCfg.EnrollmentKey,
		CsrDer:        csrDER,
	})
	if err != nil {
		return err
	}

	// проверка полученного сертификата
	cert, err := verifyCertificate(res.CertDer, res.CAChainDer, caPool, privateKey)
	if err != nil {
		return fmt.Errorf("received invalid enroll result: %w", err)
	}

	err = setupCfg.SetEnrollmentKeyUsed()
	if err != nil {
		return fmt.Errorf("failed to set entollment key used: %w", err)
	}

	var buf bytes.Buffer
	err = pem.Encode(&buf, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
	if err != nil {
		return fmt.Errorf("failed to encode cert")
	}
	for _, der := range res.CAChainDer {
		err = pem.Encode(&buf, &pem.Block{Type: "CERTIFICATE", Bytes: der})
		if err != nil {
			return fmt.Errorf("failed to encode cert")
		}
	}
	// сохраняем сертификат
	err = certStore.SaveCertificate(buf.Bytes(), true)
	if err != nil {
		return err
	}
	// сохраняем ключ
	err = certStore.SaveKey(privateKey)
	if err != nil {
		return err
	}
	return nil
}

func setupEnrollmentClient(setupCfg *config.SetupConfig, caPool *x509.CertPool) (*enrollment_client.EnrollmentClient, error) {
	tlsConfig := &tls.Config{
		RootCAs:    caPool, // для проверки сертификата сервера
		MinVersion: tls.VersionTLS12,
		// без сертификата агента, т.к. его еще нет
	}
	conn, err := grpc.NewClient(
		setupCfg.EnrollmentAddress,
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc client: %w", err)
	}
	return enrollment_client.NewEnrollmentClient(conn, setupCfg)
}

func verifyCertificate(certDER []byte, chainDER [][]byte, trustedCA *x509.CertPool, priv crypto.Signer) (*x509.Certificate, error) {
	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cert: %w", err)
	}

	intermediates := x509.NewCertPool()
	for _, c := range chainDER {
		ic, err := x509.ParseCertificate(c)
		if err != nil {
			return nil, fmt.Errorf("failed to parse intermediate CAs: %w", err)
		}
		intermediates.AddCert(ic)
	}

	// проверка сертификата
	if _, err := cert.Verify(x509.VerifyOptions{
		Roots:         trustedCA,
		Intermediates: intermediates,
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}); err != nil {
		return nil, fmt.Errorf("failed to verify cert: %w", err)
	}

	// проверка соответствия публичного ключа приватному ключу агента
	certPub, ok := cert.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("incorrect cert public key type")
	}
	myPub := priv.Public().(*ecdsa.PublicKey)
	if !certPub.Equal(myPub) {
		return nil, errors.New("cert public key doesn't match agent private key")
	}

	// проверка срока
	if time.Now().After(cert.NotAfter) {
		return nil, errors.New("cert already expired")
	}
	if cert.Subject.CommonName == "" {
		return nil, errors.New("cert without CN")
	}

	return cert, nil
}

func getAgentIDFromCert(cert tls.Certificate) (uuid.UUID, error) {
	cn, err := getCN(cert)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get agent id from CN")
	}

	if err := uuid.Validate(cn); err != nil {
		return uuid.Nil, fmt.Errorf("invalid agent uuid")
	}

	id, err := uuid.Parse(cn)
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
