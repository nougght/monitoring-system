package metrics

import (
	"context"
	"fmt"
	"log"
	"maps"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/nougght/monitoring-system/server/internal/config"
	"github.com/nougght/monitoring-system/server/internal/model"
	metrics_model "github.com/nougght/monitoring-system/server/internal/model/metrics"
	"github.com/nougght/monitoring-system/server/internal/storage/timescale/repository"
	"github.com/nougght/monitoring-system/server/internal/util"

	metrickinds "github.com/nougght/monitoring-system/shared/go/metric_kinds"
	sharedUtil "github.com/nougght/monitoring-system/shared/go/util"
)

type MetricsService struct {
	cfg            *config.Config
	transactor     model.Transactor
	metricsRepo    *repository.MetricsRepository
	snapshot       *SnapshotCache
	batcher        *util.Batcher[metrics_model.MetricsBatch]
	seriesResolver *SeriesResolver
}

func NewMetricsService(cfg *config.Config,
	transactor model.Transactor,
	metricsRepo *repository.MetricsRepository,
	seriesRepo *repository.SeriesRepository,
) (*MetricsService, error) {
	if cfg == nil || transactor == nil || metricsRepo == nil {
		return nil, fmt.Errorf("params required")
	}
	s := &MetricsService{
		cfg:            cfg,
		transactor:     transactor,
		metricsRepo:    metricsRepo,
		snapshot:       NewSnapshotCache(),
		seriesResolver: NewSeriesResolver(seriesRepo),
	}
	s.batcher = util.NewBatcher(1000, time.Second, s.resolveAndSaveBatchFunc)
	return s, nil
}

// start metrics saving
func (s *MetricsService) StartSaving(ctx context.Context) {
	s.batcher.Run(ctx)
}

func (s *MetricsService) SyncMetricKinds(ctx context.Context) error {
	err := s.metricsRepo.UpsertMetricKinds(ctx, metrickinds.MetricKindList())
	if err != nil {
		return fmt.Errorf("failed to upsert metric kinds: %w", err)
	}
	return nil
}

func (s *MetricsService) HandleMetrics(ctx context.Context, agentID uuid.UUID, metrics metrics_model.MetricsBatch) error {
	s.snapshot.UpdateMetrics(agentID, metrics)

	if err := s.batcher.Add(metrics); err != nil {
		return err
	}

	return nil
}

func (s *MetricsService) resolveAndSaveBatchFunc(ctx context.Context, batches []metrics_model.MetricsBatch) error {
	metricsCount := metrics_model.CountAllMetricsInBatchList(batches)
	seriesKeys := make(map[metrics_model.MetricSeriesKey]struct{}, metricsCount)

	for _, b := range batches {
		for _, key := range sharedUtil.Map(b.Metrics, func(m *metrics_model.MetricSample) metrics_model.MetricSeriesKey {
			return metrics_model.MetricSeriesKey{
				AgentID: b.AgentID,
				Kind:    m.Kind,
				Label:   m.Label,
			}
		}) {
			seriesKeys[key] = struct{}{}
		}
	}

	seriesIDs, err := s.seriesResolver.ResolveAllSeries(ctx, slices.Collect(maps.Keys(seriesKeys)))
	if err != nil {
		return fmt.Errorf("failed to resolve metric series: %w", err)
	}
	metricRows := make([]metrics_model.MetricRow, metricsCount)
	for _, b := range batches {
		for i, key := range sharedUtil.Map(b.Metrics, func(m *metrics_model.MetricSample) metrics_model.MetricSeriesKey {
			return metrics_model.MetricSeriesKey{
				AgentID: b.AgentID,
				Kind:    m.Kind,
				Label:   m.Label,
			}
		}) {
			id, ok := seriesIDs[key]
			if !ok {
				log.Printf("series id not found after resolving: %v", key)
				continue
			}
			metricRows = append(metricRows, metrics_model.MetricRow{
				SeriesID:  int64(id),
				Value:     b.Metrics[i].Value,
				Timestamp: b.Metrics[i].Timestamp,
			})
		}
	}

	err = s.metricsRepo.SaveRows(ctx, metricRows)
	if err != nil {
		return fmt.Errorf("failed to save metric rows: %w", err)
	}
	log.Println("metrics saved")
	return nil
}
