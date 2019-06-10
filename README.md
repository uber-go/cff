**DO NOT USE CFF UNTIL THIS NOTICE GOES AWAY.**

This is not ready for use by anyone. It does not work in the monorepo at this
time. Please reach out to the owners (see METADATA) to ask about using this
project.

---

CFF is intended to be a system which facilitates orchestrating large numbers
of interrelated functions with as much concurrent execution as possible.

# Concepts

In CFF, a **Task** is a single executable function or bound method. Tasks have
**inputs** and **outputs**: the parameters and return values of the
corresponding functions. One or more tasks come together to form a **Flow**.
Flows have zero or more **inputs** and one or more **outputs**.

A single Flow must be self-contained. That is, all inputs of all Tasks in a
Flow must come from either another task or as one of the inputs of the Flow
itself.

# Usage

To use CFF, write code using the APIs exported by the
`go.uber.org/cff` package.

```go
var result Response
err := cff.Flow(
    cff.Params(request),
    cff.Results(&response),
    cff.Tasks(
        client.GetUser,
        func(*GetUserResponse) *FooResults {
            ...
        },
        mapper.FormatResponse,
    ),
)
```

Tag the files you used CFF in with the `cff` build tag. This excludes them from
being built by the Go compiler.

```
// +build cff

package userservice
```

Finally, run the `cff` command on a specific package and specify the output.

```shell
$ bazel run //src/go.uber.org/cff/cmd/cff -- go.uber.org/cff/internal/tests/basic --file=basic.go=/tmp/basic_gen.go
```

This will generate copies of the original files without the `cff` tag, and the
relevant sections in the code replaced with generated code.
