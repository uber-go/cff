.PHONY: test
test:
	buckw test ...

.PHONY: generate
generate:
	go run ./cmd/cff ./internal/tests/...
