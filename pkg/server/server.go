package sfs

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/tymbaca/sfs/internal/codes"
	"github.com/tymbaca/sfs/internal/logger"
)

type Server struct {
	addr    string
	storage storage
}

func New(addr string, storage storage) *Server {
	return &Server{
		addr:    addr,
		storage: storage,
	}
}

func (s *Server) Addr() string {
	return s.addr
}

func (s *Server) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("can't listen addr '%s': %w", s.addr, err)
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			panic(err)
		}

		go func() {
			err := s.handleConn(ctx, conn)
			if err != nil {
				logger.Logf("can't handle conn: %s", err)
				return
			}
		}()
	}
}

func (s *Server) handleConn(ctx context.Context, conn net.Conn) error {
	defer conn.Close()

	head, err := peekByte(conn)
	if err != nil {
		logger.Logf("can't read the head of request: %s", err)
		return err
	}

	start := time.Now()
	traceID := uuid.New()
	defer func() {
		logger.Logf("request '%c', trace-id '%s', time elapsed: %s", head, traceID, time.Since(start))
	}()

	switch head {
	case '*':
		return s.handleSendChunk(ctx, conn)
	case '/':
		return s.handleRecvChunk(ctx, conn)
	case '%':
		return s.handleListIDs(ctx, conn)
	}

	return writeCodeMsg(conn, codes.InvalidReq, fmt.Sprintf("incorrect head character: '%c' (dec:%d)", head, head))
}

func peekByte(r io.Reader) (byte, error) {
	p := make([]byte, 1)
	_, err := r.Read(p)
	if err != nil {
		return 0, fmt.Errorf("can't peek first byte: %w", err)
	}

	return p[0], nil
}

func writeCode(w io.Writer, code uint64) error {
	if err := binary.Write(w, binary.LittleEndian, code); err != nil {
		return err
	}

	return nil
}

func writeCodeMsg(w io.Writer, code uint64, msg string) error {
	if err := binary.Write(w, binary.LittleEndian, code); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, uint64(len(msg))); err != nil {
		return err
	}

	_, err := w.Write([]byte(msg))
	if err != nil {
		return err
	}

	return nil
}
