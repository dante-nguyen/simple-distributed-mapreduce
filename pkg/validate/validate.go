package validate

import (
	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/fsx"
)

func EnsureIsFile(path string) error {
	is, err := fsx.IsFile(path)
	if err != nil {
		return err
	} else if !is {
		return fsx.ErrNotAFile
	}

	return nil
}

func EnsureIsDir(path string) error {
	is, err := fsx.IsDir(path)
	if err != nil {
		return err
	} else if !is {
		return fsx.ErrNotADirectory
	}

	return nil
}
