package common

import "errors"

var (
	ErrChunkNotFound = errors.New("chunk not found")
	ErrFileNotFound  = errors.New("file not found")
)
