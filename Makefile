CFF = bazel run --run_under="cd $(shell pwd) &&" //src/go.uber.org/cff/cmd/cff --

.PHONY: test
test:
	bazel test ...

.PHONY: generate
generate:
	 $(CFF) --input examples/magic.go --output examples/magic_gen.go
