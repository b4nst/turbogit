# Build config
PLUGIN_DIRS = $(sort $(dir $(wildcard cmd/*/*)))
PLUGIN_BINS = $(PLUGIN_DIRS:cmd/%/=%)
BUILD_FILES = $(shell go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}' ./...)

# Go config
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GORUN=$(GOCMD) run
LDFLAGS = -s -w

bin/$(PLUGIN_BINS): $(BUILD_FILES)
	$(GOBUILD) -trimpath -o "$@" -ldflags='$(LDFLAGS)' cmd/$(@F)/main.go

build: bin/$(PLUGIN_BINS)
.PHONY: build

test: $(BUILD_FILES)
	$(GOTEST) ./...  -coverprofile c.out
.PHONY: test

doc:
	$(GORUN) scripts/gen-doc.go
.PHONY: doc

clean:
	rm -rf bin dist c.out
.PHONY: clean
