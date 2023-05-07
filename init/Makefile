.PHONY: all build compress install lint test vendor

all: build

build:
	go build -o dist/app -mod=vendor --ldflags="-s -w" .

compress:
	upx-ucl -1 dist/app

install:
	go install .

lint:
	go run ./vendor/github.com/golangci/golangci-lint/cmd/golangci-lint/main.go run

test:
	env TEST=true go test -covermode atomic -coverprofile coverage.txt -mod=vendor ./...

vendor:
	go mod tidy
	go mod vendor
	go run vendor/github.com/goware/modvendor/main.go -copy="**/*.c **/*.h"
