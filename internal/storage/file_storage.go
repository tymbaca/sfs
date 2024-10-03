package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"strconv"

	"github.com/tymbaca/sfs/internal/chunks"
	"github.com/tymbaca/sfs/internal/common"
	"github.com/tymbaca/sfs/internal/files"
	"github.com/tymbaca/sfs/internal/logger"
)

type FileStorage struct {
	baseDir string
}

func NewFileStorage(baseDir string) *FileStorage {
	return &FileStorage{
		baseDir: baseDir,
	}
}

func (s *FileStorage) StoreChunk(ctx context.Context, chunk chunks.Chunk) error {
	f, err := files.CreateFile(path.Join(s.baseDir, chunk.Filename, strconv.Itoa(int(chunk.ID))))
	if err != nil {
		return fmt.Errorf("can't create file for %s/%d: %w", chunk.Filename, chunk.ID, err)
	}

	if _, err = io.Copy(f, chunk.Body); err != nil {
		return fmt.Errorf("can't write to file for %s/%d: %w", chunk.Filename, chunk.ID, err)
	}

	return nil
}

// GetChunk gets the chunk with file io.Reader inside. It's the called responsibility to close
// the file.
func (s *FileStorage) GetChunk(ctx context.Context, name string, id uint64) (chk chunks.Chunk, closeChk func() error, err error) {
	// Open the file
	f, err := os.Open(path.Join(s.baseDir, name, strconv.Itoa(int(id))))
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		return chunks.Chunk{}, nil, common.ErrNotFound
	} else if err != nil {
		return chunks.Chunk{}, nil, fmt.Errorf("can't get the chunk: %w", err)
	}

	// Get stat for size
	stat, err := f.Stat()
	if err != nil {
		return chunks.Chunk{}, nil, fmt.Errorf("can't get chunk file stat: %w", err)
	}

	return chunks.Chunk{
		ID:       id,
		Filename: name,
		Size:     uint64(stat.Size()),
		Body:     f,
	}, f.Close, nil
}

func (s *FileStorage) ListChunkIDs(ctx context.Context, name string) ([]uint64, error) {
	entries, err := os.ReadDir(path.Join(s.baseDir, name))
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		return nil, common.ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("can't get file's chunks IDs: %w", err)
	}

	ids := make([]uint64, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			id, err := strconv.Atoi(e.Name())
			if err != nil {
				logger.Logf("got non-int name in chunks folder: path: %s", path.Join(s.baseDir, name, e.Name()))
				continue
			}

			ids = append(ids, uint64(id))
		}
	}

	return ids, nil
}
