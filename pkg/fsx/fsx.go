package fsx

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/nlduy0310/simple-distributed-mapreduce/pkg/errx"
)

type WalkFilter func(path string, de fs.DirEntry) bool

func FilterDir(path string, de fs.DirEntry) bool {
	return de.IsDir()
}

func FilterFile(path string, de fs.DirEntry) bool {
	return de.Type().IsRegular()
}

func IsDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, errx.WithContext(err, "stat path")
	}

	return info.IsDir(), nil
}

func IsFile(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, errx.WithContext(err, "stat path")
	}

	return info.Mode().IsRegular(), nil
}

func CollectPaths(rootDir string, filters ...WalkFilter) ([]string, error) {
	res := make([]string, 0)
	err := filepath.WalkDir(rootDir, func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		for _, filter := range filters {
			if !filter(path, de) {
				return nil
			}
		}

		res = append(res, path)
		return nil
	})

	if err != nil {
		return nil, errx.WithContext(err, "walk dir")
	}

	return res, nil
}
