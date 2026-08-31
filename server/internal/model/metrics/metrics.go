package metrics_model

import (
	"time"

	"github.com/google/uuid"
)

type MetricSample struct {
	Kind      int32
	Label     string
	Value     float64
	Timestamp time.Time
}

type MetricsBatch struct {
	ID      uint64
	AgentID uuid.UUID
	Metrics []*MetricSample
}

// metric row for db
type MetricRow struct {
	Timestamp time.Time `db:"time"`
	SeriesID  int32     `db:"series_id"`
	Value     float64   `db:"value"`
}
