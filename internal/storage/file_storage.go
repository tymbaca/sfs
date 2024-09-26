package storage

import (
	"context"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/tymbaca/sfs/internal/chunk"
)

type FileStorage struct {
	baseDir string
}

func NewFileStorage(baseDir string) *FileStorage {
	return &FileStorage{
		baseDir: baseDir,
	}
}

func (s *FileStorage) StoreChunk(ctx context.Context, chunk chunk.Chunk) error {
	f, err := createFile(path.Join(s.baseDir, chunk.Filename, strconv.Itoa(int(chunk.ID))))
	if err != nil {
		return fmt.Errorf("can't create file for %s/%d: %w", chunk.Filename, chunk.ID, err)
	}

	if _, err = f.Write(chunk.Body); err != nil {
		return fmt.Errorf("can't write to file for %s/%d: %w", chunk.Filename, chunk.ID, err)
	}

	return nil
}

// what if concurrent write of same file-chunk?
func createFile(fullPath string) (*os.File, error) {
	if err := os.MkdirAll(path.Dir(fullPath), os.ModePerm); err != nil {
		return nil, err
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}

	return f, nil
}
