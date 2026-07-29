package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nougght/monitoring-system/server/internal/model/agent"
)

type EnrollmentKeysRepository struct {
	db *pgxpool.Pool
}

func NewEnrollmentKeysRepository(db *pgxpool.Pool) *EnrollmentKeysRepository {
	return &EnrollmentKeysRepository{
		db: db,
	}
}

func (r *EnrollmentKeysRepository) CreateKey(ctx context.Context, key *agent.EnrollmentKey) (*agent.EnrollmentKey, error) {
	query := `
	INSERT INTO enrollment_keys (key_hash, agent_id, expires_at) 
	VALUES($1, $2, $3)
	`
	err := r.db.QueryRow(ctx, query, key.HashString, key.AgentID, key.ExpiresAt).Scan(&key.ID)
	if err != nil {
		return nil, fmt.Errorf("insert failed: %w", err)
	}

	return key, nil
}
