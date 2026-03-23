BINARY=gurl
LDFLAGS=-ldflags "-X main.commit=$(shell git rev-parse --short HEAD)"

.PHONY: build test install clean

build:
	go build $(LDFLAGS) -o $(BINARY) .

test:
	go test -cover ./...

install:
	go install $(LDFLAGS) .

clean:
	rm -f $(BINARY)
