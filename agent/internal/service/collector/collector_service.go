package collector

import (
	"agent/internal/config"
	"agent/internal/model"
	"context"
	"log"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

type MetricsConsumer interface {
	UpdateMetric(metric model.Metric)
	UpdateSpecs(specs *model.CpuSpecs)
}

// collects data from the system
type CollectorService struct {
	config          *config.Config
	metricsConsumer MetricsConsumer
	metricsChan     chan model.Metric
	wg              sync.WaitGroup
}

func NewCollectorService(cfg *config.Config) *CollectorService {
	return &CollectorService{
		config:          cfg,
		metricsConsumer: nil,
		metricsChan:     make(chan model.Metric),
		wg:              sync.WaitGroup{},
	}
}

func (c *CollectorService) SetMetricsConsumer(metricsConsumer MetricsConsumer) {
	c.metricsConsumer = metricsConsumer
}

func (c *CollectorService) StartCollectors(ctx context.Context) {
	c.wg.Add(1)
	c.runFocusedWindowCollector(ctx)

	specs, err := c.GetSpecifications(ctx)
	if err != nil {
		log.Println("failed to get specifications")
	}
	c.metricsConsumer.UpdateSpecs(specs)

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case metric := <-c.metricsChan:
				if c.metricsConsumer == nil {
					log.Println("no metrics consumer")
					continue
				}
				c.metricsConsumer.UpdateMetric(metric)
			}
		}
	}()
}

func (c *CollectorService) GetSpecifications(ctx context.Context) (*model.CpuSpecs, error) {
	return getSpecifications(ctx)
}

func (c *CollectorService) runFocusedWindowCollector(ctx context.Context) {
	ticker := time.NewTicker(c.config.FocusedWindowInterval)
	go func() {
		defer c.wg.Done()
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				focusedWindow := GetFocusedWindowTitle()
				c.metricsChan <- model.NewFocusedWindowMetric(focusedWindow)
			}
		}
	}()
}

func getSpecifications(ctx context.Context) (*model.CpuSpecs, error) {
	cpuInfo, err := cpu.InfoWithContext(ctx)
	if err != nil {
		return nil, err
	}
	coreCount, err := cpu.CountsWithContext(ctx, false)
	if err != nil {
		return nil, err
	}
	logicalCoreCount, err := cpu.CountsWithContext(ctx, true)
	if err != nil {
		return nil, err
	}
	cpuSpecs := &model.CpuSpecs{
		ModelName:        cpuInfo[0].ModelName,
		CoreCount:        coreCount,
		LogicalCoreCount: logicalCoreCount,
	}
	return cpuSpecs, nil
}
