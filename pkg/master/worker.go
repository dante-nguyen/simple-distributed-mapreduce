package master

import (
	"context"
	"time"

	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/syncx"
)

// worker is concurrency-safe and needs to be used by pointer
type worker struct {
	addr string // no need for mutex cause we only read for now

	heartbeatMutex *syncx.CtxRWMutex
	heartbeat      time.Time
}

func newWorker(addr string) *worker {
	return &worker{
		addr:           addr,
		heartbeatMutex: syncx.NewCtxRWMutex(),
		heartbeat:      time.Now(),
	}
}

func (w *worker) address() string {
	return w.addr
}

func (w *worker) lastHeartbeat(ctx context.Context) (time.Time, error) {
	if err := w.heartbeatMutex.RLock(ctx); err != nil {
		return time.Time{}, err
	}
	defer w.heartbeatMutex.RUnlock()

	return w.heartbeat, nil
}

func (w *worker) recordHeartbeat(ctx context.Context, t time.Time) error {
	if err := w.heartbeatMutex.Lock(ctx); err != nil {
		return err
	}
	defer w.heartbeatMutex.Unlock()

	w.heartbeat = t
	return nil
}
