client:
	go run ./cmd/full

server:
	go run ./cmd/server

stress:
	go run ./cmd/stress

cli:
	go build -o bin/sfs-cli ./cmd/sfs-cli

randfile:
	go build -o bin/randfile ./cmd/randfile

.PHONY: client server
