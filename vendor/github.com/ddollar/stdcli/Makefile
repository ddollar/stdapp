.PHONY: all build lint test

all: build

build:
	go build .

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run

test:
	env go test -covermode atomic -coverprofile coverage.txt ./...