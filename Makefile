export GOBIN = $(shell pwd)/bin
export PATH := $(GOBIN):$(PATH)

MODULES ?= . ./examples ./internal/tests
TEST_FLAGS ?= -race

CFF = $(GOBIN)/cff
MOCKGEN = $(GOBIN)/mockgen
STATICCHECK = $(GOBIN)/staticcheck

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

.PHONY: cover
cover: build
	@$(foreach dir,$(MODULES),( \
		cd $(dir) && \
		echo "--- [cover] $(dir)" && \
		go test $(TEST_FLAGS) -coverprofile=cover.out -coverpkg=go.uber.org/cff/... ./... && \
		go tool cover -html=cover.out -o cover.html \
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

.PHONY: lint
lint: staticcheck

.PHONY: staticcheck
staticcheck: $(STATICCHECK)
	$(STATICCHECK) ./...

$(CFF): $(SRC_FILES)
	go install go.uber.org/cff/cmd/cff

$(MOCKGEN): go.mod
	go install github.com/golang/mock/mockgen

$(STATICCHECK): tools/go.mod
	cd tools && go install honnef.co/go/tools/cmd/staticcheck
