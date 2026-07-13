package metrics

import (
	"agent/internal/config"
	"agent/internal/model"
	"context"
	"log"
	"sync"
	"time"
)

type MetricsService struct {
	config *config.Config

	// temporary storing
	specs     *model.Specs
	specsMu   sync.RWMutex
	metrics   *model.Metrics
	metricsMu sync.RWMutex

	refreshSpecsFunc func(ctx context.Context) (*model.Specs, error)
}

func NewMetricsService(cfg *config.Config,
	refreshSpecsFunc func(ctx context.Context) (*model.Specs, error)) *MetricsService {
	return &MetricsService{
		config:           cfg,
		refreshSpecsFunc: refreshSpecsFunc,
		metrics:          new(model.Metrics),
	}
}

func (m *MetricsService) UpdateMetric(metric model.Metric) {
	// log.Println("update metric", metric.Type())
	switch metric.Type() {
	case model.MetricTypeFocusedWindow:
		m.metrics.FocusedWindow = metric.(*model.FocusedWindowMetric)
	case model.MetricTypeCpuPercent:
		m.metrics.CpuPercent = metric.(*model.CpuPercentMetric)
	case model.MetricTypeMemory:
		m.metrics.MemoryUsage = metric.(*model.MemoryMetric)
	case model.MetricTypeDisk:
		m.metrics.DiskUsage = metric.(*model.DiskMetric)
	case model.MetricTypeNet:
		netMetric := metric.(*model.NetIOMetric)
		m.metrics.NetworkUsage = netMetric
	case model.MetricTypeProcess:
		m.metrics.Process = metric.(*model.ProcessMetric)
	}

}

func (m *MetricsService) GetSpecs(ctx context.Context) (*model.Specs, error) {
	m.metricsMu.RLock()
	if m.specs == nil {
		m.metricsMu.RUnlock()
		specs, err := m.refreshSpecsFunc(ctx)
		if err != nil {
			return nil, err
		}
		m.metricsMu.Lock()
		m.specs = specs
		m.metricsMu.Unlock()
	}
	m.metricsMu.RUnlock()
	return m.specs, nil
}

func (m *MetricsService) UpdateSpecs(specs *model.Specs) {
	m.specsMu.Lock()
	defer m.specsMu.Unlock()
	log.Println("update specs", specs)
	m.specs = specs
}

func (m *MetricsService) GetFocusedWindow() *model.FocusedWindowMetric {
	m.metricsMu.RLock()
	defer m.metricsMu.RUnlock()
	if m.metrics.FocusedWindow == nil {
		return nil
	}
	return m.metrics.FocusedWindow
}

func (m *MetricsService) GetCpuPercent() *model.CpuPercentMetric {
	m.metricsMu.RLock()
	defer m.metricsMu.RUnlock()
	if m.metrics.CpuPercent == nil {
		return nil
	}
	return m.metrics.CpuPercent
}

func (m *MetricsService) GetMemory() *model.MemoryMetric {
	m.metricsMu.RLock()
	defer m.metricsMu.RUnlock()
	if m.metrics.MemoryUsage == nil {
		return nil
	}
	return m.metrics.MemoryUsage
}

func (m *MetricsService) GetDiskUsage() *model.DiskMetric {
	m.metricsMu.RLock()
	defer m.metricsMu.RUnlock()
	if m.metrics.DiskUsage == nil {
		return nil
	}
	return m.metrics.DiskUsage
}

func (m *MetricsService) GetNetworkUsage() *model.NetIOMetric {
	m.metricsMu.RLock()
	defer m.metricsMu.RUnlock()
	if m.metrics.NetworkUsage == nil {
		return nil
	}
	return m.metrics.NetworkUsage
}

func (m *MetricsService) GetProcess() *model.ProcessMetric {
	m.metricsMu.RLock()
	defer m.metricsMu.RUnlock()
	if m.metrics.Process == nil {
		return nil
	}
	return m.metrics.Process
}

func (m *MetricsService) GetTimestamp() time.Time {
	m.metricsMu.RLock()
	defer m.metricsMu.RUnlock()
	return m.metrics.Timestamp
}

func (m *MetricsService) GetMetrics() *model.Metrics {
	return &model.Metrics{
		FocusedWindow: m.GetFocusedWindow(),
		CpuPercent:    m.GetCpuPercent(),
		MemoryUsage:   m.GetMemory(),
		DiskUsage:     m.GetDiskUsage(),
		NetworkUsage:  m.GetNetworkUsage(),
		Process:       m.GetProcess(),
		Timestamp:     time.Now(),
	}
}
