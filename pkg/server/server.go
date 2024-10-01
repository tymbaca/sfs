package sfs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/tymbaca/sfs/internal/chunks"
	"github.com/tymbaca/sfs/internal/logger"
)

type Server struct {
	addr     string
	connPool chan struct{}
	storage  storage
}

func New(addr string, storage storage) *Server {
	connPoolSize := 10
	return &Server{
		addr:     addr,
		connPool: make(chan struct{}, connPoolSize),
		storage:  storage,
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

		s.queue()
		go s.handleConn(ctx, conn)
	}
}

func (s *Server) queue() {
	s.connPool <- struct{}{}
}

func (s *Server) unqueue() {
	<-s.connPool
}

func (s *Server) handleConn(ctx context.Context, conn net.Conn) {
	logger.Log("handling conn")
	defer s.unqueue()
	for {
		logger.Log("handling conn iter")
		// TODO add timeout
		chunk, err := chunks.RecvChunk(conn)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}

		if err := s.storage.StoreChunk(ctx, chunk); err != nil {
			logger.Logf("ERROR: can't store chunk: %s", err)
		}
	}
}
