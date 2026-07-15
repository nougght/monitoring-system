package collector

import (
	"agent/internal/config"
	"agent/internal/model"
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

type MetricsConsumer interface {
	UpdateMetric(metric model.Metric)
	UpdateSpecs(specs *model.Specs)
}

// collects data from the system
type CollectorService struct {
	config           *config.Config
	metricsConsumer  MetricsConsumer
	metricsChan      chan model.Metric
	wg               sync.WaitGroup
	cancel           context.CancelFunc
	lastNetInput     atomic.Uint64
	lastNetOutput    atomic.Uint64
	lastNetTimestamp atomic.Uint64 //nolint

	processes map[int32]*process.Process
}

func NewCollectorService(cfg *config.Config) *CollectorService {
	return &CollectorService{
		config:          cfg,
		metricsConsumer: nil,
		metricsChan:     make(chan model.Metric),
		wg:              sync.WaitGroup{},
		processes:       make(map[int32]*process.Process),
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
	c.runMemoryCollector(collectorsCtx)
	c.runDiskUsageCollector(collectorsCtx)
	c.runNetIOCollectior(collectorsCtx)
	c.runProcessCollector(collectorsCtx)

	specs, err := c.GetSpecifications(collectorsCtx)
	if err != nil {
		log.Println("failed to get specifications: ", err.Error())
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
				c.metricsChan <- model.NewCpuPercentMetric(float32(percent[0]))
			}
		}
	}()
}

func (c *CollectorService) runMemoryCollector(ctx context.Context) {
	ticker := time.NewTicker(c.config.MemoryInterval)
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				memory, err := mem.VirtualMemory()
				if err != nil {
					log.Printf("failed to get memory: %s", err.Error())
					continue
				}
				// log.Printf("memory: %v", memory)
				c.metricsChan <- model.NewMemoryMetric(memory.Used)
			}

		}
	}()
}

func (c *CollectorService) getDiskUsageMap(ctx context.Context) (map[string]uint64, error) {
	disks, err := disk.PartitionsWithContext(ctx, true)
	if err != nil {
		return nil, err
	}
	diskUsageMap := make(map[string]uint64, len(disks))
	for _, d := range disks {
		diskUsage, err := disk.UsageWithContext(ctx, d.Mountpoint)
		if err != nil {
			return nil, err
		}
		diskUsageMap[diskUsage.Path] = diskUsage.Used
	}
	return diskUsageMap, nil
}

func (c *CollectorService) runDiskUsageCollector(ctx context.Context) {
	ticker := time.NewTicker(c.config.DiskInterval)
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		usage, err := c.getDiskUsageMap(ctx)
		if err == nil {
			c.metricsChan <- model.NewDiskMetric(usage)
		}
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				diskUsageMap, err := c.getDiskUsageMap(ctx)
				if err != nil {
					log.Printf("failed to get disk usage: %s", err.Error())
					continue
				}
				log.Printf("disk usage: %v", diskUsageMap)
				c.metricsChan <- model.NewDiskMetric(diskUsageMap)
			}
		}
	}()
}

func (c *CollectorService) runNetIOCollectior(ctx context.Context) {
	ticker := time.NewTicker(c.config.NetInterval)
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				ioCounters, err := net.IOCounters(false)
				if err != nil {
					log.Printf("failed to get IOCounters: %s", err.Error())
					continue
				}
				if len(ioCounters) == 0 {
					log.Printf("empty IOCounters")
					continue
				}
				lastInput := c.lastNetInput.Load()
				lastOutput := c.lastNetOutput.Load()
				if lastInput != 0 && lastOutput != 0 {
					uploadMbps := float64(ioCounters[0].BytesSent-c.lastNetOutput.Load()) / c.config.NetInterval.Seconds() / 125000
					downloadMbps := float64(ioCounters[0].BytesRecv-c.lastNetInput.Load()) / c.config.NetInterval.Seconds() / 125000
					c.metricsChan <- model.NewNetIOMetric(float32(uploadMbps), float32(downloadMbps))
				}
				c.lastNetInput.Store(ioCounters[0].BytesRecv)
				c.lastNetOutput.Store(ioCounters[0].BytesSent)
			}
		}
	}()
}

func (c *CollectorService) runProcessCollector(ctx context.Context) {
	ticker := time.NewTicker(c.config.ProcessInterval)
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		ps, err := process.ProcessesWithContext(ctx)
		if err != nil {
			log.Printf("failed to get processes: %s", err.Error())
		}
		for _, p := range ps {
			c.processes[p.Pid] = p
		}
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				ps, err := process.ProcessesWithContext(ctx)
				if err != nil {
					log.Printf("failed to get processes: %s", err.Error())
				}
				processList := make([]model.Process, len(ps))
				for i, p := range ps {
					if _, ok := c.processes[p.Pid]; !ok {
						c.processes[p.Pid] = p
					}
					p = c.processes[p.Pid]
					processList[i].Pid = p.Pid
					name, err := p.NameWithContext(ctx)
					if err != nil {
						// log.Printf("failed to get process name: %s", err.Error())
						name = "NO DATA"
					}
					processList[i].Name = name
					// isBackground, _ := p.Background()
					// fmt.Printf("background: %v", isBackground)
					parent, err := p.ParentWithContext(ctx)
					if err != nil {
						// log.Printf("failed to get process name: %s", err.Error())
					}
					if parent == nil {
						processList[i].ParentPid = nil
					} else {
						processList[i].ParentPid = &parent.Pid
					}

					cpuPercent, err := p.PercentWithContext(ctx, 0)
					if err != nil {
						// log.Printf("failed to get cpu percent: %s", err.Error())
						processList[i].CPUPercent = nil
					} else {
						cpuPercent /= 16
						processList[i].CPUPercent = &cpuPercent
					}
					memory, err := p.MemoryInfoWithContext(ctx)
					if err != nil {
						// log.Printf("failed to get memory info: %s", err.Error())
						processList[i].MemoryUsed = nil
					} else {
						processList[i].MemoryUsed = &memory.RSS
					}
					// fmt.Print(p.Username())
				}
				c.metricsChan <- model.NewProcessMetric(processList)
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
	log.Printf("cpu info: %v", cpuInfo)

	disks, err := disk.PartitionsWithContext(ctx, true)
	if err != nil {
		return nil, err
	}
	log.Printf("disks: %v", disks)

	diskSpecsList := make([]model.DiskSpecs, len(disks))
	for i, d := range disks {
		diskSpecsList[i].Device = d.Mountpoint
		diskSpecsList[i].FsType = d.Fstype
		diskUsage, err := disk.UsageWithContext(ctx, d.Mountpoint)
		if err != nil {
			return nil, err
		}
		log.Printf("disk usage: %v", diskUsage)
		diskSpecsList[i].Total = diskUsage.Total
	}
	// userInfo, err := host.Users()
	// if err != nil {
	// 	return nil, err
	// }
	// log.Printf("user info: %v", userInfo)
	virtual, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	log.Printf("virtual memory: %v", virtual)
	physicalMemoryList, err := GetMemoryDetails()
	if err != nil {
		return nil, err
	}
	log.Printf("physical memory: %v", physicalMemoryList)
	// swap, err := mem.SwapMemory()
	// if err != nil {
	// 	return nil, err
	// }
	// log.Printf("swap memory: %v", swap)
	// swapDevices, err := mem.SwapDevices()
	// if err != nil {
	// 	return nil, err
	// }
	// log.Printf("swap devices: %v", swapDevices)
	// c, err := disk.IOCounters()
	// if err != nil {
	// 	return nil, err
	// }
	// log.Printf("disk io: %v", c)
	// label, err := disk.Label(disks[0].Device)
	// if err != nil {
	// 	return nil, err
	// }
	// log.Printf("disk label: %v", label)
	netIO, err := net.IOCounters(false)
	if err != nil {
		return nil, err
	}
	log.Printf("net io: %v", netIO)
	// ps, err := process.ProcessesWithContext(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	// log.Printf("processes: %v", ps)
	// coreCount, err := cpu.CountsWithContext(ctx, false)
	// if err != nil {
	// 	return nil, err
	// }
	// logicalCoreCount, err := cpu.CountsWithContext(ctx, true)
	// if err != nil {
	// 	return nil, err
	// }

	hostSpecs := model.HostSpecs{
		Hostname:        hostInfo.Hostname,
		OsType:          hostInfo.OS,
		Os:              hostInfo.Platform,
		OsVersion:       hostInfo.PlatformVersion,
		OsKernelVersion: hostInfo.KernelVersion,
		OsArch:          hostInfo.KernelArch,
	}

	cpuSpecs, err := GetProcessorDetails()
	if err != nil {
		return nil, err
	}

	memorySpecs := model.MemorySpecs{
		Total:              virtual.Total,
		PhysicalMemoryList: physicalMemoryList,
	}
	return &model.Specs{
		Host:   hostSpecs,
		CPU:    cpuSpecs[0],
		Disk:   diskSpecsList,
		Memory: memorySpecs,
	}, nil
}
