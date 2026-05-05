package fsx

import "errors"

var (
	ErrNotADirectory = errors.New("not a directory")
	ErrNotAFile      = errors.New("not a file")
)
