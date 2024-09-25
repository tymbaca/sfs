package sfs

import (
	"context"
	"encoding/binary"
	"errors"
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
	chunks, err := split(r, c.chunkSize)
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
				conn.Write([]byte("$"))

				err := binary.Write(conn, binary.LittleEndian, chunk.ID)
				if err != nil {
					panic(err)
				}

				err = binary.Write(conn, binary.LittleEndian, chunk.Size)
				if err != nil {
					panic(err)
				}

				_, err = conn.Write(chunk.Body)
				if err != nil {
					panic(err)
				}
			}
		}()
	}

	wg.Wait()

	return errors.New("not implemented")
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

func split(r io.Reader, size uint64) (<-chan chunk, error) {
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
				ID:   uint64(id),
				Size: uint64(n),
				Body: buf[:n],
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
