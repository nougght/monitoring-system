package model

import (
	"time"
)

type MetricType string

const (
	MetricTypeFocusedWindow MetricType = "focused_window" // активное окно
	MetricTypeCpuPercent    MetricType = "cpu_percent"    // процент использования процессора
	MetricTypeMemory        MetricType = "memory"         // использование памяти
	MetricTypeDisk          MetricType = "disk"           // использование диска
	MetricTypeNet           MetricType = "net"            // использование сети
	MetricTypeProcess       MetricType = "process"        // процессы
)

const (
	EmptyFocusedWindow = "EMPTY_FOCUSED_WINDOW"
)

type Metric interface {
	Type() MetricType
	Timestamp() time.Time
} // @name Metric

type CpuPercentMetric struct {
	value     float32
	timestamp time.Time
} // @name CpuPercentMetric

func NewCpuPercentMetric(value float32) *CpuPercentMetric {
	return &CpuPercentMetric{
		value:     value,
		timestamp: time.Now(),
	}
}

func (m *CpuPercentMetric) Type() MetricType {
	return MetricTypeCpuPercent
}
func (m *CpuPercentMetric) Timestamp() time.Time {
	return m.timestamp
}
func (m *CpuPercentMetric) Value() float32 {
	return m.value
}

type FocusedWindowMetric struct {
	value     string
	timestamp time.Time
} // @name FocusedWindowMetric

func NewFocusedWindowMetric(value string) *FocusedWindowMetric {
	return &FocusedWindowMetric{
		value:     value,
		timestamp: time.Now(),
	}
}

func (m *FocusedWindowMetric) Type() MetricType {
	return MetricTypeFocusedWindow
}
func (m *FocusedWindowMetric) Timestamp() time.Time {
	return m.timestamp
}
func (m *FocusedWindowMetric) Value() string {
	return m.value
}

// метрика памяти
//
//	минимум полей: Total есть из MemorySpecs, остальное вычисляется
type MemoryMetric struct {
	used      uint64
	timestamp time.Time
} // @name MemoryMetric

func NewMemoryMetric(used uint64) *MemoryMetric {
	return &MemoryMetric{
		used:      used,
		timestamp: time.Now(),
	}
}
func (m *MemoryMetric) Type() MetricType {
	return MetricTypeMemory
}
func (m *MemoryMetric) Timestamp() time.Time {
	return m.timestamp
}
func (m *MemoryMetric) Value() uint64 {
	return m.used
}

type DiskMetric struct {
	used      map[string]uint64
	timestamp time.Time
} // @name DiskMetric

func NewDiskMetric(used map[string]uint64) *DiskMetric {
	return &DiskMetric{
		used:      used,
		timestamp: time.Now(),
	}
}
func (m *DiskMetric) Type() MetricType {
	return MetricTypeDisk
}
func (m *DiskMetric) Timestamp() time.Time {
	return m.timestamp
}
func (m *DiskMetric) Value() map[string]uint64 {
	return m.used
}

type NetIOMetric struct {
	uploadMbps   float32
	downloadMbps float32
	timestamp    time.Time
}

func NewNetIOMetric(upMbps float32, downMbps float32) *NetIOMetric {
	return &NetIOMetric{
		uploadMbps:   upMbps,
		downloadMbps: downMbps,
		timestamp:    time.Now(),
	}
}

func (m *NetIOMetric) Type() MetricType {
	return MetricTypeNet
}

func (m *NetIOMetric) Timestamp() time.Time {
	return m.timestamp
}

func (m *NetIOMetric) UploadMbps() float32 {
	return m.uploadMbps
}

func (m *NetIOMetric) DownloadMbps() float32 {
	return m.downloadMbps
}

type ProcessMetric struct {
	processes []Process
	timestamp time.Time
}

func NewProcessMetric(processList []Process) *ProcessMetric {
	return &ProcessMetric{
		processes: processList,
		timestamp: time.Now(),
	}
}
func (m *ProcessMetric) Type() MetricType {
	return MetricTypeProcess
}

func (m *ProcessMetric) Timestamp() time.Time {
	return m.timestamp
}

func (m *ProcessMetric) ProcessList() []Process {
	return m.processes
}
