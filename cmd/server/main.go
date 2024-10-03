package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tymbaca/sfs/internal/storage"
	sfs "github.com/tymbaca/sfs/pkg/server"
)

func main() {
	ctx := context.Background()

	// logStorage := logStorage{}
	storage1 := storage.NewFileStorage("cmd/output/server/1st-node")
	storage2 := storage.NewFileStorage("cmd/output/server/2nd-node")
	storage3 := storage.NewFileStorage("cmd/output/server/3rd-node")

	server1 := sfs.New(":6886", storage1)
	go func() {
		log.Fatal(server1.Run(ctx))
	}()

	server2 := sfs.New(":6887", storage2)
	go func() {
		log.Fatal(server2.Run(ctx))
	}()

	server3 := sfs.New(":6888", storage3)
	go func() {
		log.Fatal(server3.Run(ctx))
	}()

	fmt.Println("started nodes on addrs:", server1.Addr(), server2.Addr(), server3.Addr())
	<-make(chan struct{})
}
