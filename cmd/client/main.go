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
)

const (
	KiB uint64 = 1 << 10
	MiB        = KiB << 10
	GiB        = MiB << 10
)

func main() {
	ctx := context.Background()

	go func() {
		lis, err := net.Listen("tcp", "localhost:6886")
		if err != nil {
			panic(err)
		}

		for {
			conn, err := lis.Accept()
			if err != nil {
				panic(err)
			}

			for {
				chunk, err := sfs.ReadChunk(conn)
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					panic(err)
				}

				fmt.Printf("server: got chunk: %v\n", chunk)
			}
		}
	}()

	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}

	client := sfs.NewClient("localhost:6886", 1*KiB)
	err = client.Upload(ctx, "shit.txt", f)
	if err != nil {
		panic(err)
	}

	time.Sleep(500 * time.Millisecond)
}
