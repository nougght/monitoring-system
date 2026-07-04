package collector

import (
	"agent/internal/config"
	"agent/internal/model"
	"agent/internal/utils"
	"context"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
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
	memoryDetails, err := GetMemoryDetails()
	if err != nil {
		return nil, err
	}
	log.Printf("memory details: %v\n", memoryDetails)
	physicalMemoryList := make([]model.PhysicalMemoryInfo, len(memoryDetails))
	for i, m := range memoryDetails {
		physicalMemoryList[i].BankLabel = m.BankLabel
		physicalMemoryList[i].Capacity = m.Capacity
		physicalMemoryList[i].ConfiguredClockSpeed = m.ConfiguredClockSpeed
		physicalMemoryList[i].DeviceLocator = m.DeviceLocator
		physicalMemoryList[i].FormFactor = utils.ConvertWinPhysicalFormFactor(m.FormFactor)
		physicalMemoryList[i].HotSwappable = m.HotSwappable
		physicalMemoryList[i].Manufacturer = m.Manufacturer
		physicalMemoryList[i].MemoryType = utils.ConvertWinPhysicalMemoryType(m.SMBIOSMemoryType)
		physicalMemoryList[i].ModelName = m.PartNumber
		physicalMemoryList[i].Removable = m.Removable
		physicalMemoryList[i].Replaceable = m.Replaceable
		physicalMemoryList[i].SerialNumber = m.SerialNumber
	}
	cpuDetails, err := GetProcessorDetails()
	if err != nil {
		return nil, err
	}
	log.Printf("processor details: %v\n", cpuDetails[0])

	osArch := runtime.GOARCH
	if osArch == "amd64" {
		osArch = "x64"
	} else if osArch == "386" {
		osArch = "x86"
	}
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
		ModelName:                     cpuInfo[0].ModelName,
		Architecture:                  utils.ConvertWinCpuArch(cpuDetails[0].Architecture),
		Availability:                  utils.ConvertWinCpuAvailability(cpuDetails[0].Availability),
		CurrentClockSpeed:             cpuDetails[0].CurrentClockSpeed,
		DataWidth:                     cpuDetails[0].DataWidth,
		L2CacheSize:                   cpuDetails[0].L2CacheSize,
		L3CacheSize:                   cpuDetails[0].L3CacheSize,
		Manufacturer:                  cpuDetails[0].Manufacturer,
		MaxClockSpeed:                 cpuDetails[0].MaxClockSpeed,
		NumberOfCores:                 uint32(coreCount),
		NumberOfEnabledCore:           uint32(logicalCoreCount),
		NumberOfLogicalProcessors:     uint32(logicalCoreCount),
		ProcessorId:                   cpuDetails[0].ProcessorId,
		SocketDesignation:             cpuDetails[0].SocketDesignation,
		Stepping:                      cpuDetails[0].Stepping,
		VirtualizationFirmwareEnabled: cpuDetails[0].VirtualizationFirmwareEnabled,
	}

	memorySpecs := model.MemorySpecs{
		Total:              virtual.Total,
		PhysicalMemoryList: physicalMemoryList,
	}
	return &model.Specs{
		Host:   hostSpecs,
		CPU:    cpuSpecs,
		Disk:   diskSpecsList,
		Memory: memorySpecs,
	}, nil
}
