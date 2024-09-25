package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/tymbaca/sfs/internal/chunk"
)

type Server struct {
	addr string
	// connPoolSize uint64 // TODO
	// connPool     uint64
	storage storage
}

func New(addr string) *Server {
	return &Server{
		addr: addr,
		// connPoolSize: 10,
		// connPool:     0,
	}
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("can't listen addr '%s': %w", s.addr, err)
	}

	for {
		ctx := context.Background()
		conn, err := lis.Accept()
		if err != nil {
			panic(err)
		}

		go s.handleConn(ctx, conn)
	}
}

func (s *Server) handleConn(ctx context.Context, conn net.Conn) {
	for {
		// TODO add timeout
		chunk, err := chunk.ReadChunk(conn)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}

		if err := s.storage.StoreChunk(ctx, chunk); err != nil {
			log.Printf("ERROR: can't store chunk: %s", err)
		}
	}
}
