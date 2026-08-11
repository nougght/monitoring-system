package cert_store

import (
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

type CertStore struct {
	certPath string
	keyPath  string
	caPath   string
}

func NewCertStore(certPath string, keyPath string, caPath string) *CertStore {
	return &CertStore{
		certPath: certPath,
		keyPath:  keyPath,
		caPath:   caPath,
	}
}
func (s *CertStore) LoadCertificate() (*tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(s.certPath, s.keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load agent cert")
	}
	return &cert, nil
}

// TODO: сделать безопасное сохранение/редактирование файлов
func (s *CertStore) SaveCertificate(cert []byte) error {
	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	}

	certFile, err := os.Open(s.certPath)
	if err != nil {
		return err
	}
	_, err = certFile.Write(pem.EncodeToMemory(block))
	if err != nil {
		return err
	}
	return nil
}

func (s *CertStore) LoadKey() (*ecdsa.PrivateKey, error) {
	keyFile, err := os.ReadFile(s.keyPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyFile)
	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse agent key")
	}
	return key, nil
}

func (s *CertStore) SaveKey(key *ecdsa.PrivateKey) error {
	keyBytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return err
	}
	block := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: keyBytes,
	}

	keyFile, err := os.Open(s.keyPath)
	if err != nil {
		return err
	}
	_, err = keyFile.Write(pem.EncodeToMemory(block))
	if err != nil {
		return err
	}
	return nil
}

func (s *CertStore) LoadCA() (*x509.CertPool, error) {
	caFile, err := os.ReadFile(s.caPath)
	if err != nil {
		return nil, err
	}

	caPool := x509.NewCertPool()
	ok := caPool.AppendCertsFromPEM(caFile)
	if !ok {
		return nil, fmt.Errorf("failed to add CA cert")
	}
	return caPool, nil
}
