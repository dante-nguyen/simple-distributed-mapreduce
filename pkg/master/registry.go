package master

import (
	"errors"
	"sync"
	"time"
)

var (
	errWorkerNotFound = errors.New("worker not found")
	errWorkerExists   = errors.New("worker already exists")
)

type registry struct {
	mu      sync.RWMutex
	workers map[string]*worker
}

func newRegistry() *registry {
	return &registry{
		workers: make(map[string]*worker),
	}
}

func (r *registry) exist(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.workers[name]
	return ok
}

func (r *registry) register(name, address string) error {
	if r.exist(name) {
		return errWorkerExists
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.workers[name] = newWorker(address)
	return nil
}

func (r *registry) lastHeartbeat(name string) (time.Time, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	w, ok := r.workers[name]
	if !ok {
		return time.Time{}, errWorkerNotFound
	}

	return w.lastHeartbeat(), nil
}

func (r *registry) recordHeartbeat(name string, at time.Time) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	w, ok := r.workers[name]
	if !ok {
		return errWorkerNotFound
	}

	w.recordHeartbeat(at)
	return nil
}

func (r *registry) remove(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.workers, name)
}

func (r *registry) names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ret := make([]string, 0, len(r.workers))
	for name := range r.workers {
		ret = append(ret, name)
	}
	return ret
}
