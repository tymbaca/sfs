package transport

import (
	"context"
	"errors"
	"net"

	"github.com/tymbaca/sfs/internal/chunks"
)

type Transport interface {
	// Sends the chunk to peer
	SendChunk(ctx context.Context, chunk chunks.Chunk) error
	// Returns the chunk ids of the file that respondent has.
	ListIDs(ctx context.Context, name string) ([]uint64, error)
	RecvChunk(ctx context.Context, name string, id uint64) (chunks.Chunk, error)
	Close() error
}

type TCPTransport struct {
	addr string
	conn net.Conn
}

func NewTCPTransport(addr string) *TCPTransport {
	return &TCPTransport{
		addr: addr,
	}
}

func (t *TCPTransport) ensureDial() (err error) {
	if t.conn == nil {
		t.conn, err = net.Dial("tcp", t.addr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *TCPTransport) SendChunk(ctx context.Context, chk chunks.Chunk) error {
	if err := t.ensureDial(); err != nil {
		return err
	}

	_, err := t.conn.Write([]byte("*"))
	if err != nil {
		return err
	}

	return chunks.SendChunk(t.conn, chk)
}

func (t *TCPTransport) ListIDs(ctx context.Context, name string) ([]uint64, error) {
	return nil, errors.New("list ids: not implemented")
}

func (t *TCPTransport) RecvChunk(ctx context.Context, name string, id uint64) (chunks.Chunk, error) {
	return chunks.Chunk{}, errors.New("recv chunk: not implemented")
}

func (t *TCPTransport) Close() error {
	return t.conn.Close()
}
