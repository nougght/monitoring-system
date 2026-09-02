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
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	// begin real or pseudo transaction
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Repositories struct {
	agentRpository           *AgentRepository
	enrollmentKeysRepository *EnrollmentKeysRepository
	specsRepository          *SpecsRepository
	metricsRepository        *MetricsRepository
	seriesRepository         *SeriesRepository
}

func New(db DB) *Repositories {
	return &Repositories{
		agentRpository:           NewAgentRepository(db),
		enrollmentKeysRepository: NewEnrollmentKeysRepository(db),
		specsRepository:          NewSpecsRepository(db),
		metricsRepository:        NewMetricsRepository(db),
		seriesRepository:         NewSeriesRepository(db),
	}
}

func (r *Repositories) AgentRepository() *AgentRepository {
	return r.agentRpository
}

func (r *Repositories) EnrollmentKeysRepository() *EnrollmentKeysRepository {
	return r.enrollmentKeysRepository
}

func (r *Repositories) SpecsRepository() *SpecsRepository {
	return r.specsRepository
}

func (r *Repositories) MetricsRepository() *MetricsRepository {
	return r.metricsRepository
}

func (r *Repositories) SeriesRepository() *SeriesRepository {
	return r.seriesRepository
}
