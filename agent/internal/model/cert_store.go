package model

import "crypto/tls"

type CertStore interface {
	LoadCertificate() (*tls.Certificate, error)
	SaveCertificate(cert []byte) error
}
