package client

import (
	"context"

	sfs "github.com/tymbaca/sfs/client"
)

const (
	KiB uint64 = 1 << 10
	MiB        = KiB << 10
	GiB        = MiB << 10
)

func main() {
	ctx := context.Background()

	client := sfs.NewClient("localhost:6886", 10)
	client.Upload()
}
