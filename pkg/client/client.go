package sfs

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spaolacci/murmur3"
	"github.com/tymbaca/sfs/internal/chunks"
	"github.com/tymbaca/sfs/internal/logger"
	"github.com/tymbaca/sfs/internal/transport"
	"github.com/tymbaca/sfs/pkg/chunkio"
	"github.com/tymbaca/sfs/pkg/mem"
)

type Client struct {
	addrs     []string
	chunkSize int64 // bytes
}

func NewClient(addrs string, chunkSize int64) *Client {
	return &Client{
		addrs:     strings.Split(addrs, ","),
		chunkSize: chunkSize,
	}
}

func (c *Client) UploadFile(ctx context.Context, name string, f *os.File) error {
	stat, err := f.Stat()
	if err != nil {
		return fmt.Errorf("can't get file stats: %w", err)
	}

	return c.Upload(ctx, name, f, stat.Size())
}

func (c *Client) Upload(ctx context.Context, name string, r io.ReaderAt, totalSize int64) error {
	chunks, err := formChunks(r, totalSize, name, c.chunkSize)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for chunk := range chunks {
		wg.Add(1)
		go c.uploadChunk(ctx, chunk, &wg)
	}

	wg.Wait()

	return nil
}

func (c *Client) uploadChunk(ctx context.Context, chunk chunks.Chunk, wg *sync.WaitGroup) error {
	defer wg.Done()
	start := time.Now()
	logger.Debugf("starting to upload the %d chunk", chunk.ID)
	defer func() {
		logger.Debugf("uploaded %d chunk, %.2f MiB, time elapsed: %s", chunk.ID, float32(chunk.Size)/float32(mem.MiB), time.Since(start))
	}()

	trans, err := c.getTransport(chunk)
	if err != nil {
		return fmt.Errorf("can't get transport: %w", err)
	}
	defer trans.Close()

	if err = trans.SendChunk(ctx, chunk); err != nil {
		return fmt.Errorf("can't send chunk: %w", err)
	}

	return nil
}

func (c *Client) getTransport(chunk chunks.Chunk) (transport.Transport, error) {
	addr := c.resolveNodeByChunk(chunk.Filename, chunk.ID)
	trans := transport.NewTCPTransport(addr)

	return trans, nil
}

func (c *Client) resolveNodeByChunk(name string, id uint64) string {
	key := []byte(name + fmt.Sprint(id))
	hash := murmur3.Sum32(key)
	logger.Debugf("name '%s', id %d, hash = %d", name, id, hash)

	idx := int(hash) % len(c.addrs)
	return c.addrs[idx]
}

func formChunks(r io.ReaderAt, totalSize int64, name string, size int64) (<-chan chunks.Chunk, error) {
	if size < 1 {
		panic("can't split byte non-positive size")
	}

	chks := chunkio.Split(r, totalSize, size)

	ch := make(chan chunks.Chunk)
	go func() {
		defer close(ch)
		for id, chnk := range chks {
			ch <- chunks.Chunk{
				ID:       uint64(id),
				Filename: name,
				Size:     uint64(chnk.Size()),
				Body:     chnk,
			}
		}
	}()

	return ch, nil
}

func closeConns(conns []net.Conn) {
	for _, conn := range conns {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}
}
