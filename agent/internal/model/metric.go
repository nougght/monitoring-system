package model

import "time"

type MetricType string

const (
	MetricTypeCpuPercent    MetricType = "cpu_percent"
	MetricTypeFocusedWindow MetricType = "focused_window"
)

type Metric interface {
	Type() MetricType
	Timestamp() time.Time
}

type CpuPercentMetric struct {
	value     float64
	timestamp time.Time
}

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
}

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

type Metrics struct {
	FocusedWindow string    `json:"focusedWindow"`
	CpuPercent    float64   `json:"cpuPercent"`
	Timestamp     time.Time `json:"timestamp"`
}
