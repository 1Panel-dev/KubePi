GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BASEPATH := $(shell pwd)
BUILDDIR=$(BASEPATH)/dist/bin
UIDIR=$(BASEPATH)/web/dashboard
TERMINALDIR=$(BASEPATH)/web/terminal

MAIN= $(BASEPATH)/cmd/server/main.go

APP_NAME=ekko

GOPROXY="https://goproxy.cn,direct"

build_web:
	cd $(UIDIR) && npm install && npm run-script build
	cd $(TERMINALDIR) && npm install && npm run-script build


build_bin_darwin: build_web
	cd $(BASEPATH)
	GOOS=darwin GOARCH=amd64  $(GOBUILD) -o $(BUILDDIR)/darwin/amd64/$(APP_NAME) $(MAIN)
