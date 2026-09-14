.PHONY: build test check install clean

BINARY_NAME=mcp-sonar-gitlab

build:
	go build -o bin/$(BINARY_NAME) .

test:
	go test -v -count=1 ./...

check:
	go vet ./...
	go test -race -count=1 ./...

install:
	go install .

clean:
	rm -rf bin/
