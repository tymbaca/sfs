package chunk

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Chunk struct {
	ID       uint64
	Filename string
	Size     uint64
	Body     io.Reader
}

const _delimiter = '*'

const chunkFmt = `{
        ID: %d,
        Filename: %s,
        Size: %d,
        Body (utf-8): '%s'
}`

func (ch Chunk) String() string {
	return fmt.Sprintf(chunkFmt, ch.ID, ch.Filename, ch.Size, ch.Body)
}

func WriteChunk(w io.Writer, chunk Chunk) error {
	if _, err := w.Write([]byte{_delimiter}); err != nil {
		return fmt.Errorf("can't write delimiter: %w", err)
	}

	// we need len of bytes, not len of utf-8 symbols, so we use [len]
	if err := binary.Write(w, binary.LittleEndian, uint64(len(chunk.Filename))); err != nil {
		return fmt.Errorf("can't write filename size: %w", err)
	}

	if _, err := w.Write([]byte(chunk.Filename)); err != nil {
		return fmt.Errorf("can't write filename: %w", err)
	}

	if err := binary.Write(w, binary.LittleEndian, chunk.ID); err != nil {
		return fmt.Errorf("can't write chunk ID: %w", err)
	}

	if err := binary.Write(w, binary.LittleEndian, chunk.Size); err != nil {
		return fmt.Errorf("can't write chunk size: %w", err)
	}

	if _, err := io.Copy(w, chunk.Body); err != nil {
		return fmt.Errorf("can't write chunk body: %w", err)
	}

	return nil
}

func ReadChunk(r io.Reader) (Chunk, error) {
	delim := make([]byte, 1)
	_, err := r.Read(delim)
	if err != nil {
		return Chunk{}, fmt.Errorf("can't read first byte of chunk: %w", err)
	}

	if delim[0] != _delimiter {
		return Chunk{}, fmt.Errorf("incorrect delimiter, expected: '%s', got '%s'", string(_delimiter), delim)
	}

	var filenameSize uint64
	if err := binary.Read(r, binary.LittleEndian, &filenameSize); err != nil {
		return Chunk{}, fmt.Errorf("can't read filename size from chunk: %w", err)
	}

	filename := make([]byte, filenameSize)
	n, err := r.Read(filename)
	if err != nil {
		return Chunk{}, fmt.Errorf("can't read filename from chunk: %w", err)
	}

	if n < int(filenameSize) {
		return Chunk{}, fmt.Errorf("can't read filename from chunk: got (%d) less then expected (%d)", n, filenameSize)
	}

	var id uint64
	if err := binary.Read(r, binary.LittleEndian, &id); err != nil {
		return Chunk{}, fmt.Errorf("can't read ID from chunk: %w", err)
	}

	var bodySize uint64
	if err := binary.Read(r, binary.LittleEndian, &bodySize); err != nil {
		return Chunk{}, fmt.Errorf("can't read body size from chunk: %w", err)
	}

	return Chunk{
		ID:       id,
		Filename: string(filename),
		Size:     bodySize,
		Body:     io.LimitReader(r, int64(bodySize)), // to not suck in next chunks
	}, nil
}
