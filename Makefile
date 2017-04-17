VERSION := $(shell git describe --always --long --dirty)

.PHONY: all

all:
	go install -ldflags="-X main.buildVersion=${VERSION}"
