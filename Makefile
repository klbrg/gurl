BINARY=gurl

.PHONY: build clean

build:
	go build -o $(BINARY) .

clean:
	rm -f $(BINARY)
