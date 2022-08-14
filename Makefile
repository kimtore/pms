VERSION := $(shell git describe --always --long --dirty)

.PHONY: install pms test

pms:
	mkdir -p build/
	go build -o build/pms -ldflags="-X main.buildVersion=${VERSION}" main.go

install:
	go install -ldflags="-X main.buildVersion=${VERSION}"

test:
	go test ./...
