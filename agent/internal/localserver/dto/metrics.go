package dto

import "time"

type Process struct {
	PID        int32    `json:"pid"`
	Name       string   `json:"name"`
	ParentPid  *int32   `json:"parentPid,omitempty"`
	CPUPercent *float64 `json:"cpuPercent,omitempty"`
	MemoryUsed *uint64  `json:"memoryUsed,omitempty"`
} // @name Process
type Metrics struct {
	FocusedWindow *string            `json:"focusedWindow,omitempty"`
	ProcessList   []Process          `json:"processList,omitempty"`
	CpuPercent    *float32           `json:"cpuPercent,omitempty"`
	MemoryUsed    *uint64            `json:"memoryUsed,omitempty"`
	DiskUsage     *map[string]uint64 `json:"diskUsage,omitempty"`
	UploadMbps    *float32           `json:"uploadMbps,omitempty"`
	DownloadMbps  *float32           `json:"downloadMbps,omitempty"`
	Timestamp     time.Time          `json:"timestamp"`
} // @name Metrics
