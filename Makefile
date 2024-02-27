.PHONY: all build lint test vendor

all: build

build:
	go generate >/dev/null
	go build -o dist -mod=vendor --ldflags="-s -w" . ./cmd/sa

lint:
	go run ./vendor/github.com/golangci/golangci-lint/cmd/golangci-lint/main.go run

test:
	env TEST=true go test -covermode atomic -coverprofile coverage.txt -mod=vendor ./...

vendor:
	go mod tidy
	go work vendor
