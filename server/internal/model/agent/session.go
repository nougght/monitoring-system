package agent

import (
	"time"

	"github.com/google/uuid"
)

type AgentSession struct {
	AgentID      uuid.UUID
	ConnectionID uuid.UUID

	ConnectedAt time.Time
	LastSeenAt  time.Time // TODO: update db column from this field
}

func NewAgentSession(agentID uuid.UUID) *AgentSession {
	return &AgentSession{
		AgentID:      agentID,
		ConnectionID: uuid.New(),
		ConnectedAt:  time.Now(),
		LastSeenAt:   time.Now(),
	}
}
