package transport

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/tymbaca/sfs/internal/chunks"
	"github.com/tymbaca/sfs/internal/codes"
	"github.com/tymbaca/sfs/internal/logger"
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
	if err := t.ensureDial(); err != nil {
		return nil, err
	}

	if _, err := t.conn.Write([]byte("%")); err != nil {
		return nil, err
	}

	// we need len of bytes, not len of utf-8 symbols, so we use [len]
	if err := binary.Write(t.conn, binary.LittleEndian, uint64(len(name))); err != nil {
		return nil, fmt.Errorf("can't write filename size: %w", err)
	}

	if _, err := t.conn.Write([]byte(name)); err != nil {
		return nil, fmt.Errorf("can't write filename: %w", err)
	}

	code, err := readCode(t.conn)
	if err != nil {
		return nil, fmt.Errorf("can't read the code: %w", err)
	}

	switch code {
	case codes.Ok:
		// read the count of ids
		var count uint64
		if err := binary.Read(t.conn, binary.LittleEndian, &count); err != nil {
			return nil, fmt.Errorf("can't read the ids count: %w", err)
		}

		// read the ids
		ids := make([]uint64, 0, count)
		for i := range count {
			var id uint64
			if err := binary.Read(t.conn, binary.LittleEndian, &id); err != nil {
				return nil, fmt.Errorf("can't read the #%d id: %w", i, err)
			}

			ids = append(ids, id)
		}

		return ids, nil

	case codes.Internal:
		msg, err := readMsg(t.conn)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("got error from server: %s", msg)
	}

	return nil, fmt.Errorf("list ids: unsupported response code: %d", code)
}

func (t *TCPTransport) RecvChunk(ctx context.Context, name string, id uint64) (chunks.Chunk, error) {
	// TODO should we just create conn for each trans endpoint and defer close it?
	if err := t.ensureDial(); err != nil {
		return chunks.Chunk{}, err
	}

	if _, err := t.conn.Write([]byte("/")); err != nil {
		return chunks.Chunk{}, err
	}

	// we need len of bytes, not len of utf-8 symbols, so we use [len]
	if err := binary.Write(t.conn, binary.LittleEndian, uint64(len(name))); err != nil {
		return chunks.Chunk{}, fmt.Errorf("can't write filename size: %w", err)
	}

	if _, err := t.conn.Write([]byte(name)); err != nil {
		return chunks.Chunk{}, fmt.Errorf("can't write filename: %w", err)
	}

	if err := binary.Write(t.conn, binary.LittleEndian, id); err != nil {
		return chunks.Chunk{}, fmt.Errorf("can't write chunk ID: %w", err)
	}

	// Starting to read response
	code, err := readCode(t.conn)
	if err != nil {
		return chunks.Chunk{}, fmt.Errorf("can't read the code: %w", err)
	}

	logger.Debugf("got code server: %d", code)

	switch code {
	case codes.Ok:
		chk, err := chunks.RecvChunk(t.conn)
		if err != nil {
			return chunks.Chunk{}, fmt.Errorf("can't receive chunk from server: %w", err)
		}

		return chk, nil

	case codes.NotFound:
		return chunks.Chunk{}, fmt.Errorf("file not found: %d", code)

	case codes.Internal, codes.InvalidReq:
		msg, err := readMsg(t.conn)
		if err != nil {
			return chunks.Chunk{}, err
		}

		return chunks.Chunk{}, fmt.Errorf("got error from server, code %d, msg: %s", code, msg)
	}
	return chunks.Chunk{}, fmt.Errorf("recv chunk: unsupported response code: %d", code)
}

func (t *TCPTransport) Close() error {
	if t.conn == nil {
		return nil
	}

	conn := t.conn
	t.conn = nil
	return conn.Close()
}

func readCode(r io.Reader) (code codes.Code, err error) {
	err = binary.Read(r, binary.LittleEndian, &code)
	return
}

func readMsg(r io.Reader) (string, error) {
	var size uint64
	if err := binary.Read(r, binary.LittleEndian, &size); err != nil {
		return "", fmt.Errorf("can't read the message: %w", err)
	}

	msgBuf := make([]byte, size)
	if _, err := io.ReadFull(r, msgBuf); err != nil {
		return "", fmt.Errorf("can't read the message: %w", err)
	}

	return string(msgBuf), nil
}
