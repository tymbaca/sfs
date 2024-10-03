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
	"github.com/tymbaca/sfs/pkg/mem"
	sfs_server "github.com/tymbaca/sfs/pkg/server"
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
	storage1 := storage.NewFileStorage("cmd/output/server/1st-node")
	storage2 := storage.NewFileStorage("cmd/output/server/2nd-node")
	storage3 := storage.NewFileStorage("cmd/output/server/3rd-node")

	server1 := sfs_server.New(":6886", storage1)
	go func() {
		log.Fatal(server1.Run(ctx))
	}()

	server2 := sfs_server.New(":6887", storage2)
	go func() {
		log.Fatal(server2.Run(ctx))
	}()

	server3 := sfs_server.New(":6888", storage3)
	go func() {
		log.Fatal(server3.Run(ctx))
	}()

	//--------------------------------------------------------------------------------------------------

	f1, err := os.Open("cmd/input/odin-macos-arm64-dev-2024-09.zip")
	if err != nil {
		panic(err)
	}
	f2, err := os.Open("cmd/input/odin-macos-arm64-dev-2024-10.zip")
	if err != nil {
		panic(err)
	}

	client := sfs_client.NewClient("localhost:6886,localhost:6887,localhost:6888", 8*mem.MiB)

	// UploadFile
	err = client.UploadFile(ctx, path.Base(f1.Name()), f1)
	if err != nil {
		panic(err)
	}

	err = client.UploadFile(ctx, path.Base(f2.Name()), f2)
	if err != nil {
		panic(err)
	}

	// Download
	downloadAndSave(ctx, client, path.Base(f1.Name()))
	downloadAndSave(ctx, client, path.Base(f2.Name()))
}

func downloadAndSave(ctx context.Context, client *sfs_client.Client, name string) {
	r, cls, size, err := client.Download(ctx, name)
	if err != nil {
		panic(err)
	}
	defer cls()
	fmt.Println(size)
	pth := path.Join("cmd/output/client", name)

	out, err := os.Create(pth)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(out, r)
	if err != nil {
		panic(err)
	}
}

func getFileSize(f *os.File) int64 {
	stat, err := f.Stat()
	if err != nil {
		panic(err)
	}

	return stat.Size()
}
