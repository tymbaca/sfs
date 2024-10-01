package sfs

import (
	"context"

	"github.com/tymbaca/sfs/internal/chunks"
)

type storage interface {
	StoreChunk(ctx context.Context, chunk chunks.Chunk) error
	GetChunk(ctx context.Context, name string, id uint64) (chunks.Chunk, error)
	ListChunkIDs(ctx context.Context, name string) ([]uint64, error)
}
