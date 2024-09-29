package chunkio

import (
	"bytes"
	"testing"
)

type writerAt struct {
	buf []byte
}

func (w *writerAt) WriteAt(p []byte, off int64) (n int, err error) {
	if int(off)+len(p) > len(w.buf) {
		p = p[:len(w.buf)-int(off)]
	}
}

func TestWriter(t *testing.T) {
	bytes.NewBuffer(nil)
}
