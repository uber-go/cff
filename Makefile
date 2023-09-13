SHELL = /bin/bash

# Cross-platform way to find the directory holding this Makefile.
PROJECT_ROOT = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

export GOBIN = $(PROJECT_ROOT)/bin
export PATH := $(GOBIN):$(PATH)

MODULES ?= . ./examples ./internal/tests ./tools ./docs
TEST_FLAGS ?= -race

CFF = bin/cff
MOCKGEN = bin/mockgen
STATICCHECK = bin/staticcheck
MDOX = bin/mdox

# 'make cover' should not run on docs by default.
# We run that separately explicitly on a specific platform.
COVER_MODULES ?= $(filter-out ./docs ./tools,$(MODULES))

SRC_FILES = $(shell find . '(' -path './.*' -o -path '*test*' -o -path '*/examples/*' -o -path './docs/*' -prune ')' -o -name '*.go' -print)

# All go files are in scope for formatting -- even if they're generated.
GOFMT_FILES = $(shell find . -path './.*' -prune -o -name '*.go' -print)

.PHONY: build
build: $(CFF)

.PHONY: test
test: build
	@$(foreach dir,$(COVER_MODULES),( \
		cd $(dir) && \
		echo "--- [test] $(dir)" && \
		go test $(TEST_FLAGS) ./... \
	) &&) true

.PHONY: cover
cover: build
	@$(foreach dir,$(COVER_MODULES),( \
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

.PHONY: fmt
fmt:
	@gofmt -w -l $(GOFMT_FILES)

.PHONY: lint
lint: staticcheck checkfmt docs-check

.PHONY: staticcheck
staticcheck: $(STATICCHECK)
	$(STATICCHECK) ./...

.PHONY: checkfmt
checkfmt:
	@DIFF=$$(gofmt -d $(GOFMT_FILES)); \
	if [[ -n "$$DIFF" ]]; then \
		echo "--- gofmt would cause changes:"; \
		echo "$$DIFF"; \
		exit 1; \
	fi

.PHONY: docs
docs:
	cd docs && yarn build

.PHONY: docs-check
docs-check: $(MDOX)
	@echo "Checking documentation"
	@make -C docs check

$(CFF): $(SRC_FILES)
	go install go.uber.org/cff/cmd/cff

$(MOCKGEN): go.mod
	go install github.com/golang/mock/mockgen

$(STATICCHECK): tools/go.mod
	cd tools && go install honnef.co/go/tools/cmd/staticcheck

$(MDOX): tools/go.mod
	cd tools && go install github.com/bwplotka/mdox
