package model

import (
	"crypto/ecdsa"
	"crypto/tls"
)

type CertStore interface {
	LoadCertificate() (*tls.Certificate, error)
	SaveCertificate(cert []byte) error
	SaveKey(key *ecdsa.PrivateKey) error
	LoadKey() (*ecdsa.PrivateKey, error)
}
