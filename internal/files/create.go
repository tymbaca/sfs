package files

import (
	"os"
	"path"
)

// what if concurrent write of same file-chunk?
func Create(fullPath string) (*os.File, error) {
	if err := os.MkdirAll(path.Dir(fullPath), os.ModePerm); err != nil {
		return nil, err
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}

	return f, nil
}
