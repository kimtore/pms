VERSION := $(shell git describe --always --long --dirty)

.PHONY: visp all install test

visp:
	go build -ldflags="-X main.buildVersion=${VERSION}" -o visp cmd/visp/visp.go

all: test install

install:
	go install -ldflags="-X main.buildVersion=${VERSION}"

test:
	go test ./...
