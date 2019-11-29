GOPATH = $(shell pwd)

.PHONY: all
all: build

.PHONY: build
build:
	GOPATH=$(GOPATH) go build \
		   -o bin/go-evt \
		   go-evt
