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

## Bazel Rule

Create cff Bazel target for your package. 

A full example: 
```
load("//rules:cff.bzl", "cff")

cff(
    name = "cff",
    srcs = [
        "//src/go.uber.org/cff/internal/tests/sandwich:afunc.go",
    ],
    cff_srcs = ["aflow.go", "bflow.go"],
    importpath = "go.uber.org/cff/internal/tests/sandwich",
    visibility = ["//visibility:public"],
    deps = [
        "//src/go.uber.org/cff:go_default_library",
    ],
)
```


so let's break down each argument.

In your `BUILD.bazel` file, add import to cff rule:
```
load("//rules:cff.bzl", "cff")
```

This is your only CFF target for the package, so name something relevant to CFF. 
This name will be used by `go_library` to export your package containing CFF flow:  
```
name = "cff"
```

Internal functions that are dependencies of your CFF flows. They are declared to
be added onto $GOPATH when compiling your CFF flow. In this case, `aflow.go` uses
`afunc.go` within the flow:
````
srcs = [
        "//src/go.uber.org/cff/internal/tests/sandwich:afunc.go",
    ],
````

CFF sources containing `+build cff` tag. You can refer to sources by the full label including the package 
or just the file name. In this case, `aflow.go` and `bflow.go` are within the same
package.
````
cff_srcs = ["aflow.go", "bflow.go"],
````

Import path leading to package containing CFF source code:
````
importpath = "go.uber.org/cff/internal/tests/sandwich",
````

Dependencies of the CFF flows. Must contain at least the interface to supported CFF 
options.
````
deps = [
        "//src/go.uber.org/cff:go_default_library",
    ],
````

CFF Bazel rule outputs all generated files with current convention of `_gen.go` 
appended to the source file. You can view the generated code by adding a 
`--output_groups=debug_files ` argument when building your target: 
```
bazel build --output_groups=debug_files //src/go.uber.org/cff/internal/tests/nested_parent:cff
```

### Using CFF Rule ###

After creating the cff target, you can build a library out of the generated files via
```
go_library(
    name = "flowcaller",
    srcs = [
        "afunc.go",
        "flowcaller.go",
        ":cff",  # keep
    ],
    importpath = "go.uber.org/cff/internal/tests/sandwich",
    visibility = ["//visibility:public"],
)
```
where `:cff` is the target we created above. 

Note that until gazelle rules are in place, `#keep` signals to gazelle not to 
delete this source.  

### CLI ###
Alternatively, to experiment with CFF you can run the `cff` command on a specific package and specify the output.

```shell
$ bazel run //src/go.uber.org/cff/cmd/cff -- go.uber.org/cff/internal/tests/basic --file=basic.go=/tmp/basic_gen.go
```

This will generate copies of the original files without the `cff` tag, and the
relevant sections in the code replaced with generated code.
