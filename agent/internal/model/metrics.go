package model

import "time"

type Metrics struct {
	FocusedWindow *string   `json:"focusedWindow,omitempty"`
	CpuPercent    *float64  `json:"cpuPercent,omitempty"`
	MemoryUsed    *uint64   `json:"memoryUsed,omitempty"`
	Timestamp     time.Time `json:"timestamp"`
} // @name Metrics
