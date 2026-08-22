package model

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type AgentState struct {
	mu              sync.RWMutex
	agentID         uuid.UUID
	connected       bool
	pending         bool
	lastConnectedAt *time.Time
}

func NewAgentState(agentID uuid.UUID) *AgentState {
	return &AgentState{
		agentID: agentID,
	}
}

func (s *AgentState) AgentID() uuid.UUID {
	return s.agentID
}

func (s *AgentState) Connected() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.connected
}

func (s *AgentState) Pending() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.pending
}

func (s *AgentState) LastConnectedAt() *time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastConnectedAt
}

func (s *AgentState) SetConnected(isConnected bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connected = isConnected
	if isConnected && s.pending {
		s.pending = false
	}
	log.Printf("connected: %b, pending: %b", s.connected, s.pending)
}
func (s *AgentState) SetPending(isPending bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pending = isPending
	if isPending && s.connected {
		s.connected = false
	}
	log.Printf("connected: %b, pending: %b", s.connected, s.pending)
}

func (s *AgentState) SetLastConnectedAt(t time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastConnectedAt = &t
}
