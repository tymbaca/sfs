package chunkio

import (
	"io"
)

// NewWriter creates new [Writer]
func NewWriter(w io.WriterAt, start, limit int64) *Writer {
	return &Writer{
		w:      w,
		start:  start,
		offset: start,
		limit:  limit,
	}
}

type Writer struct {
	w      io.WriterAt
	start  int64
	offset int64
	limit  int64
}

// Write writes p to underlying [io.WriterAt]. When it hits the w.limit and
// wrote less then len(p) bytes it returns written bytes count and [io.ErrShortWrite].
// If limit already exceeded - it returns 0 and [io.EOF]
func (w *Writer) Write(p []byte) (int, error) {
	isShortWrite := false

	// Check if we've reached the end of the chunk
	if w.offset >= w.limit {
		return 0, io.EOF
	}

	// If the full p read will jump over limit - shorten p
	if w.offset+int64(len(p)) > w.limit {
		p = p[:(w.limit - w.offset)]
		isShortWrite = true
	}

	n, err := w.w.WriteAt(p, int64(w.offset))
	if err != nil {
		return n, err
	}

	w.offset += int64(n)

	if isShortWrite {
		return n, io.ErrShortWrite
	}

	return n, nil
}
