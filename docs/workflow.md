# Workflow

At a high-level, the workflow for using cff is as follows:

1. Add `//go:build cff` to the top of the file if it's not already present.
2. Use functions defined in the go.uber.org/cff package.
3. Run `go generate` if you added a `//go:generate` directive,
   or manually run `cff ./...` if you didn't.
4. Run `go build` or `go test` as usual.
5. Commit the generated code to your repository.
