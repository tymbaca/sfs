package sfs

import (
	"context"

	"github.com/tymbaca/sfs/internal/chunks"
)

type storage interface {
	StoreChunk(ctx context.Context, chunk chunks.Chunk) error
}
