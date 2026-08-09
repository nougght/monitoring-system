package model

import "time"

type EnrollParams struct {
	EnrollmentKey string
	CsrDer        []byte
}

type EnrollResult struct {
	CertDer    []byte   // agent certificate
	CAChainDer [][]byte // intermediate CA
	NotAfter   time.Time
}
