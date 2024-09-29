package chunkio

import (
	"fmt"
	"io"
	"os"
)

func Split(r io.ReaderAt, totalSize, chunkSize int64) []*Reader {
	offset := int64(0)
	chunks := make([]*Reader, 0, totalSize/chunkSize+1)

	for offset < totalSize {
		limit := offset + chunkSize
		if limit > totalSize {
			limit = totalSize
		}

		chunk := NewReader(r, offset, limit)
		chunks = append(chunks, chunk)
		offset += chunkSize
	}

	return chunks
}

func SplitFile(f *os.File, chunkSize int64) ([]*Reader, error) {
	stat, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("can't split file into chunks: %w", err)
	}

	size := stat.Size()

	return Split(f, size, chunkSize), nil
}

// Four different cases:
// ------|------|--
// -------------|--
// -------------|
// ---------    |

func NewReader(r io.ReaderAt, start, limit int64) *Reader {
	return &Reader{
		r:      r,
		start:  start,
		offset: start,
		limit:  limit,
	}
}

type Reader struct {
	r      io.ReaderAt
	start  int64
	offset int64
	limit  int64
}

func (r *Reader) Read(p []byte) (int, error) {
	// Check if we've reached the end of the chunk
	if r.offset >= r.limit {
		return 0, io.EOF
	}

	// If the full p read will jump over limit - shorten p
	if r.offset+int64(len(p)) > r.limit {
		p = p[:(r.limit - r.offset)]
	}

	n, err := r.r.ReadAt(p, int64(r.offset))
	if err != nil {
		return 0, err
	}

	r.offset += int64(n)

	return n, nil
}

func (r *Reader) Size() int64 {
	return r.limit - r.start
}
