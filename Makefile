VERSION := $(shell git describe --always --long --dirty)

.PHONY: all test

all:
	go install -ldflags="-X main.buildVersion=${VERSION}"

test:
	go test ./...
