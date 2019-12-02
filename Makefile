VERSION := $(shell git describe --always --long --dirty)

.PHONY: all install test

all: test install

install:
	go install -ldflags="-X main.buildVersion=${VERSION}"

test:
	go test ./...
