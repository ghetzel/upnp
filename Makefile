.PHONY: all deps fmt bin/upnpc
.EXPORT_ALL_VARIABLES:

GO111MODULE ?= on
BIN         ?= bin/upnpc-$(shell go env GOOS)-$(shell go env GOARCH)

all: deps fmt $(BIN)

deps:
	go get ./...
	go mod tidy

fmt:
	go fmt ./...

$(BIN):
	go build -o $(BIN) cmd/*.go
