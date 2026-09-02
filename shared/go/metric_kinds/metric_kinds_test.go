package metrickinds

import (
	"fmt"
	"testing"

	agentv1 "github.com/nougght/monitoring-system/shared/go/proto/gen/agent/v1"
	"github.com/stretchr/testify/require"
)

func TestRegisteredMetricKinds(t *testing.T) {
	for kind, name := range agentv1.MetricKind_name {
		_, ok := MetricKinds[kind]
		require.True(t, ok || kind == 0, fmt.Sprintf("info for metric kind '%s' is not registered", name))
	}
}

func TestMetricKindInfoUniqueKey(t *testing.T) {
	keys := make(map[string]struct{}, len(MetricKinds))
	for kind, info := range MetricKinds {
		_, exists := keys[info.Key]
		require.False(t, exists, "key for metric kind '%s' is duplicated", agentv1.MetricKind_name[int32(kind)])
	}
}
