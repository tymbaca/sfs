package sfs

import (
	"context"
	"errors"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/tymbaca/sfs/internal/chunk"
)

type Client struct {
	addrs     []string
	chunkSize uint64 // bytes
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

func uploadChunks(_ context.Context, conns []net.Conn, chunks <-chan chunk.Chunk) error {
	var wg sync.WaitGroup
	wg.Add(len(conns))

	for _, conn := range conns {
		conn := conn
		go func() {
			defer wg.Done()
			for chk := range chunks {
				err := chunk.WriteChunk(conn, chk)
				if err != nil {
					panic(err)
				}
			}
		}()
	}

	wg.Wait()

	return nil
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

func formChunks(r io.Reader, name string, size uint64) (<-chan chunk.Chunk, error) {
	if size < 1 {
		panic("can't split byte non-positive size")
	}

	ch := make(chan chunk.Chunk)
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

			ch <- chunk.Chunk{
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
