package workerregistry

import (
	"fmt"
	"sync"
	"time"
)

var workerAlreadyExistsError = fmt.Errorf("worker already exists")

type WorkerInfo struct {
	LastContactTime time.Time
}

// Registry is concurrency-safe.
// It contains a lock and must always be passed by pointer.
type Registry struct {
	mu         sync.RWMutex
	workersMap map[string]WorkerInfo // by worker id
}

func NewRegistry() *Registry {
	return &Registry{
		workersMap: make(map[string]WorkerInfo),
	}
}

func (r *Registry) Size() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.workersMap)
}

func (r *Registry) hasWorker(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.workersMap[id]
	return ok
}

func (r *Registry) Register(workerId string) error {
	switch {
	case r.hasWorker(workerId):
		return workerAlreadyExistsError
	default:
		r.doRegister(workerId)
		return nil
	}
}

func (r *Registry) doRegister(workerId string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.workersMap[workerId] = WorkerInfo{
		LastContactTime: time.Now(),
	}
}
