package agent

import (
	"time"

	"github.com/google/uuid"
)

type AgentStatus string

const (
	AgentStatusNotEnrolled AgentStatus = "not_enrolled"
	AgentStatusActive      AgentStatus = "active"
)

type Agent struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description *string   `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	LastSeenAt  time.Time `db:"last_seen_at"`
	Status      string    `db:"status"`
	IsOnline    bool
}

type CreateAgentResult struct {
	Agent
	EnrollmentKey string
}
