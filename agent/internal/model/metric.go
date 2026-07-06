package model

import "time"

type MetricType string

const (
	MetricTypeFocusedWindow MetricType = "focused_window" // активное окно
	MetricTypeCpuPercent    MetricType = "cpu_percent"    // процент использования процессора
	MetricTypeMemory        MetricType = "memory"         // использование памяти
)

const (
	EmptyFocusedWindow = "EMPTY_FOCUSED_WINDOW"
)

type Metric interface {
	Type() MetricType
	Timestamp() time.Time
} // @name Metric

type CpuPercentMetric struct {
	value     float64
	timestamp time.Time
} // @name CpuPercentMetric

func NewCpuPercentMetric(value float64) *CpuPercentMetric {
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
func (m *CpuPercentMetric) Value() float64 {
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
//  минимум полей: Total есть из MemorySpecs, остальное вычисляется
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
