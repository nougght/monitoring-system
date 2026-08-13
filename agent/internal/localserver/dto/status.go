package dto

import "github.com/google/uuid"

type Status struct {
	AgentID          uuid.UUID `json:"agentId"`
	ConnectionStatus string    `json:"connectionStatus"`
	CertInfo         string    `json:"CertInfo"`
	ServerInfo       string    `json:"ServerInfo"`
}
