package model

import "time"

type Metrics struct {
	FocusedWindow *FocusedWindowMetric `json:"focusedWindow,omitempty"`
	CpuPercent    *CpuPercentMetric    `json:"cpuPercent,omitempty"`
	MemoryUsage   *MemoryMetric        `json:"memoryUsage,omitempty"`
	DiskUsage     *DiskMetric          `json:"diskUsage,omitempty"`
	NetworkUsage  *NetIOMetric         `json:"networkUsage,omitempty"`
	Process       *ProcessMetric       `json:"process,omitempty"`
	Timestamp     time.Time            `json:"timestamp"`
} // @name Metrics
