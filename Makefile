.PHONY: tools test lint

all: test

lint:
	golangci-lint run --tests=false

test:
	go test -v -race -cover ./...

# tools:
# 	go install github.com/golangci/golangci-lint/cmd/golangci-lint

