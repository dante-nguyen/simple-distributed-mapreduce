package master

import (
	"sync"

	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/task"
)

type maptask struct {
	mu sync.RWMutex

	resource   task.MapResource
	assigned   bool
	assignedTo string
}

func newMaptask(path string) *maptask {
	return &maptask{
		resource:   task.MapResource{Path: path},
		assigned:   false,
		assignedTo: "",
	}
}

func (mt *maptask) path() string {
	mt.mu.RLock()
	defer mt.mu.RUnlock()

	return mt.resource.Path
}
