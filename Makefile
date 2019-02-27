.PHONY: test
test:
	buckw test ...

.PHONY: generate
generate:
	go test -run TestCodeIsUpToDate ./internal --generate
