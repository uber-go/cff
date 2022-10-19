export GOBIN = $(shell pwd)/bin
export PATH := $(GOBIN):$(PATH)

MODULES ?= . ./examples ./internal/tests
TEST_FLAGS ?= -race

CFF = $(GOBIN)/cff
MOCKGEN = $(GOBIN)/mockgen

SRC_FILES = $(shell find . '(' -path '*test*' -o -path '*/examples/*' -prune ')' -o -name '*.go' -print)

.PHONY: build
build: $(CFF)

.PHONY: test
test: build
	@$(foreach dir,$(MODULES),( \
		cd $(dir) && \
		echo "--- [test] $(dir)" && \
		go test $(TEST_FLAGS) ./... \
	) &&) true

# We run 'go generate' with '-tags cff'
# because build tags can be inside cff-tagged files.
.PHONY: generate
generate: $(CFF) $(MOCKGEN)
	@$(foreach dir,$(MODULES),( \
		cd $(dir) && \
		echo "--- [generate] $(dir)" && \
		go generate -tags cff -x ./... \
	) &&) true

.PHONY: tidy
tidy:
	@$(foreach dir,$(MODULES),( \
		cd $(dir) && \
		echo "--- [tidy] $(dir)" && \
		go mod tidy \
	) &&) true

$(CFF): $(SRC_FILES)
	go install go.uber.org/cff/cmd/cff

$(MOCKGEN): go.mod
	go install github.com/golang/mock/mockgen
