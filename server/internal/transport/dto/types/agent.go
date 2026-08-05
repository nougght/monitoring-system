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
} // @Name Agent

type EnrollmentResponse struct {
	EnrollmentKey string
	AgentID       uuid.UUID
} // @Name EnrollmentResponse

type CreateAgentBody struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
} // @Name CreateAgentBody

type CreateAgentResponse struct {
	AgentDTO
	EnrollmentKey string `json:"enrollmentKey"`
} // @Name CreateAgentResponse

type AgentConfigBody struct {
	EnrollmentKey string `json:"enrollmentKey"`
}

type AgentConfigResponse struct {
	EnrollmentKey     string `yaml:"enrollment_key"`
	EnrollmentAddress string `yaml:"enrollment_address"`
	ServerAddress     string `yaml:"server_address"`
}
