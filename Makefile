GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BASEPATH := $(shell pwd)
BUILDDIR=$(BASEPATH)/dist/bin
MAIN= $(BASEPATH)/main.go

APP_NAME=ekko

GOPROXY="https://goproxy.cn,direct"

build_darwin:
	GOOS=darwin GOARCH=amd64  $(GOBUILD) -o $(BUILDDIR)/darwin/amd64/$(APP_NAME) $(MAIN)
