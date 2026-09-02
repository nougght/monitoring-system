package metrics

import (
	"context"
	"fmt"
	"log"
	"slices"
	"sync"

	metrics_model "github.com/nougght/monitoring-system/server/internal/model/metrics"
	metrickinds "github.com/nougght/monitoring-system/shared/go/metric_kinds"
)

type SeriesRepository interface {
	BatchCreateOrLoadSeries(ctx context.Context, seriesList []metrics_model.MetricSeriesKey) ([]metrics_model.MetricSeries, error)
	GetAllSeries(ctx context.Context) ([]*metrics_model.MetricSeries, error)
	// GetSeriesIDs(ctx context.Context, seriesKeys []metrics_model.MetricSeriesKey) (map[metrics_model.MetricSeriesKey]uuid.UUID, error)
}

type SeriesResolver struct {
	mu    sync.RWMutex
	cache map[metrics_model.MetricSeriesKey]metrics_model.MetricSeriesID
	repo  SeriesRepository
}

func NewSeriesResolver(seriesRepo SeriesRepository) *SeriesResolver {
	return &SeriesResolver{
		cache: make(map[metrics_model.MetricSeriesKey]metrics_model.MetricSeriesID),
		repo:  seriesRepo,
	}
}

func (c *SeriesResolver) LoadAllSeries(ctx context.Context) error {
	series, err := c.repo.GetAllSeries(ctx)
	if err != nil {
		return fmt.Errorf("failed to get all series: %w", err)
	}
	c.mu.Lock()
	for _, s := range series {
		c.cache[s.MetricSeriesKey] = s.ID
	}

	c.mu.Unlock()
	return nil
}

func (c *SeriesResolver) ResolveAllSeries(ctx context.Context, keys []metrics_model.MetricSeriesKey) (resultIDs map[metrics_model.MetricSeriesKey]metrics_model.MetricSeriesID, err error) {
	log.Printf("resolving series %#v", keys)
	// filter unknown metric kinds
	keys = slices.DeleteFunc(keys, func(key metrics_model.MetricSeriesKey) bool {
		_, ok := metrickinds.MetricKinds[key.Kind]
		if !ok {
			log.Printf("resolving unknown metric kind: %d", key.Kind)
			return true
		}
		return false
	})

	missingKeys := c.getMissingKeys(ctx, keys)
	if len(missingKeys) > 0 {
		err = c.addOrLoadSeries(ctx, missingKeys)
		if err != nil {
			return nil, fmt.Errorf("failed to add or load series: %w", err)
		}
	}
	resultIDs, err = c.getSeriesIDsByKeys(ctx, keys)
	if err != nil {
		return nil, fmt.Errorf("failed to get series ids by keys after loading misses: %w", err)
	}

	return resultIDs, nil
}

// returns list of keys, that are missing in cache
func (c *SeriesResolver) getMissingKeys(ctx context.Context, keys []metrics_model.MetricSeriesKey) (missingKeys []metrics_model.MetricSeriesKey) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, key := range keys {
		_, ok := c.cache[key]
		if !ok {
			missingKeys = append(missingKeys, key)
		}
	}
	return missingKeys
}

func (c *SeriesResolver) getSeriesID(ctx context.Context, key metrics_model.MetricSeriesKey) (id metrics_model.MetricSeriesID, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	id, ok = c.cache[key]
	return
}

func (c *SeriesResolver) getSeriesIDsByKeys(ctx context.Context, keys []metrics_model.MetricSeriesKey) (map[metrics_model.MetricSeriesKey]metrics_model.MetricSeriesID, error) {
	res := make(map[metrics_model.MetricSeriesKey]metrics_model.MetricSeriesID, len(keys))
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, key := range keys {
		id, ok := c.cache[key]
		if !ok {
			return nil, fmt.Errorf("not found series key: %#v", key)
		}
		res[key] = id
	}
	return res, nil
}

func (c *SeriesResolver) addOrLoadSeries(ctx context.Context, seriesKeys []metrics_model.MetricSeriesKey) error {
	series, err := c.repo.BatchCreateOrLoadSeries(ctx, seriesKeys)
	if err != nil {
		return fmt.Errorf("create or load series failed: %w", err)
	}
	c.mu.Lock()
	log.Println("addOrLoad count: %d", len(series))
	for _, s := range series {
		c.cache[s.MetricSeriesKey] = s.ID
	}
	c.mu.Unlock()
	return nil
}
