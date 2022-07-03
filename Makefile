# Build config
BUILD_FILES = $(shell go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}' ./...)

DATE=$(shell date -u "+%a %b %d %T %Y")
TUG_COMMIT ?= $(shell git rev-parse --short HEAD)
TUG_VERSION ?= dev

LDFLAGS = -s -w
LDFLAGS += -X "github.com/b4nst/turbogit/cmd.BuildDate=$(DATE)"
LDFLAGS += -X "github.com/b4nst/turbogit/cmd.Commit=$(TUG_COMMIT)"
LDFLAGS += -X "github.com/b4nst/turbogit/cmd.Version=$(TUG_VERSION)"

# Go config
BUILD_ARGS=-trimpath -tags=static -ldflags='$(LDFLAGS)'
GOCMD=go
GOBUILD=$(GOCMD) build $(BUILD_ARGS)
GOTEST=$(GOCMD) test -tags=static
GORUN=$(GOCMD) run

dist/bin/tug: $(BUILD_FILES)
	$(GOBUILD) -o "$@" ./main.go

build: libgit2 dist/bin/tug
.PHONY: build

libgit2:
	$(MAKE) -C ./git2go install-static

test: $(BUILD_FILES)
	$(GOTEST) ./...  -coverprofile c.out
.PHONY: test

doc:
	$(GORUN) scripts/gen-doc.go
.PHONY: doc

clean:
	rm -rf bin dist c.out
.PHONY: clean
