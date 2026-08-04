package agentcore

import (
	"crypto/tls"
	"fmt"
)

type CertStore struct {
	certPath string
}

func NewCertStore(certPath string) *CertStore {
	return &CertStore{
		certPath: certPath,
	}
}
func (s *CertStore) LoadCertificate() (*tls.Certificate, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *CertStore) SaveCertificate(cert []byte) error {
	return fmt.Errorf("not implemented")
}
