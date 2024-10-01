package chunks

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

const chunkFmt = `{
        ID: %d,
        Filename: %s,
        Size: %d,
        Body (utf-8): '%s'
}`

func (ch Chunk) String() string {
	return fmt.Sprintf(chunkFmt, ch.ID, ch.Filename, ch.Size, ch.Body)
}

func SendChunk(w io.Writer, chunk Chunk) error {
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

func RecvChunk(r io.Reader) (Chunk, error) {
	var filenameSize uint64
	if err := binary.Read(r, binary.LittleEndian, &filenameSize); err != nil {
		return Chunk{}, fmt.Errorf("can't read filename size from chunk: %w", err)
	}

	filename := make([]byte, filenameSize)
	_, err := io.ReadFull(r, filename)
	if err != nil {
		return Chunk{}, fmt.Errorf("can't read filename from chunk: %w", err)
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
