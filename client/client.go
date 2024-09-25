package sfs

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

type Client struct {
	addrs     []string
	chunkSize uint64 // bytes
}

type chunk struct {
	ID       uint64
	Filename string
	Size     uint64
	Body     []byte
}

func NewClient(addrs string, chunkSize uint64) *Client {
	return &Client{
		addrs:     strings.Split(addrs, ","),
		chunkSize: chunkSize,
	}
}

func (c *Client) Upload(ctx context.Context, name string, r io.Reader) error {
	chunks, err := formChunks(r, name, c.chunkSize)
	if err != nil {
		return err
	}

	conns, err := c.connect(ctx)
	if err != nil {
		return err
	}
	defer closeConns(conns)

	err = uploadChunks(ctx, conns, chunks)
	if err != nil {
		return err
	}

	return nil
}

func uploadChunks(_ context.Context, conns []net.Conn, chunks <-chan chunk) error {
	var wg sync.WaitGroup
	wg.Add(len(conns))

	for _, conn := range conns {
		conn := conn
		go func() {
			defer wg.Done()
			for chunk := range chunks {
				err := writeChunk(conn, chunk)
				if err != nil {
					panic(err)
				}
			}
		}()
	}

	wg.Wait()

	return nil
}

func writeChunk(w io.Writer, chunk chunk) error {
	if _, err := w.Write([]byte("$")); err != nil {
		return err
	}

	// we need len of bytes, not len of utf-8 symbols, so we use [len]
	if err := binary.Write(w, binary.LittleEndian, uint64(len(chunk.Filename))); err != nil {
		return err
	}

	if _, err := w.Write([]byte(chunk.Filename)); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, chunk.ID); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, chunk.Size); err != nil {
		return err
	}

	if _, err := w.Write(chunk.Body); err != nil {
		return err
	}

	return nil
}

const _delimiter = '$'

func ReadChunk(r io.Reader) (chunk, error) {
	delim := make([]byte, 1)
	_, err := r.Read(delim)
	if err != nil {
		return chunk{}, fmt.Errorf("can't read first byte of chunk: %w", err)
	}

	if delim[0] != '$' {
		return chunk{}, fmt.Errorf("incorrect delimiter, expected: '%s', got '%s'", string(_delimiter), delim)
	}

	var filenameSize uint64
	if err := binary.Read(r, binary.LittleEndian, &filenameSize); err != nil {
		return chunk{}, fmt.Errorf("can't read filename size from chunk: %w", err)
	}

	filename := make([]byte, filenameSize)
	n, err := r.Read(filename)
	if err != nil {
		return chunk{}, fmt.Errorf("can't read filename from chunk: %w", err)
	}

	if n < int(filenameSize) {
		return chunk{}, fmt.Errorf("can't read filename from chunk: got (%d) less then expected (%d)", n, filenameSize)
	}

	var id uint64
	if err := binary.Read(r, binary.LittleEndian, &id); err != nil {
		return chunk{}, fmt.Errorf("can't read ID from chunk: %w", err)
	}

	var bodySize uint64
	if err := binary.Read(r, binary.LittleEndian, &bodySize); err != nil {
		return chunk{}, fmt.Errorf("can't read body size from chunk: %w", err)
	}

	body := make([]byte, bodySize)
	n, err = r.Read(body)
	if err != nil {
		return chunk{}, fmt.Errorf("can't read body from chunk: %w", err)
	}

	if n < int(bodySize) {
		return chunk{}, fmt.Errorf("can't read body from chunk: got (%d) less then expected (%d)", n, bodySize)
	}

	return chunk{
		ID:       id,
		Filename: string(filename),
		Size:     bodySize,
		Body:     body,
	}, nil
}

func (c *Client) connect(_ context.Context) ([]net.Conn, error) {
	conns := make([]net.Conn, 0, len(c.addrs))
	for _, addr := range c.addrs {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			closeConns(conns)
			return nil, err
		}

		conns = append(conns, conn)
	}

	return conns, nil
}

func formChunks(r io.Reader, name string, size uint64) (<-chan chunk, error) {
	if size < 1 {
		panic("can't split byte non-positive size")
	}

	ch := make(chan chunk)
	go func() {
		defer close(ch)
		for id := 0; ; id++ {
			buf := make([]byte, size)
			n, err := r.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}

				panic("shit happens")
			}

			ch <- chunk{
				ID:       uint64(id),
				Filename: name,
				Size:     uint64(n),
				Body:     buf[:n],
			}
		}
	}()

	return ch, nil
}

func closeConns(conns []net.Conn) {
	for _, conn := range conns {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}
}
