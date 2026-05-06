package master

import (
	"sync"
	"time"
)

// worker is concurrency-safe and needs to be used by pointer
type worker struct {
	addr string // no need for mutex cause we only read for now

	heartbeatMutex sync.RWMutex
	heartbeat      time.Time
}

func newWorker(addr string) *worker {
	return &worker{
		addr:      addr,
		heartbeat: time.Now(),
	}
}

func (w *worker) address() string {
	return w.addr
}

func (w *worker) lastHeartbeat() time.Time {
	w.heartbeatMutex.RLock()
	defer w.heartbeatMutex.RUnlock()

	return w.heartbeat
}

func (w *worker) recordHeartbeat(t time.Time) {
	w.heartbeatMutex.Lock()
	defer w.heartbeatMutex.Unlock()

	w.heartbeat = t
}
