package flagx

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/errx"
)

var (
	errFailToStat    = errors.New("failed to stat")
	errNotADirectory = errors.New("not a directory")
	errFailToAbs     = errors.New("failed to get absolute path")
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
		return errx.Chain(errFailToStat, err)
	}

	if !info.IsDir() {
		return errNotADirectory
	}

	abs, err := filepath.Abs(val)
	if err != nil {
		return errx.Chain(errFailToAbs, err)
	}

	dv.Path = abs
	return nil
}
