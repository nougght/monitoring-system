package metrics

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nougght/monitoring-system/server/internal/config"
	"github.com/nougght/monitoring-system/server/internal/model"
	metrics_model "github.com/nougght/monitoring-system/server/internal/model/metrics"
	"github.com/nougght/monitoring-system/server/internal/storage/timescale/repository"
	"github.com/nougght/monitoring-system/server/internal/util"
)

type MetricsService struct {
	cfg         *config.Config
	transactor  model.Transactor
	metricsRepo *repository.MetricsRepository
	snapshot    *SnapshotCache
	batcher     *util.Batcher[metrics_model.MetricsBatch]
}

func NewMetricsService(cfg *config.Config,
	transactor model.Transactor,
	metricsRepo *repository.MetricsRepository,
) (*MetricsService, error) {
	if cfg == nil || transactor == nil || metricsRepo == nil {
		return nil, fmt.Errorf("params required")
	}
	return &MetricsService{
		cfg:         cfg,
		transactor:  transactor,
		metricsRepo: metricsRepo,
	}, nil
}

func (s *MetricsService) HandleMetrics(ctx context.Context, agentID uuid.UUID, metrics metrics_model.MetricsBatch) error {
	s.snapshot.UpdateMetrics(agentID, metrics)

	if err := s.batcher.Add(metrics); err != nil {
		return err
	}

	return nil
}
