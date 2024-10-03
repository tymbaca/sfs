package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/tymbaca/sfs/internal/files"
	"github.com/tymbaca/sfs/internal/logger"
	sfs_client "github.com/tymbaca/sfs/pkg/client"
	"github.com/tymbaca/sfs/pkg/mem"
)

func main() {
	ctx := context.Background()

	f, err := os.Open("cmd/input/random-8gb")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	client := sfs_client.NewClient("localhost:6886,localhost:6887,localhost:6888", 64*mem.MiB)

	for i := range 5 {
		go worker(ctx, i, client, f)
	}

	<-make(chan struct{})
}

func worker(ctx context.Context, workerID int, client *sfs_client.Client, f *os.File) {
	i := 0
	for {
		i++
		start := time.Now()
		logger.Debugf("starting %d iteration", i)

		name := path.Join(fmt.Sprint(workerID), path.Base(f.Name()))

		// UploadFile
		err := client.UploadFile(ctx, name, f)
		if err != nil {
			panic(err)
		}

		// Download
		downloadAndSave(ctx, client, name)

		logger.Debugf("iteration %d ended, time elapsed: %s", i, time.Since(start))
	}
}

func downloadAndSave(ctx context.Context, client *sfs_client.Client, name string) {
	r, cls, size, err := client.Download(ctx, name)
	if err != nil {
		panic(err)
	}
	defer cls()
	fmt.Println(size)
	pth := path.Join("cmd/output/client", name)

	out, err := files.CreateFile(pth)
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
