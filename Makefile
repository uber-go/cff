export GOBIN = $(shell pwd)/bin
export PATH := $(GOBIN):$(PATH)

CFF = $(GOBIN)/cff
MOCKGEN = $(GOBIN)/mockgen

SRC_FILES = $(shell find . '(' -path '*test*' -o -path '*/examples/*' -prune ')' -o -name '*.go' -print)

.PHONY: build
build: $(CFF)

.PHONY: test
test: build
	go test -race ./...

# We run 'go generate' with '-tags cff'
# because build tags can be inside cff-tagged files.
.PHONY: generate
generate: $(CFF) $(MOCKGEN)
	go generate -tags cff -x ./...

$(CFF): $(SRC_FILES)
	go install go.uber.org/cff/cmd/cff

$(MOCKGEN): go.mod
	go install github.com/golang/mock/mockgen
