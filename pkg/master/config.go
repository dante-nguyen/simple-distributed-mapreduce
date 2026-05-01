package master

import (
	"fmt"

	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/errx"
	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/fsx"
)

type Config struct {
	InputFiles []string
}

func validateConfig(cfg Config) error {
	for _, path := range cfg.InputFiles {
		if is, err := fsx.IsFile(path); err != nil {
			return errx.WithContextMsg(err, fmt.Sprintf("file %s", path))
		} else if !is {
			return errx.WithContextMsg(fsx.ErrNotAFile, fmt.Sprintf("file %s", path))
		}
	}

	return nil
}
