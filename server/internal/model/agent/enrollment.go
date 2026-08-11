package agent

import (
	"time"

	"github.com/google/uuid"
)

const (
	DefaultAgentCertificateDuration time.Duration = time.Hour * 24 * 90
)

type EnrollmentKey struct {
	HashString string     `db:"key_hash"`
	AgentID    uuid.UUID  `db:"agent_id"`
	ExpiresAt  time.Time  `db:"expires_at"`
	UsedAt     *time.Time `db:"used_at"`
}

type EnrollParams struct {
	EnrollmentKey string
	CsrDer        []byte
}

type EnrollResult struct {
	CertDer    []byte   // agent certificate
	CAChainDer [][]byte // intermediate CA
	NotAfter   time.Time
}
