package agent_server

import (
	"github.com/google/uuid"
	"github.com/nougght/monitoring-system/server/internal/model/metrics"
	agentv1 "github.com/nougght/monitoring-system/shared/go/proto/gen/agent/v1"
	"github.com/nougght/monitoring-system/shared/go/util"
)

func convertMetricSampleFromProto(p *agentv1.MetricSample) *metrics.MetricSample {
	if p == nil {
		return nil
	}

	return &metrics.MetricSample{
		Kind:      int32(p.Kind),
		Label:     p.Label,
		Value:     p.Value,
		Timestamp: p.Timestamp.AsTime(),
	}
}

func convertMetricsBatchWithAgentIDFromProto(p *agentv1.MetricsBatch, agetnID uuid.UUID) *metrics.MetricsBatch {
	if p == nil {
		return nil
	}

	return &metrics.MetricsBatch{
		AgentID: agetnID,
		ID:      p.BatchId,
		Metrics: util.Map(p.Samples, convertMetricSampleFromProto),
	}
}
