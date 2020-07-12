#@IgnoreInspection BashAddShebang
export ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))
export CGO_ENABLED=0

export BIN_NAME=api

all: format build test

format:
	gofmt -s -w $(ROOT)

build:
	go build -o $(ROOT)/bin/$(BIN_NAME) -ldflags="-s -w" $(ROOT)/cmd/realworld/*.go

run:
	go run $(ROOT)/cmd/realworld/*.go

test:
	CGO_ENABLED=1 go test -race -coverprofile=coverage.txt -covermode=atomic $(ROOT)/...
