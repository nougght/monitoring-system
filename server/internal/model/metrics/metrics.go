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

func CountAllMetricsInBatchList(batchList []MetricsBatch) (count int) {
	for i := range batchList {
		count += len(batchList[i].Metrics)
	}
	return count
}

// metric row for db
type MetricRow struct {
	Timestamp time.Time `db:"time"`
	SeriesID  int64     `db:"series_id"`
	Value     float64   `db:"value"`
}

type MetricSeries struct {
	ID MetricSeriesID `db:"id"`
	MetricSeriesKey
}

type MetricSeriesID int64

type MetricSeriesKey struct {
	AgentID uuid.UUID `db:"agent_id"`
	Kind    int32     `db:"kind"`
	Label   string    `db:"lable"`
}
