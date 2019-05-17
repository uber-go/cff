.PHONY: test
test:
	buckw test ...

.PHONY: generate
generate:
	go run ./cmd/cff --input ./examples/magic.go --output /tmp/magic_gen.go
