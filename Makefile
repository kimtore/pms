VERSION := $(shell git describe --always --long --dirty)

.PHONY: all install test get-deps

all: get-deps test install

install:
	go install -ldflags="-X main.buildVersion=${VERSION}"

test:
	go test ./...

get-deps:
	go get github.com/ambientsound/gompd/mpd
	go get github.com/blevesearch/bleve
	go get github.com/gdamore/tcell
	go get github.com/jessevdk/go-flags
	go get github.com/stretchr/testify
