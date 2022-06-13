# Build config
PLUGIN_DIRS = $(sort $(dir $(wildcard cmd/*/*)))
PLUGIN_BINS = $(addprefix dist/bin/, $(PLUGIN_DIRS:cmd/%/=%))
BUILD_FILES = $(shell go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}' ./...)

# Go config
LDFLAGS=-s -w
BUILD_ARGS=-trimpath -tags=static -ldflags='$(LDFLAGS)'
GOCMD=go
GOBUILD=$(GOCMD) build $(BUILD_ARGS)
GOTEST=$(GOCMD) test -tags=static
GORUN=$(GOCMD) run

$(PLUGIN_BINS): $(BUILD_FILES)
	$(GOBUILD) -o "$@" cmd/$(@F)/main.go

build: libgit2 $(PLUGIN_BINS)
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
