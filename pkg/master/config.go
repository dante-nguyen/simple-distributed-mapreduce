package master

import (
	"errors"
)

type Config struct {
	InputFiles []string // paths relative to NFS root
	MaxWorkers int
}

func validateConfig(cfg Config) error {
	switch {
	case len(cfg.InputFiles) == 0:
		return errors.New("no input files")
	case cfg.MaxWorkers <= 0:
		return errors.New("invalid number of max workers")
	default:
		return nil
	}
}
