package sfs

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/tymbaca/sfs/internal/chunks"
	"github.com/tymbaca/sfs/internal/transport"
	"github.com/tymbaca/sfs/pkg/chunkio"
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
	addrIdx := c.resolveNodeIndex(chunk.Filename, chunk.ID)
	trans := transport.NewTCPTransport(c.addrs[addrIdx])

	return trans, nil
}

func (c *Client) resolveNodeIndex(name string, id uint64) int {
	// TODO i'm too lazy for this shit
	// but i need consistent hashing
	// sum := sha1.Sum([]byte(name + fmt.Sprint(id))) // not good
	// i := new(big.Int)
	// i.SetBytes()

	return rand.Intn(len(c.addrs))
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
