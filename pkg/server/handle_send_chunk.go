package sfs

import (
	"context"
	"fmt"
	"io"

	"github.com/tymbaca/sfs/internal/chunks"
	"github.com/tymbaca/sfs/internal/codes"
)

func (s *Server) handleSendChunk(ctx context.Context, conn io.ReadWriter) error {
	chk, err := chunks.RecvChunk(conn)
	if err != nil {
		return fmt.Errorf("can't receive chunk from client: %w", err)
	}

	if err = s.storage.StoreChunk(ctx, chk); err != nil {
		err = fmt.Errorf("can't store the chunk: %w", err)
		writeCodeMsg(conn, codes.Internal, err.Error())
		return err
	}

	return writeCodeMsg(conn, codes.Ok, "uploaded")
}
