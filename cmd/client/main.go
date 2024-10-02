package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"

	"github.com/tymbaca/sfs/internal/chunks"
	"github.com/tymbaca/sfs/internal/logger"
	"github.com/tymbaca/sfs/internal/storage"
	sfs_client "github.com/tymbaca/sfs/pkg/client"
	sfs_server "github.com/tymbaca/sfs/pkg/server"
)

const (
	KiB int64 = 1 << 10
	MiB       = KiB << 10
	GiB       = MiB << 10
)

type logStorage struct{}

func (s logStorage) StoreChunk(ctx context.Context, chunk chunks.Chunk) error {
	logger.Logf("server: chunk: %s\n", chunk)
	time.Sleep(100 * time.Millisecond)
	return nil
}

func main() {
	ctx := context.Background()

	// logStorage := logStorage{}
	storage := storage.NewFileStorage("cmd/output/server")
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

	client := sfs_client.NewClient("localhost:6886", 64*MiB)
	err = client.UploadFile(ctx, path.Base(f1.Name()), f1)
	if err != nil {
		panic(err)
	}
	// err = client.Upload(ctx, "small2", f2)
	// if err != nil {
	// 	panic(err)
	// }

	r, cls, size, err := client.Download(ctx, path.Base(f1.Name()))
	if err != nil {
		panic(err)
	}
	defer cls()

	fmt.Println(size)

	out, err := os.Create(path.Join("cmd/output/client", path.Base(f1.Name())))
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(out, r)
	if err != nil {
		panic(err)
	}

	time.Sleep(1000 * time.Millisecond)
}

func getFileSize(f *os.File) int64 {
	stat, err := f.Stat()
	if err != nil {
		panic(err)
	}

	return stat.Size()
}
