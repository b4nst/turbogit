# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

BUILD_FILES = $(shell go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}' ./...)

bin/tug: $(BUILD_FILES)
	$(GOBUILD) -trimpath -o "$@" ./main.go
build: bin/tug
.PHONY: build

test: $(BUILD_FILES)
	$(GOTEST) ./...
.PHONY: test