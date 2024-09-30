package transport

import (
	"context"

	"github.com/tymbaca/sfs/internal/chunk"
)

type Transport interface {
	// Sends the chunk to peer
	SendChunk(ctx context.Context, chunk chunk.Chunk) error
	// Returns the chunk ids of the file that respondent has.
	ListIDs(ctx context.Context, name string) ([]uint64, error)
	RecvChunk(ctx context.Context, name string, id uint64) (chunk.Chunk, error)
}

type TCPTransport struct {
	addr string
}

func NewTCPTransport(addr string) *TCPTransport {
	return &TCPTransport{
		addr: addr,
	}
}

func (t *TCPTransport) SendChunk(ctx context.Context, chunk chunk.Chunk) error

func (t *TCPTransport) ListIDs(ctx context.Context, name string) ([]uint64, error)

func (t *TCPTransport) RecvChunk(ctx context.Context, name string, id uint64) (chunk.Chunk, error)
