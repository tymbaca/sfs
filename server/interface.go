package sfs

import (
	"context"

	"github.com/tymbaca/sfs/internal/chunk"
)

type storage interface {
	StoreChunk(ctx context.Context, chunk chunk.Chunk) error
}
