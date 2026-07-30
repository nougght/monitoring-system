package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nougght/monitoring-system/server/internal/model"
	"github.com/nougght/monitoring-system/server/internal/model/agent"
)

type EnrollmentKeysRepository struct {
	pool DB
}

func NewEnrollmentKeysRepository(db DB) *EnrollmentKeysRepository {
	return &EnrollmentKeysRepository{
		pool: db,
	}
}

func (r *EnrollmentKeysRepository) db(ctx context.Context) DB {
	res := r.pool
	if tx := ctx.Value(model.TxKey); tx != nil {
		res = tx.(DB)
	}
	return res
}

func (r *EnrollmentKeysRepository) CreateKey(ctx context.Context, key *agent.EnrollmentKey) (*agent.EnrollmentKey, error) {
	query := `
	INSERT INTO enrollment_keys (key_hash, agent_id, expires_at) 
	VALUES($1, $2, $3)
	`
	_, err := r.db(ctx).Exec(ctx, query, key.HashString, key.AgentID, key.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("insert failed: %w", err)
	}

	return key, nil
}

func (r *EnrollmentKeysRepository) GetKeyByAgentId(ctx context.Context, agentID uuid.UUID) (*agent.EnrollmentKey, error) {
	query := `
	SELECT * FROM enrollment_keys k WHERE k.agent_id = $1
	VALUES($1)
	`
	rows, err := r.db(ctx).Query(ctx, query, agentID)
	if err != nil {
		return nil, fmt.Errorf("select failed: %w", err)
	}

	key, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[agent.EnrollmentKey])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("row collect failed: %w", err)
	}
	return &key, nil
}

func (r *EnrollmentKeysRepository) SetUsed(ctx context.Context, agentID uuid.UUID, usedAt time.Time) error {
	query := `
		UPDATE enrollment_keys keys 
		SET used_at = $1
		WHERE agent_id = $2 AND used_at = NULL
	`
	res, err := r.db(ctx).Exec(ctx, query, usedAt, agentID)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	if res.RowsAffected() == 0 {
		return ErrNoAffectedRows
	}

	return nil
}
