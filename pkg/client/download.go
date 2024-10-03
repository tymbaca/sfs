package sfs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/tymbaca/sfs/internal/chunks"
	"github.com/tymbaca/sfs/internal/transport"
	"golang.org/x/sync/errgroup"
)

func (c *Client) DownloadTo(ctx context.Context, w io.WriterAt, name string) error {
	// Get id-addr mapping to know where to go for each chunk
	idToAddr, err := c.resolveChunksAddrs(ctx, name)
	if err != nil {
		return fmt.Errorf("can't download the file: can't resolve chunks for '%s': %w", name, err)
	}

	chks := make([]chunks.Chunk, len(idToAddr))
	closes := make([]func() error, len(idToAddr))
	var g errgroup.Group

	// Receive all chunks
	for id, addr := range idToAddr {
		id, addr := id, addr
		g.Go(func() error {
			trans := transport.NewTCPTransport(addr)

			chk, err := trans.RecvChunk(ctx, name, id)

			chks[id] = chk // result will be in order of IDs: 0, 1, 2, etc
			closes[id] = trans.Close
			return err
		})
	}

	err = g.Wait()
	if err != nil {
		return fmt.Errorf("can't download the file: %w", err)
	}

	// defer close all transport connetions after download
	defer func() {
		for _, cls := range closes {
			_ = cls()
		}
	}()

	// Write chunks to the writer
	return chunks.WriteTo(w, chks...)
}

func (c *Client) resolveChunksAddrs(ctx context.Context, name string) (map[uint64]string, error) {
	addrToIDs := make(map[string][]uint64, len(c.addrs))
	var mu sync.Mutex
	var g errgroup.Group

	// Get chunk ids from each address
	for _, addr := range c.addrs {
		addr := addr
		g.Go(func() (err error) {
			trans := transport.NewTCPTransport(addr)

			ids, err := trans.ListIDs(ctx, name)
			if err != nil {
				return err
			}

			mu.Lock()
			addrToIDs[addr] = ids
			mu.Unlock()
			return nil
		})
	}

	err := g.Wait()
	if err != nil {
		return nil, fmt.Errorf("can't detect chunk ids: %w", err)
	}

	// WARN: need to rebalance
	// mapping ids to addrs
	idToAddr := make(map[uint64]string)
	for addr, ids := range addrToIDs {
		for _, id := range ids {
			idToAddr[id] = addr
		}
	}

	if len(idToAddr) == 0 {
		return nil, errors.New("file not found")
	}

	if !isIDToAddrContinuous(idToAddr) {
		return nil, errors.New("file chunks are incomplete")
	}

	return idToAddr, nil
}

func isIDToAddrContinuous(m map[uint64]string) bool {
	// if ids are correct, map must contain ids: 0, 1, 2, 3 ... etc
	for i := range len(m) {
		_, ok := m[uint64(i)]
		if !ok {
			return false
		}
	}

	return true
}
