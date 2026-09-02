package metrics

import (
	"maps"
	"sync"

	"github.com/google/uuid"
	metrics_model "github.com/nougght/monitoring-system/server/internal/model/metrics"
	model "github.com/nougght/monitoring-system/server/internal/model/metrics"
)

type SnapshotCache struct {
	mu    sync.RWMutex
	cache map[uuid.UUID]*model.Snapshot
}

func NewSnapshotCache() *SnapshotCache {
	return &SnapshotCache{
		cache: make(map[uuid.UUID]*model.Snapshot),
	}
}

func (s *SnapshotCache) UpdateMetrics(agentID uuid.UUID, batch metrics_model.MetricsBatch) {
	s.mu.Lock()
	snapshot, ok := s.cache[agentID]
	s.mu.Unlock()
	if !ok || snapshot == nil {
		snapshot = model.NewSnapshot(agentID)
		s.cache[agentID] = snapshot
	}
	snapshot.UpdateMetrics(batch.Metrics)
}

func (s *SnapshotCache) Get(agentID uuid.UUID) (snapshot *model.Snapshot, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	snapshot, ok = s.cache[agentID]
	return
}

func (s *SnapshotCache) All() map[uuid.UUID]*model.Snapshot {
	// copy of map
	return maps.Clone(s.cache)
}
