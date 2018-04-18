SOURCEDIR=./src
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

VERSION := $(shell git describe --abbrev=0 --tags)
SHA := $(shell git rev-parse --short HEAD)

GOPATH ?= /usr/local/go
GOPATH := ${CURDIR}:${GOPATH}
export GOPATH

all: ./bin/teleport

./bin/teleport: $(SOURCES)
	go build -o ./bin/teleport -ldflags "-X main.BuildVersion=$(VERSION)-$(SHA)" $(SOURCEDIR)/cmd/main.go

.DEFAULT_GOAL: all
include Makefile.git
