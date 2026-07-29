package agent

import (
	"time"

	"github.com/google/uuid"
)

type EnrollmentKey struct {
	ID         uuid.UUID  `db:"id"`
	HashString string     `db:"key_hash"`
	AgentID    uuid.UUID  `db:"agent_id"`
	ExpiresAt  time.Time  `db:"expires_at"`
	UsedAt     *time.Time `db:"used_at"`
}
