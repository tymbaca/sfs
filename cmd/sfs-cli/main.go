package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/tymbaca/sfs/internal/files"
	sfs "github.com/tymbaca/sfs/pkg/client"
	"github.com/tymbaca/sfs/pkg/mem"
)

const addrsEnv = "SFS_ADDRS"

func main() {
	ctx := context.Background()
	addrs := os.Getenv(addrsEnv)
	if strings.TrimSpace(addrs) == "" {
		fmt.Printf("set server nodes addresses in env var %s\n", addrsEnv)
		os.Exit(1)
	}

	client := sfs.NewClient(addrs, 2*mem.MiB)

	if len(os.Args) < 2 {
		fmt.Println("specify the operation")
		os.Exit(1)
	}

	op := os.Args[1]

	switch op {
	case "upload":
		if len(os.Args) < 3 {
			fmt.Println("specify the input file")
			os.Exit(1)
		}
		pathToFile := os.Args[2]

		f, err := os.Open(pathToFile)
		if err != nil {
			fmt.Printf("can't open the file: %s\n", err)
			os.Exit(1)
		}

		err = client.UploadFile(ctx, path.Base(pathToFile), f)
		if err != nil {
			fmt.Printf("error while uploading: %s\n", err)
			os.Exit(1)
		}
	case "download":
		if len(os.Args) < 4 {
			fmt.Println("specify the target filename and destination path")
			os.Exit(1)
		}
		name := os.Args[2]
		dstPath := os.Args[3]

		dst, err := files.CreateFile(dstPath)
		if err != nil {
			fmt.Printf("error while creating destination file: %s\n", err)
			os.Exit(1)
		}

		src, cls, expectedSize, err := client.Download(ctx, name)
		if err != nil {
			fmt.Printf("error while downloading: %s\n", err)
			os.Exit(1)
		}
		defer cls()

		actualSize, err := io.Copy(dst, src)
		if err != nil {
			fmt.Printf("error while downloading: %s\n", err)
			os.Exit(1)
		}

		if expectedSize != actualSize {
			fmt.Printf("error while downloading: final size dismatch: expected %d, got %d %s\n", expectedSize, actualSize, err)
			os.Exit(1)
		}

	default:
		fmt.Println("unknown operation")
	}
}
