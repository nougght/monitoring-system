package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nougght/monitoring-system/server/internal/model"
	metrics_model "github.com/nougght/monitoring-system/server/internal/model/metrics"
)

type SeriesRepository struct {
	database DB
}

func NewSeriesRepository(db DB) *SeriesRepository {
	return &SeriesRepository{
		database: db,
	}

}

func (r *SeriesRepository) db(ctx context.Context) DB {
	res := r.database
	if tx := ctx.Value(model.ContextKeyTx); tx != nil {
		res = tx.(DB)
	}
	return res
}
func (r *SeriesRepository) BatchCreateSeries(ctx context.Context, seriesList []metrics_model.MetricSeries) ([]metrics_model.MetricSeries, error) {
	query := `
            INSERT INTO metric_series (agent_id, kind, label) VALUES ($1,$2,$3)
            ON CONFLICT (agent_id, kind, label) DO NOTHING
            RETURNING id
        `

	batch := &pgx.Batch{}
	for _, s := range seriesList {
		batch.Queue(query, s.AgentID, s.Kind, s.Label)
	}

	res := r.db(ctx).SendBatch(ctx, batch)
	defer func() {
		_ = res.Close()
	}()

	for i := range seriesList {
		err := res.QueryRow().Scan(&seriesList[i].AgentID)
		// do nothing if no row
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("batch QueryRow scan error: %w", err)
		}
	}
	return seriesList, nil
}

func (r *SeriesRepository) GetSeriesIDs(ctx context.Context, seriesKeys []metrics_model.SeriesKey) (map[metrics_model.SeriesKey]uuid.UUID, error) {
	query := `
		SELECT s.id, s.agent_id, s.kind, s.label
		FROM metric_series s
		JOIN unnest($1::uuid[], $2::smallint[], $3::text[]) AS u(agent_id, kind, label)
  		ON s.agent_id = u.agent_id AND s.kind = u.kind AND s.label = u.label;`
	agentIDList, kindList, labelList := make([]uuid.UUID, len(seriesKeys)), make([]string, len(seriesKeys)), make([]string, len(seriesKeys))
	rows, err := r.db(ctx).Query(ctx, query, agentIDList, kindList, labelList)
	if err != nil {
		return nil, fmt.Errorf("select failed: %w", err)
	}
	defer func() {
		rows.Close()
	}()

	var id uuid.UUID
	var key metrics_model.SeriesKey
	res := make(map[metrics_model.SeriesKey]uuid.UUID, len(seriesKeys))
	for rows.Next() {
		err = rows.Scan(&id, &key.AgentID, &key.Kind, &key.Label)
		if err != nil {
			return nil, fmt.Errorf("rows scan failed: %w", err)
		}
		res[key] = id
	}
	return res, nil
}

func (r *SeriesRepository) GetAllSeries(ctx context.Context) ([]*metrics_model.MetricSeries, error) {
	query := `SELECT * FROM metric_series`
	rows, err := r.db(ctx).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("select failed: %w", err)
	}
	defer rows.Close()
	series, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[metrics_model.MetricSeries])
	if err != nil {
		return nil, fmt.Errorf("collect rows failed: %w", err)
	}
	return series, nil
}
