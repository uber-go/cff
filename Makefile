CFF = bazel run --run_under="cd $(shell pwd) &&" //src/go.uber.org/cff/cmd/cff --

.PHONY: test
test:
	bazel test ...

.PHONY: generate
generate:
	 $(CFF) go.uber.org/cff/internal/tests/basic --file=basic.go=/tmp/basic_gen.go
