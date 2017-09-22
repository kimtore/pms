VERSION := $(shell git describe --always --long --dirty)

.PHONY: all install test get-deps

all: get-deps test install

install:
	go install -ldflags="-X main.buildVersion=${VERSION}"

test:
	go test ./...

get-deps:
	go get -u github.com/ambientsound/gompd/mpd
	go get -u github.com/blevesearch/bleve
	go get -u github.com/jessevdk/go-flags
	go get -u github.com/gdamore/tcell
	go get -u github.com/nsf/termbox-go
	go get -u github.com/stretchr/testify
