package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	sfs_client "github.com/tymbaca/sfs/client"
	"github.com/tymbaca/sfs/internal/chunk"
	"github.com/tymbaca/sfs/internal/storage"
	sfs_server "github.com/tymbaca/sfs/server"
)

const (
	KiB uint64 = 1 << 10
	MiB        = KiB << 10
	GiB        = MiB << 10
)

type logStorage struct {
}

func (s logStorage) StoreChunk(ctx context.Context, chunk chunk.Chunk) error {
	fmt.Printf("server: chunk: %s\n", chunk)
	time.Sleep(100 * time.Millisecond)
	return nil
}

func main() {
	ctx := context.Background()

	// logStorage := logStorage{}
	storage := storage.NewFileStorage("cmd/output/data")
	server := sfs_server.New(":6886", storage)
	go func() {
		log.Fatal(server.Run())
	}()

	f1, err := os.Open("cmd/input/odin-macos-arm64-dev-2024-09.zip")
	if err != nil {
		panic(err)
	}
	// f2, err := os.Open("cmd/input/small2.txt")
	// if err != nil {
	// 	panic(err)
	// }

	client := sfs_client.NewClient("localhost:6886", 32*MiB)
	err = client.Upload(ctx, "small1", f1)
	if err != nil {
		panic(err)
	}
	// err = client.Upload(ctx, "small2", f2)
	// if err != nil {
	// 	panic(err)
	// }

	time.Sleep(1000 * time.Millisecond)
}
