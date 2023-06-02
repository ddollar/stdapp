.PHONY: all build lint test vendor

all: build

build:
	make -C web build
	go build -o dist/app -mod=vendor --ldflags="-s -w" .
	go build -o dist/radiusd -mod=vendor --ldflags="-s -w" ./cmd/radiusd

lint:
	make -C web lint
	go run ./vendor/github.com/golangci/golangci-lint/cmd/golangci-lint/main.go run

test:
	make -C web test
	env TEST=true go test -covermode atomic -coverprofile coverage.txt -mod=vendor ./...

vendor:
	go mod tidy
	go mod vendor
	go run vendor/github.com/goware/modvendor/main.go -copy="**/*.c **/*.h"
