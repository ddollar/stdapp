.PHONY: all build clean release

gopath := $(shell mktemp -d)
work   := $(gopath)/src/github.com/ddollar/stdapp/cmd/sa

all: build

build:
	tar cz -C ../.. . | docker run -v .:/src -i ubuntu tar xz -C /sa
	mkdir -p $(work) && cd $(work) && env GOPATH=$(gopath) $(shell go env GOPATH)/bin/xgo -out pkg/sa -targets 'darwin/amd64,linux/amd64' -ldflags "-X main.version=$(VERSION) -X main.image=$(IMAGE)" .
	mkdir -p pkg && docker run -v $(gopath):/gopath -i ubuntu tar czv /gopath/src/github.com/sa/sa/cmd/sa/pkg | tar xzv -C pkg --strip-components 8
	mkdir -p $(work) && cd $(work) && env GOPATH=$(gopath) $(shell go env GOPATH)/bin/xgo -out pkg/sa -targets 'darwin/amd64,linux/amd64' -ldflags "-X main.version=$(VERSION) -X main.image=$(IMAGE)" .

clean:
	rm -f pkg/sa-*%
