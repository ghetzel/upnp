.PHONY: all deps fmt bin/upnpc
.EXPORT_ALL_VARIABLES:

GO111MODULE ?= on

all: deps fmt bin/upnpc

deps:
	go get ./...
	go mod tidy

fmt:
	go fmt ./...

bin/upnpc:
	go build -o bin/upnpc cmd/*.go
