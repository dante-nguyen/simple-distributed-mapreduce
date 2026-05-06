package master

import (
	"errors"
	"fmt"

	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/errx"
	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/fsx"
)

type Config struct {
	InputFiles []string
}

func validateConfig(cfg Config) error {
	if len(cfg.InputFiles) == 0 {
		return errors.New("no input files")
	}

	for _, path := range cfg.InputFiles {
		if is, err := fsx.IsFile(path); err != nil {
			return errx.WithContext(err, fmt.Sprintf("path %s", path))
		} else if !is {
			return errx.WithContext(fsx.ErrNotAFile, fmt.Sprintf("path %s", path))
		}
	}

	return nil
}
