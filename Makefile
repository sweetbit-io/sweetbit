PKG := github.com/the-lightning-land/sweetd

GO_BIN := ${GOPATH}/bin
VERSION := $(shell git describe --tags)
COMMIT := $(shell git rev-parse HEAD)
DATE := $(shell date +%Y-%m-%d)

LDFLAGS := "-X main.Commit=$(COMMIT) -X main.Version=$(VERSION) -X main.Date=$(DATE)"

GOBUILD := go build
RM := rm -f

PACKR2_PKG := github.com/gobuffalo/packr/v2
PACKR2_BIN := $(GO_BIN)/packr2
PACKR2_COMMIT := $(shell cat go.mod | \
    grep $(PACKR2_PKG) | \
    tail -n1 | \
    awk -F " " '{ print $$2 }' | \
    awk -F "/" '{ print $$1 }')

# commands

default: build

packr2:
	@$(call print, "Installing packr2.")
	go get $(PACKR2_PKG)@$(PACKR2_COMMIT)
	go install $(PACKR2_PKG)/packr2

pack: packr2
	@$(call print, "Getting node dependencies.")
	(cd pos && npm install)
	@$(call print, "Compiling point-of-sale assets.")
	(cd pos && npm run export)
	@$(call print, "Packaging static assets.")
	$(PACKR2_BIN)

compile: pack
	@$(call print, "Building sweetd.")
	$(GOBUILD) -o sweetd -ldflags $(LDFLAGS) $(PKG)

test:
	@$(call print, "Testing sweetd.")
	go test -v ./...

clean:
	@$(call print, "Cleaning static asset packages.")
	packr2 clean
	@$(call print, "Cleaning builds and module cache")
	$(RM) ./sweetd

clean-cache:
	@$(call print, "Cleaning go module cache")
	go clean --modcache

build: compile
