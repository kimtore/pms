VERSION := $(shell git describe --always --long --dirty)

.PHONY: all test

all:
	go install -ldflags="-X main.buildVersion=${VERSION}"

test:
	go test ./...

get-deps:
	go get github.com/ambientsound/gompd/mpd
	go get github.com/blevesearch/bleve
	go get github.com/gdamore/tcell
	go get github.com/jessevdk/go-flags
