package metrics

import (
	"agent/internal/config"
	"agent/internal/model"
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type MetricsService struct {
	config      *config.Config
	metricsChan <-chan model.Metric

	// temporary storing
	specs         *model.Specs
	specsMu       sync.Mutex
	focusedWindow atomic.Value
	cpuPercent    atomic.Value

	refreshSpecsFunc func(ctx context.Context) (*model.Specs, error)
}

func NewMetricsService(cfg *config.Config,
	refreshSpecsFunc func(ctx context.Context) (*model.Specs, error)) *MetricsService {
	return &MetricsService{
		config:           cfg,
		refreshSpecsFunc: refreshSpecsFunc,
	}
}

func (m *MetricsService) UpdateMetric(metric model.Metric) {
	// log.Println("update metric", metric.Type())
	switch metric.Type() {
	case model.MetricTypeFocusedWindow:
		m.focusedWindow.Store(metric.(*model.FocusedWindowMetric).Value())
	case model.MetricTypeCpuPercent:
		m.cpuPercent.Store(metric.(*model.CpuPercentMetric).Value())
	}
}

func (m *MetricsService) GetSpecs(ctx context.Context) (*model.Specs, error) {
	m.specsMu.Lock()
	defer m.specsMu.Unlock()
	if m.specs == nil {
		specs, err := m.refreshSpecsFunc(ctx)
		if err != nil {
			return nil, err
		}
		m.specs = specs
	}
	return m.specs, nil
}

func (m *MetricsService) UpdateSpecs(specs *model.Specs) {
	m.specsMu.Lock()
	defer m.specsMu.Unlock()
	log.Println("update specs", specs)
	m.specs = specs
}

func (m *MetricsService) GetFocusedWindow() *string {
	val := m.focusedWindow.Load()
	if val == nil {
		return nil
	}
	title := val.(string)
	return &title
}

func (m *MetricsService) GetCpuPercent() *float64 {
	val := m.cpuPercent.Load()
	if val == nil {
		return nil
	}
	percent := val.(float64)
	return &percent
}

func (m *MetricsService) GetMetrics() *model.Metrics {
	return &model.Metrics{
		FocusedWindow: m.GetFocusedWindow(),
		CpuPercent:    m.GetCpuPercent(),
		Timestamp:     time.Now(),
	}
}
