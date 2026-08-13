package model

import (
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
)

type Certs struct {
	CA     *x509.Certificate
	Key    *ecdsa.PrivateKey
	RootCA *x509.CertPool
}
type CertStore interface {
	LoadCertificate() (*tls.Certificate, error)
	SaveCertificate(cert []byte) error
	SaveKey(key *ecdsa.PrivateKey) error
	LoadKey() (*ecdsa.PrivateKey, error)
	LoadCA() (*x509.CertPool, error)
}
