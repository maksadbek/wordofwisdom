all: build-server build-client


build-server:
	go build -o build/server ./cmd/server

build-client:
	go build -o build/client ./cmd/client
