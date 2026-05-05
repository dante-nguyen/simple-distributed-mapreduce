package master

import (
	"context"
	"errors"
	"time"

	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/syncx"
)

var (
	errWorkerNotFound = errors.New("worker not found")
	errWorkerExists   = errors.New("worker already exists")
)

type registry struct {
	mu      *syncx.CtxRWMutex
	workers map[string]*worker
}

func newRegistry() *registry {
	return &registry{
		mu:      syncx.NewCtxRWMutex(),
		workers: make(map[string]*worker),
	}
}

func (r *registry) exist(ctx context.Context, name string) (bool, error) {
	if err := r.mu.RLock(ctx); err != nil {
		return false, err
	}
	defer r.mu.RUnlock()

	_, ok := r.workers[name]
	return ok, nil
}

func (r *registry) register(ctx context.Context, name, address string) error {
	if exist, err := r.exist(ctx, name); err != nil {
		return err
	} else if exist {
		return errWorkerExists
	}

	if err := r.mu.Lock(ctx); err != nil {
		return err
	}
	defer r.mu.Unlock()

	r.workers[name] = newWorker(address)
	return nil
}

func (r *registry) lastHeartbeat(ctx context.Context, name string) (time.Time, error) {
	if err := r.mu.RLock(ctx); err != nil {
		return time.Time{}, err
	}
	defer r.mu.RUnlock()

	w, ok := r.workers[name]
	if !ok {
		return time.Time{}, errWorkerNotFound
	}

	ret, err := w.lastHeartbeat(ctx)
	if err != nil {
		return time.Time{}, err
	}

	return ret, nil
}

func (r *registry) recordHeartbeat(ctx context.Context, name string, at time.Time) error {
	if err := r.mu.RLock(ctx); err != nil {
		return err
	}
	defer r.mu.RUnlock()

	w, ok := r.workers[name]
	if !ok {
		return errWorkerNotFound
	}

	return w.recordHeartbeat(ctx, at)
}

func (r *registry) remove(ctx context.Context, name string) error {
	if err := r.mu.Lock(ctx); err != nil {
		return err
	}
	defer r.mu.Unlock()

	delete(r.workers, name)
	return nil
}

func (r *registry) names(ctx context.Context) ([]string, error) {
	if err := r.mu.RLock(ctx); err != nil {
		return nil, err
	}
	defer r.mu.RUnlock()

	ret := make([]string, 0, len(r.workers))
	for name := range r.workers {
		ret = append(ret, name)
	}
	return ret, nil
}
