package task

import "sync"

// Task is an abstraction of a worker task (either map or reduce).
// This struct is concurrency-safe and must not be copied after first use.
type Task struct {
	mu        sync.RWMutex
	typ       Type
	mapRes    MapResource
	reduceRes ReduceResource
}

func New() Task {
	return Task{
		typ: TypeNone,
	}
}

func (t *Task) IsSet() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.typ != TypeNone
}

// SetMap returns false if the task is already set.
func (t *Task) SetMap(path string) bool {
	if t.IsSet() {
		return false
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	t.typ = TypeMap
	t.mapRes = MapResource{Path: path}
	return true
}

func (t *Task) GetMap() (MapResource, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.typ != TypeMap {
		return MapResource{}, false
	}

	return t.mapRes, true
}
