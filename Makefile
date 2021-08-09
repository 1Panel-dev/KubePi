GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOARCH=$(shell go env GOARCH)
GOOS=$(shell go env GOOS )

BASEPATH := $(shell pwd)
BUILDDIR=$(BASEPATH)/dist/bin
EKKODIR=$(BASEPATH)/web/ekko
DASHBOARDDIR=$(BASEPATH)/web/dashboard
TERMINALDIR=$(BASEPATH)/web/terminal


MAIN= $(BASEPATH)/cmd/server/main.go

APP_NAME=ekko

GOPROXY="https://goproxy.cn,direct"

build_web_ekko:
	cd $(EKKODIR) && npm install && npm run-script build
build_web_dashboard:
	cd $(DASHBOARDDIR) && npm install && npm run-script build
build_web_terminal:
	cd $(TERMINALDIR) && npm install && npm run-script build

build_web: build_web_ekko build_web_dashboard build_web_terminal

build_bin:
	GOOS=$(GOOS) GOARCH=$(GOARCH)  $(GOBUILD) -trimpath  -ldflags "-s -w"  -o $(BUILDDIR)/$(GOOS)/$(GOARCH)/$(APP_NAME) $(MAIN)
