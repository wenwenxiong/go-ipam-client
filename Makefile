.ONESHELL:
SHA := $(shell git rev-parse --short=8 HEAD)
GITVERSION := $(shell git describe --long --all)
BUILDDATE := $(shell date -Iseconds)
VERSION := $(or ${VERSION},devel)

CGO_ENABLED := $(or ${CGO_ENABLED},0)
GO := go
GO111MODULE := on
LINKMODE := -extldflags '-static -s -w'

.EXPORT_ALL_VARIABLES:

all:   client

.PHONY: client
client:
	go build -tags netgo,osusergo,urfave_cli_no_docs \
		 -ldflags "$(LINKMODE) -X 'github.com/metal-stack/v.Version=$(VERSION)' \
								   -X 'github.com/metal-stack/v.Revision=$(GITVERSION)' \
								   -X 'github.com/metal-stack/v.GitSHA1=$(SHA)' \
								   -X 'github.com/metal-stack/v.BuildDate=$(BUILDDATE)'" \
	   -o bin/go-ipam-client cmd/go-ipam-client.go
	strip bin/go-ipam-client
