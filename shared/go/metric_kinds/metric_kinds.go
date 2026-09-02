package metrickinds

import (
	"maps"
	"slices"
	"time"

	agentv1 "github.com/nougght/monitoring-system/shared/go/proto/gen/agent/v1"
)

type MetricKindInfo struct {
	Kind            int32  `db:"kind"`
	Key             string `db:"key"`
	Unit            string `db:"unit"`
	Agg             int32  `db:"agg"`
	LabelName       string `db:"label_name"`
	Description     string `db:"description"`
	DefaultInterval time.Duration
}

// TEMP
const (
	AggMax = iota + 1
	AggMin
	AggAvg
)

var MetricKinds = map[int32]MetricKindInfo{
	int32(agentv1.MetricKind_CPU_PERCENT): {
		Kind:      int32(agentv1.MetricKind_CPU_PERCENT),
		Key:       "cpu_usage",
		Unit:      "percent",
		Agg:       AggAvg,
		LabelName: "core",
	},
	int32(agentv1.MetricKind_MEM_USED): {
		Kind:      int32(agentv1.MetricKind_MEM_USED),
		Key:       "memory_usage",
		Unit:      "bytes",
		Agg:       AggAvg,
		LabelName: "type",
	},
	int32(agentv1.MetricKind_DISK_USED): {
		Kind:      int32(agentv1.MetricKind_DISK_USED),
		Key:       "disk_usage",
		Unit:      "bytes",
		Agg:       AggAvg,
		LabelName: "mount_point",
	},
	int32(agentv1.MetricKind_NET_UPLOAD): {
		Kind:      int32(agentv1.MetricKind_NET_UPLOAD),
		Key:       "network_upload",
		Unit:      "bytes",
		Agg:       AggAvg,
		LabelName: "interface",
	},
	int32(agentv1.MetricKind_NET_DOWNLOAD): {
		Kind:      int32(agentv1.MetricKind_NET_DOWNLOAD),
		Key:       "network_download",
		Unit:      "bytes",
		Agg:       AggAvg,
		LabelName: "interface",
	},
}

func MetricKindList() []MetricKindInfo {
	return slices.Collect(maps.Values(MetricKinds))
}
