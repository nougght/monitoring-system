package metrics_model

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type AgentSeriesKey struct {
	Kind  int32
	Label string
}

type MetricState struct {
	Value float64
	Ts    time.Time
	// Ring
}

type Snapshot struct {
	AgentID  uuid.UUID
	LastSeen time.Time

	mu     sync.RWMutex
	values map[AgentSeriesKey]MetricState
}

func (s *Snapshot) UpdateMetrics(metrics []*MetricSample) {
	s.mu.Lock()
	for _, m := range metrics {
		key := AgentSeriesKey{Kind: m.Kind, Label: m.Label}
		s.values[key] = MetricState{
			Value: m.Value,
			Ts:    m.Timestamp,
		}
	}
	s.mu.Unlock()
}
