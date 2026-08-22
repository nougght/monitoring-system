package metrics_model

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type seriesKey struct {
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
	values map[seriesKey]MetricState
}

func (s *Snapshot) UpdateMetrics(metrics []*MetricSample) {
	s.mu.Lock()
	for _, m := range metrics {
		key := seriesKey{Kind: m.Kind, Label: m.Label}
		s.values[key] = MetricState{
			Value: m.Value,
			Ts:    m.Timestamp,
		}
	}
	s.mu.Unlock()
}
