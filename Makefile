.PHONY: all lint test

all: lint

lint:
	go run ./vendor/github.com/golangci/golangci-lint/cmd/golangci-lint/main.go run

test:
	env TEST=true go test -covermode atomic -coverprofile coverage.txt -mod=vendor ./...
