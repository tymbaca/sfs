package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	sfs "github.com/tymbaca/sfs/client"
	"github.com/tymbaca/sfs/internal/chunk"
)

const (
	KiB uint64 = 1 << 10
	MiB        = KiB << 10
	GiB        = MiB << 10
)

func main() {
	ctx := context.Background()

	go func() {
		lis, err := net.Listen("tcp", ":6886")
		if err != nil {
			panic(err)
		}

		for {
			conn, err := lis.Accept()
			if err != nil {
				panic(err)
			}

			for {
				chunk, err := chunk.ReadChunk(conn)
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					panic(err)
				}

				fmt.Printf("server: got chunk: %s\n", chunk)
			}
		}
	}()

	f, err := os.Open("cmd/client/main.go")
	if err != nil {
		panic(err)
	}

	client := sfs.NewClient("localhost:6886", 256)
	err = client.Upload(ctx, "man/file", f)
	if err != nil {
		panic(err)
	}

	time.Sleep(500 * time.Millisecond)
}
