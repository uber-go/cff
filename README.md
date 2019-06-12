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

### Developing on CFF ###

`internal/compile.go` and `internal/gen.go` contain the code for the static analysis and Go code generation respectively. 

#### Tests ####

"Golden" tests are under the `internal/tests` folder, which is written as one folder per test, and a single CFF 
source file matching the directory name (e.g. `internal/tests/basic/basic.go`). These have CFF sources that we want 
to assert (1) the CFF compiler works on them correctly, and (2) that the behavior of the generated code is as we
expect. (1) is enforced by the bazel rule in each directory, and (2) is enforced by `*_test.go` files in each directory. 

Failing test cases are in `internal/failing_tests` and are processed by `aquaregia_test.go` which does **not** use the 
bazel rule for CFF, because we want to assert (1) that the source code fails the CFF compiler, and (2) assert on the
error that was returned for the compiler. 

##### Benchmarks #####

Benchmarks are a special case of golden tests. They can be invoked using bazel test as follows:

```shell
$ bazel run //src/go.uber.org/cff/internal/tests/benchmark:go_default_test -- --test.v --test.bench=. --test.benchmem
Executing tests from //src/go.uber.org/cff/internal/tests/benchmark:go_default_test
-----------------------------------------------------------------------------
goos: darwin
goarch: amd64
BenchmarkSimple-8         	  500000	      2162 ns/op	      80 B/op	       6 allocs/op
BenchmarkSimpleNative-8   	 1000000	      1567 ns/op	      32 B/op	       3 allocs/op
BenchmarkMetrics-8        	  200000	     10817 ns/op	    3128 B/op	      35 allocs/op
```
