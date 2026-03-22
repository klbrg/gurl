BINARY=gurl

.PHONY: build test clean

build:
	go build -o $(BINARY) .

test:
	go test -cover ./...

clean:
	rm -f $(BINARY)
