package flagx

import (
	"os"
	"path/filepath"

	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/errx"
	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/fsx"
)

type DirValue struct {
	Path string
}

func (dv *DirValue) String() string {
	return "Dir:" + dv.Path
}

func (dv *DirValue) Set(val string) error {
	info, err := os.Stat(val)
	if err != nil {
		return errx.WithContext(err, "stat path")
	}

	if !info.IsDir() {
		return fsx.ErrNotADirectory
	}

	abs, err := filepath.Abs(val)
	if err != nil {
		return errx.WithContext(err, "get absolute path")
	}

	dv.Path = abs
	return nil
}
