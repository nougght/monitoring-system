package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nougght/monitoring-system/server/internal/model/agent"
)

type AgentRepository struct {
	db *pgxpool.Pool
}

func NewAgentRepository(db *pgxpool.Pool) *AgentRepository {
	return &AgentRepository{
		db: db,
	}
}

func (r *AgentRepository) CreateAgent(ctx context.Context, agent *agent.Agent) (*agent.Agent, error) {
	query := `
	INSERT INTO agents (name, description) 
	VALUES($1, $2)
	`
	err := r.db.QueryRow(ctx, query, agent.Name, agent.Description).Scan(&agent.ID, &agent.CreatedAt, &agent.LastSeenAt)
	if err != nil {
		return nil, fmt.Errorf("insert failed: %w", err)
	}

	return agent, nil
}
