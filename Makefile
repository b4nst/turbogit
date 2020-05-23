# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

DATE=$(shell date -u "+%a %b %d %T %Y")
COMMIT=$(shell git rev-parse --short HEAD)

BUILD_FILES = $(shell go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}' ./...)

bin/tug: $(BUILD_FILES)
	$(GOBUILD) -trimpath -o "$@" -ldflags='-s -w -X "github.com/b4nst/turbogit/cmd.BuildDate=$(DATE)" -X github.com/b4nst/turbogit/cmd.Commit=$(COMMIT)' ./main.go 
build: bin/tug
.PHONY: build

test: $(BUILD_FILES)
	$(GOTEST) ./...  -coverprofile c.out
.PHONY: test