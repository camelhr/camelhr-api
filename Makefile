.PHONY: build test

build:
	go build -o bin/ ./...

test:
	go test -v -cover ./...