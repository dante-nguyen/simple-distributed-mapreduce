package worker

import (
	"fmt"
	"math/rand"
)

func genId() string {
	return fmt.Sprintf("worker-%d", rand.Uint32())
}
