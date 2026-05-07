package master

import (
	"errors"
)

type Config struct {
	InputFiles []string // paths relative to NFS root
}

func validateConfig(cfg Config) error {
	if len(cfg.InputFiles) == 0 {
		return errors.New("no input files")
	}

	return nil
}
