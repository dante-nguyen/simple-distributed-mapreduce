package inputregistry

import (
	"fmt"

	"github.com/nlduy0310/simple-distributed-mapreduce/filex"
)

type Registry struct {
	inputFilePaths []string
}

func FromPaths(paths []string) (*Registry, error) {
	validatePaths(paths)

	return &Registry{
		inputFilePaths: paths,
	}, nil
}

func (r *Registry) Size() int {
	return len(r.inputFilePaths)
}

func validatePaths(filePaths []string) error {
	var err error
	for _, path := range filePaths {
		if err = validatePath(path); err != nil {
			return err
		}
	}

	return nil
}

func validatePath(filePath string) error {
	switch {
	case !filex.IsFile(filePath):
		return fmt.Errorf("%s is not a file", filePath)
	case !filex.IsReadable(filePath):
		return fmt.Errorf("file %s is not readable", filePath)
	default:
		return nil
	}
}
