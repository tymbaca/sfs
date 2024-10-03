client:
	go run ./cmd/full

server:
	go run ./cmd/server

cli:
	go build -o bin/sfs-cli ./cmd/sfs-cli

.PHONY: client server
