CFF = bazel run //src/go.uber.org/cff/cmd/cff --

.PHONY: test
test:
	bazel test ...

.PHONY: generate
generate:
	 $(CFF) --input ./examples/magic.go --output /tmp/magic_gen.go
