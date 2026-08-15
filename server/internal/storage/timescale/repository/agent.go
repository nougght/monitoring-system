package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nougght/monitoring-system/server/internal/model"
	agent_model "github.com/nougght/monitoring-system/server/internal/model/agent"
)

type AgentRepository struct {
	pool DB
}

func NewAgentRepository(db DB) *AgentRepository {
	return &AgentRepository{
		pool: db,
	}

}

func (r *AgentRepository) db(ctx context.Context) DB {
	res := r.pool
	if tx := ctx.Value(model.ContextKeyTx); tx != nil {
		res = tx.(DB)
	}
	return res
}

func (r *AgentRepository) CreateAgent(ctx context.Context, agent *agent_model.Agent) (*agent_model.Agent, error) {
	query := `
	INSERT INTO agents (name, description) 
	VALUES($1, $2)
	RETURNING id, created_at, last_seen_at
	`
	err := r.db(ctx).QueryRow(ctx, query, agent.Name, agent.Description).Scan(&agent.ID, &agent.CreatedAt, &agent.LastSeenAt)
	if err != nil {
		return nil, fmt.Errorf("insert failed: %w", err)
	}

	return agent, nil
}

func (r *AgentRepository) GetAllAgents(ctx context.Context) (res []*agent_model.Agent, err error) {
	query := `
		SELECT * FROM agents;
		`
	// SELECT *,  EXISTS (
	//     SELECT 1
	//     FROM enrollment_keys keys
	//     WHERE keys.agent_id = agent.id AND keys.used_at != NULL
	// ) AS is_enrolled FROM agents
	// `
	rows, err := r.db(ctx).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("select failed: %w", err)
	}
	defer rows.Close()

	res, err = pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[agent_model.Agent])
	if err != nil {
		return nil, fmt.Errorf("collect rows failed: %w", err)
	}

	return res, nil
}

func (r *AgentRepository) UpdateStatus(ctx context.Context, agentID uuid.UUID, status agent_model.AgentStatus) error {
	query := `
	UPDATE agents SET status = $1 WHERE ID = $2
	`
	_, err := r.db(ctx).Exec(ctx, query, status, agentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("insert failed: %w", ErrNotFound)
		}
		return fmt.Errorf("insert failed: %w", err)
	}

	return nil
}
