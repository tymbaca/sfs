package sfs

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/tymbaca/sfs/internal/chunks"
	"github.com/tymbaca/sfs/internal/codes"
	"github.com/tymbaca/sfs/internal/common"
)

func (s *Server) handleRecvChunk(ctx context.Context, conn io.ReadWriter) error {
	var filenameSize uint64
	if err := binary.Read(conn, binary.LittleEndian, &filenameSize); err != nil {
		return fmt.Errorf("can't read filename size from chunk: %w", err)
	}

	filename := make([]byte, filenameSize)
	_, err := io.ReadFull(conn, filename)
	if err != nil {
		return fmt.Errorf("can't read filename from chunk: %w", err)
	}

	var id uint64
	if err := binary.Read(conn, binary.LittleEndian, &id); err != nil {
		return fmt.Errorf("can't read ID from chunk: %w", err)
	}

	chk, closeChk, err := s.storage.GetChunk(ctx, string(filename), id)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			// normal case, not an error
			return writeCode(conn, codes.NotFound)
		}

		err = fmt.Errorf("can't get chunk from storage: %w", err)
		writeCodeMsg(conn, codes.Internal, err.Error())
		return err
	}
	defer closeChk()

	if err := writeCode(conn, codes.Ok); err != nil {
		return fmt.Errorf("can't write OK: %w", err)
	}

	return chunks.SendChunk(conn, chk)
}
