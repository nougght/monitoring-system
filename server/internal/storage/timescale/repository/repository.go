package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DB interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Repositories struct {
	agentRpository           *AgentRepository
	enrollmentKeysRepository *EnrollmentKeysRepository
}

func New(db DB) *Repositories {
	return &Repositories{
		agentRpository:           NewAgentRepository(db),
		enrollmentKeysRepository: NewEnrollmentKeysRepository(db),
	}
}

func (r *Repositories) AgentRepository() *AgentRepository {
	return r.agentRpository
}

func (r *Repositories) EnrollmentKeysRepository() *EnrollmentKeysRepository {
	return r.enrollmentKeysRepository
}
