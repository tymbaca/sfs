package sfs

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/tymbaca/sfs/internal/codes"
)

func (s *Server) handleListIDs(ctx context.Context, conn net.Conn) error {
	req, err := readListIDsReq(conn)
	if err != nil {
		return err
	}

	ids, err := s.storage.ListChunkIDs(ctx, req.name)
	if err != nil {
		writeCodeMsg(conn, codes.Internal, err.Error())
		return err
	}

	return writeListIDsResp(conn, ids)
}

type listIDsReq struct {
	name string
}

func readListIDsReq(r io.Reader) (listIDsReq, error) {
	var filenameSize uint64
	if err := binary.Read(r, binary.LittleEndian, &filenameSize); err != nil {
		return listIDsReq{}, fmt.Errorf("can't read filename size from request: %w", err)
	}

	filename := make([]byte, filenameSize)
	_, err := io.ReadFull(r, filename)
	if err != nil {
		return listIDsReq{}, fmt.Errorf("can't read filename from request: %w", err)
	}

	return listIDsReq{
		name: string(filename),
	}, nil
}

func writeListIDsResp(w io.Writer, ids []uint64) error {
	if err := writeCode(w, codes.Ok); err != nil {
		return fmt.Errorf("can't write OK: %w", err)
	}

	if err := binary.Write(w, binary.LittleEndian, uint64(len(ids))); err != nil {
		return fmt.Errorf("can't write ids len: %w", err)
	}

	for _, id := range ids {
		if err := binary.Write(w, binary.LittleEndian, id); err != nil {
			return fmt.Errorf("can't write chunk id: %w", err)
		}
	}

	return nil
}
