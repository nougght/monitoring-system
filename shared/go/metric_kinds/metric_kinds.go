package metrickinds

import (
	"time"
)

type MetricKindInfo struct {
	Kind            int32
	Key             string
	Unit            string
	Agg             string
	LabelName       string
	DefaultInterval time.Duration
}

var MetricKinds = map[int32]MetricKindInfo{}
