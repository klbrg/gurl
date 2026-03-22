BINARY=gurl

.PHONY: build test install clean

build:
	go build -o $(BINARY) .

test:
	go test -cover ./...

install:
	go install .

clean:
	rm -f $(BINARY)
