package syncx

import (
	"context"

	"golang.org/x/sync/semaphore"
)

var maxReaders int64 = 1 << 30

type CtxRWMutex struct {
	sem *semaphore.Weighted
}

func NewCtxRWMutex() *CtxRWMutex {
	return &CtxRWMutex{
		sem: semaphore.NewWeighted(maxReaders),
	}
}

func (m *CtxRWMutex) RLock(ctx context.Context) error {
	return m.sem.Acquire(ctx, 1)
}

func (m *CtxRWMutex) RUnlock() {
	m.sem.Release(1)
}

func (m *CtxRWMutex) Lock(ctx context.Context) error {
	return m.sem.Acquire(ctx, maxReaders)
}

func (m *CtxRWMutex) Unlock() {
	m.sem.Release(maxReaders)
}
