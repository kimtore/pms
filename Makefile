VERSION := $(shell git describe --always --long --dirty)
DATE := $(shell date +%s)
LDFLAGS := -ldflags="-X main.buildVersion=${VERSION}"

.PHONY: install pms test linux-amd64 linux-arm64 linux-arm darwin-amd64 darwin-arm64 windows-amd64.exe

install: pms
	sh ./install.sh

pms:
	go build ${LDFLAGS} -o build/pms main.go

test:
	go test ./...

linux-amd64:
	GOOS=linux GOARCH=amd64 \
	go build ${LDFLAGS} -o build/pms-linux-amd64 main.go

linux-arm64:
	GOOS=linux GOARCH=arm64 \
	go build ${LDFLAGS} -o build/pms-linux-arm64 main.go

linux-arm:
	GOOS=linux GOARCH=arm \
	go build ${LDFLAGS} -o build/pms-linux-arm main.go

darwin-amd64:
	GOOS=darwin GOARCH=amd64 \
	go build ${LDFLAGS} -o build/pms-darwin-amd64 main.go

darwin-arm64:
	GOOS=darwin GOARCH=arm64 \
	go build ${LDFLAGS} -o build/pms-darwin-arm64 main.go

windows-amd64.exe:
	GOOS=windows GOARCH=amd64 \
	go build ${LDFLAGS} -o build/pms-windows-amd64.exe main.go
