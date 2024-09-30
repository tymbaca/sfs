package transport

import (
	"context"
	"net"

	"github.com/tymbaca/sfs/internal/chunk"
)

type Transporter interface {
	// Sends the chunk to peer
	SendChunk(ctx context.Context, chunk chunk.Chunk) error
	// Returns all chunks of the file with name that respondent has.
	// If there is no chunks that match ids - empty silce will be returned without error.
	RecvAll(ctx context.Context, name string) ([]chunk.Chunk, error)
	// Returns chunks that respondent has.
	// If there is no chunks that match ids - empty silce will be returned without error.
	RecvChunk(ctx context.Context, name string, ids []uint64) (result []chunk.Chunk, err error)
}

type NetTransport struct {
	addr net.Addr
}

func (t *NetTransport) SendChunk(ctx context.Context, chunk chunk.Chunk) error

// WARN: need conn for each chunk
func (t *NetTransport) RecvAll(ctx context.Context, name string) ([]chunk.Chunk, error)

func (t *NetTransport) RecvChunk(ctx context.Context, name string, ids []uint64) (result []chunk.Chunk, err error)
