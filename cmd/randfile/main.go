package main

import (
	"crypto/rand"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	path := os.Args[1]
	sizeStr := os.Args[2]
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.CopyN(f, rand.Reader, int64(size))
	if err != nil {
		log.Fatal(err)
	}
}
