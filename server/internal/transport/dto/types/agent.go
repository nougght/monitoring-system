package dto

import (
	"time"

	"github.com/google/uuid"
	agent_model "github.com/nougght/monitoring-system/server/internal/model/agent"
)

type AgentDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	Status      *string   `json:"status,omitempty"`
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
} // @Name AgentConfigBody

type AgentConfigResponse struct {
	EnrollmentKey     string `yaml:"enrollment_key"`
	EnrollmentAddress string `yaml:"enrollment_address"`
	ServerAddress     string `yaml:"server_address"`
} // @Name AgentConfigResponse

// TODO: check if need
type SpecsDTO struct {
	AgentID     uuid.UUID                `json:"agentId"`
	UpdatedAt   time.Time                `json:"updatedAt"`
	HostSpecs   agent_model.HostSpecs    `json:"hostSpecs"`
	CpuSpecs    agent_model.CpuSpecs     `json:"cpuSpecs"`
	DiskSpecs   []*agent_model.DiskSpecs `json:"diskSpecs"`
	MemorySpecs agent_model.MemorySpecs  `json:"memorySpecs"`
} // @Name AgentSpecs
