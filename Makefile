#@IgnoreInspection BashAddShebang
export ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))
export CGO_ENABLED=0

export BIN_NAME=api

all: format build test

format:
	gofmt -s -w $(ROOT)

build:
	go build -o $(BIN_NAME) -ldflags="-s -w" $(ROOT)/*.go

run:
	go run $(ROOT)/*.go

test:
	CGO_ENABLED=1 go test -race -coverprofile=coverage.txt -covermode=atomic $(ROOT)/...
