.PHONY: all build dev clean lint

all: build

build:
	go build -o gonova ./cmd/main.go

dev:
	gonova watchexec -e go -e tmpl go build -o gonova ./cmd/main.go

clean:
	go clean
	rm -f main

lint:
	golangci-lint run ./...
