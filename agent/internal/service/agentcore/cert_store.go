package agentcore

import (
	"crypto/ecdsa"
	"crypto/tls"
	"fmt"
)

type CertStore struct {
	certPath string
	keyPath  string
}

func NewCertStore(certPath string, keyPath string) *CertStore {
	return &CertStore{
		certPath: certPath,
		keyPath:  keyPath,
	}
}
func (s *CertStore) LoadCertificate() (*tls.Certificate, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *CertStore) SaveCertificate(cert []byte) error {
	return fmt.Errorf("not implemented")
}

func (s *CertStore) LoadKey() (*ecdsa.PrivateKey, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *CertStore) SaveKey(key *ecdsa.PrivateKey) error {
	return fmt.Errorf("not implemented")
}
