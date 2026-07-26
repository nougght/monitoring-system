package metrics

import (
	"maps"
	"sync"

	"github.com/google/uuid"
	model "github.com/nougght/monitoring-system/server/internal/model/event/metrics"
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

func (s *SnapshotCache) Update(agentID uuid.UUID, snapshot *model.Snapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache[agentID] = snapshot
}

func (s *SnapshotCache) Get(agentID uuid.UUID) (snapshot *model.Snapshot, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	snapshot, ok = s.cache[agentID]
	return
}

func (s *SnapshotCache) All() map[uuid.UUID]*model.Snapshot {
	// returning copy of map
	return maps.Clone(s.cache)
}
