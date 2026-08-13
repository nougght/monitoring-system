package model

import (
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
)

type CertStore interface {
	LoadCertificate() (*tls.Certificate, error)
	SaveCertificate(cert []byte, encoded bool) error
	SaveKey(key *ecdsa.PrivateKey) error
	LoadKey() (*ecdsa.PrivateKey, error)
	LoadCA() (*x509.CertPool, error)
}
