PKG := github.com/the-lightning-land/sweetd

GO_BIN := ${GOPATH}/bin
VERSION := $(shell git describe --tags)
COMMIT := $(shell git rev-parse HEAD)
DATE := $(shell date +%Y-%m-%d)

LDFLAGS := "-X main.Commit=$(COMMIT) -X main.Version=$(VERSION) -X main.Date=$(DATE) ${LDFLAGS}"

GOBUILD := GO111MODULE=on go build -v
RM := rm -f

# commands

default: build

compile:
	@$(call print, "Building sweetd.")
	$(GOBUILD) -o sweetd -ldflags $(LDFLAGS) $(PKG)

test:
	@$(call print, "Testing sweetd.")
	go test -v ./...

clean:
	@$(call print, "Cleaning builds and module cache")
	$(RM) ./sweetd

clean-cache:
	@$(call print, "Cleaning go module cache")
	go clean --modcache

build: compile

