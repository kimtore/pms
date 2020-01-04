VERSION := $(shell git describe --always --long --dirty)
DATE := $(shell date +%s)

.PHONY: visp all install test

visp:
	go build -ldflags="-X github.com/ambientsound/pms/version.Version=${VERSION} -X github.com/ambientsound/pms/version.buildDate=${DATE}" -o visp cmd/visp/visp.go

all: test install

install:
	go install -ldflags="-X main.buildVersion=${VERSION}"

test:
	go test ./...
