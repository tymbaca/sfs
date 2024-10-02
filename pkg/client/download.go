package sfs

import (
	"context"
	"fmt"
	"io"

	"github.com/tymbaca/sfs/internal/chunks"
	"github.com/tymbaca/sfs/internal/logger"
	"github.com/tymbaca/sfs/internal/transport"
	"go.uber.org/multierr"
	"golang.org/x/sync/errgroup"
)

func (c *Client) Download(ctx context.Context, name string) (io.Reader, func() error, int64, error) {
	// Get id-addr mapping to know where to go for each chunk
	idToAddr, err := c.resolveChunksAddrs(ctx, name)
	if err != nil {
		return nil, nil, 0, err
	}

	logger.Debugf("idToAddr: %v", idToAddr)

	// Receive all chunks
	chunks := make([]chunks.Chunk, len(idToAddr))
	closes := make([]func() error, len(idToAddr))
	var g errgroup.Group
	for id, addr := range idToAddr {
		id, addr := id, addr
		g.Go(func() error {
			logger.Debugf("starting id = %d, addr = %s", id, addr)
			defer logger.Debugf("exiting id = %d, addr = %s", id, addr)

			trans := transport.NewTCPTransport(addr)

			chk, err := trans.RecvChunk(ctx, name, id)
			logger.Debugf("got chunk header, id = %d, addr = %s", id, addr)

			chunks[id] = chk // result will be in order of IDs: 0, 1, 2, etc
			closes[id] = trans.Close
			return err
		})
	}

	err = g.Wait()
	logger.Debugf("chunks: %v", chunks)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("can't download the file: %w", err)
	}

	// Merge readers
	readers := make([]io.Reader, 0, len(chunks))
	size := int64(0)
	for _, chk := range chunks {
		size += int64(chk.Size)
		readers = append(readers, chk.Body)
	}

	mergedReader := io.MultiReader(readers...)
	closeFn := func() (err error) {
		for _, cls := range closes {
			multierr.Append(err, cls())
		}
		return err
	}

	return mergedReader, closeFn, size, nil
}

func (c *Client) resolveChunksAddrs(ctx context.Context, name string) (map[uint64]string, error) {
	addrToIDs := make(map[string][]uint64, len(c.addrs))
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

			addrToIDs[addr] = ids
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

	return idToAddr, nil
}
