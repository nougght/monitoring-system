package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/nougght/monitoring-system/server/internal/model"
	metrics_model "github.com/nougght/monitoring-system/server/internal/model/metrics"
)

type MetricsRepository struct {
	database DB
}

const (
	metricSamplesTableName = "metric_samples"
	metricSamplesSeriesID  = "series_id"
	metricSamplesValue     = "value"
	metricSamplesTimestamp = "time"
)

type rowSource struct {
	rows []metrics_model.MetricRow
	i    int
	cur  [3]any
}

func (s *rowSource) Next() bool {
	if s.i >= len(s.rows) {
		return false
	}
	r := s.rows[s.i]
	s.i++
	s.cur[0], s.cur[1], s.cur[2] = r.Timestamp, r.SeriesID, r.Value
	return true
}
func (s *rowSource) Values() ([]any, error) {
	return s.cur[:], nil
}
func (s *rowSource) Err() error {
	return nil
}

func NewMetricsRepository(db DB) *MetricsRepository {
	return &MetricsRepository{
		database: db,
	}

}

func (r *MetricsRepository) db(ctx context.Context) DB {
	res := r.database
	if tx := ctx.Value(model.ContextKeyTx); tx != nil {
		res = tx.(DB)
	}
	return res
}

func (r *MetricsRepository) SaveRows(ctx context.Context, rows []metrics_model.MetricRow) error {
	_, err := r.db(ctx).CopyFrom(ctx, pgx.Identifier{metricSamplesTableName},
		[]string{metricSamplesSeriesID, metricSamplesValue, metricSamplesTimestamp},
		pgx.CopyFromSource(&rowSource{rows: rows}))

	if err != nil {
		return fmt.Errorf("copy from failed: %w", err)
	}
	return nil
}
