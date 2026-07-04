package collector

import (
	"agent/internal/config"
	"agent/internal/model"
	"context"
	"log"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
)

type MetricsConsumer interface {
	UpdateMetric(metric model.Metric)
	UpdateSpecs(specs *model.Specs)
}

// collects data from the system
type CollectorService struct {
	config          *config.Config
	metricsConsumer MetricsConsumer
	metricsChan     chan model.Metric
	wg              sync.WaitGroup
	cancel          context.CancelFunc
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

func (c *CollectorService) GetSpecifications(ctx context.Context) (*model.Specs, error) {
	return getSpecifications(ctx)
}

// starts parallel metrics collectors
func (c *CollectorService) StartCollectors(ctx context.Context) {
	log.Println("starting collectors")
	collectorsCtx, cancel := context.WithCancel(ctx)
	c.cancel = cancel
	c.runMetricsSender(collectorsCtx)

	c.runFocusedWindowCollector(collectorsCtx)
	c.runCpuPercentCollector(collectorsCtx)

	specs, err := c.GetSpecifications(collectorsCtx)
	if err != nil {
		log.Println("failed to get specifications")
	}
	c.metricsConsumer.UpdateSpecs(specs)

}

// stops collectors with waiting
func (c *CollectorService) StopCollectors() {
	c.cancel()
	c.wg.Wait()
}

// sends metrics from channel to consumer
func (c *CollectorService) runMetricsSender(ctx context.Context) {
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

func (c *CollectorService) runFocusedWindowCollector(ctx context.Context) {
	ticker := time.NewTicker(c.config.FocusedWindowInterval)
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				focusedWindow, err := GetFocusedWindowTitle()
				if err != nil {
					log.Printf("failed to get focused window: %s", err.Error())
					continue
				}
				// log.Println("focused window " + focusedWindow)
				c.metricsChan <- model.NewFocusedWindowMetric(focusedWindow)
			}
		}
	}()
}

func (c *CollectorService) runCpuPercentCollector(ctx context.Context) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				percent, err := cpu.PercentWithContext(ctx, c.config.CpuPercentInterval, false)
				if err != nil {
					log.Printf("failed to get cpu percent: %s", err.Error())
					continue
				}
				// log.Println("cpu percent", percent[0])
				c.metricsChan <- model.NewCpuPercentMetric(percent[0])
			}
		}
	}()
}

func getSpecifications(ctx context.Context) (*model.Specs, error) {
	hostInfo, err := host.InfoWithContext(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("host info: %v", *hostInfo)

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
	hostSpecs := model.HostSpecs{
		Hostname:        hostInfo.Hostname,
		OsType:          hostInfo.OS,
		Os:              hostInfo.Platform,
		OsVersion:       hostInfo.PlatformVersion,
		OsKernelVersion: hostInfo.KernelVersion,
		OsArch:          hostInfo.KernelArch,
	}
	cpuSpecs := model.CpuSpecs{
		ModelName:        cpuInfo[0].ModelName,
		CoreCount:        coreCount,
		LogicalCoreCount: logicalCoreCount,
	}
	return &model.Specs{
		Host: hostSpecs,
		CPU:  cpuSpecs,
	}, nil
}
