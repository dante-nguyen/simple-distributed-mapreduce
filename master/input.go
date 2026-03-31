package master

import (
	"fmt"

	"github.com/nlduy0310/simple-distributed-mapreduce/filex"
)

func validateInput(files []string) error {
	for _, file := range files {
		if err := validateFile(file); err != nil {
			return err
		}
	}

	return nil
}

func validateFile(file string) error {
	if !filex.IsFile(file) {
		return fmt.Errorf("%s is not a file", file)
	} else if !filex.IsReadable(file) {
		return fmt.Errorf("file %s is not readable", file)
	}

	return nil
}
