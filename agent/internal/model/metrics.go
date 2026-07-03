package model

import "time"

type Metrics struct {
	FocusedWindow *string   `json:"focusedWindow,omitempty"`
	CpuPercent    *float64  `json:"cpuPercent,omitempty"`
	Timestamp     time.Time `json:"timestamp"`
}
