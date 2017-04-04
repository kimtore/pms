VERSION := $(shell git describe --always --long --dirty)

.PHONY: all

all:
	go install -ldflags="-X main.build_version=${VERSION}"
