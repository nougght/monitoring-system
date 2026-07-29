package dto

import (
	"time"

	"github.com/google/uuid"
)

type AgentDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	LastSeenAt  time.Time `json:"lastSeenAt"`
	IsOnline    bool      `json:"isOnline"`
}

type EnrollmentResponse struct {
	EnrollmentKey string
	AgentID       uuid.UUID
}
