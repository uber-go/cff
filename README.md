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
    cff.Task(
        client.GetUser),
    cff.Task(
        func(*GetUserResponse) *FooResults {
            ...
        }),
    cff.Task( mapper.FormatResponse),
)
```

## Bazel Rules

If you write the code above and then run gazelle, you'll see your `BUILD.bazel`
file change. Gazelle will add two bazel rules: `cff` and `go_library`. The `cff`
rule defines how to generate CFF code and the `go_library` rule defines a library
that uses the generated code.

An example:
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

Let's break down each argument.

This allows the usage of the CFF bazel rule.
```
load("//rules:cff.bzl", "cff")
```


The name from the CFF rule is used in a `go_library` to access the generated
code.
```
name = "cff"
```

Internal functions that are dependencies of CFF flows. They are declared to
be added onto $GOPATH when compiling a CFF flow. In this case, `aflow.go` uses
`afunc.go` within the flow:
````
srcs = [
        "//src/go.uber.org/cff/internal/tests/sandwich:afunc.go",
    ],
````

CFF sources for the cff rule. CFF source files are defined as any Go files that
contain a CFF code generation directive such as `cff.Flow` or `cff.Task`.
````
cff_srcs = ["aflow.go", "bflow.go"],
````

Import path leading to package containing CFF source code:
````
importpath = "go.uber.org/cff/internal/tests/sandwich",
````

Dependencies of the CFF flows. This will contain at least the interface to supported CFF
options.
````
deps = [
        "//src/go.uber.org/cff:go_default_library",
    ],
````

CFF Bazel rule outputs all generated files with current convention of `_gen.go`
appended to the source file. You can view the generated code by adding a
`--output_groups=go_generated_srcs ` argument when building your target:
```
bazel build --output_groups=go_generated_srcs //src/go.uber.org/cff/internal/tests/nested_parent:cff
```

A library is created out of the generated files via
```
go_library(
    name = "flowcaller",
    srcs = [
        "afunc.go",
        "flowcaller.go",
        ":cff",
    ],
    importpath = "go.uber.org/cff/internal/tests/sandwich",
    visibility = ["//visibility:public"],
)
```
where `:cff` is the target created above.

In order to auto generate all of these rules, just run `bin/gazelle` on your
chosen package.

### CLI ###
To experiment with CFF you can run the `cff` command on a specific package and specify the output path.
The script lives under monorepo root in `bin/cff` eg `$GOPATH/bin/cff` and if `$PATH` contains it, can be called via
`cff` for example,

```shell
$ cff go.uber.org/cff/internal/tests/basic --file=basic.go=/tmp/basic_gen.go
```

This will generate `basic.go` inside `go.uber.org/cff/internal/tests/basic` to
`/tmp/basic_gen.go`, and the relevant sections eg `cff.Flow` in the code will be replaced
with generated code.

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

##### Test speedup #####

Currently, benchmark files after transpilation add up to over 120,000 LOC, and even though we dont run benchmarks via
`bazel test //src/go.uber.org/cff/...`, those files are still compiled. To speed up your unit tests,
you can run with a tag that excludes the slow target: `bazel test --test_tag_filters="-cff_slow_compile" --build_tag_filters="-cff_slow_compile" //src/go.uber.org/cff/...`

##### Benchmarks #####

Benchmarks are a special case of golden tests. They can be invoked using bazel test as follows:

```shell
$ bazel run //src/go.uber.org/cff/internal/tests/benchmark:go_default_test -- --test.v --test.bench=. --test.benchmem
Executing tests from //src/go.uber.org/cff/internal/tests/benchmark:go_default_test
-----------------------------------------------------------------------------
goos: darwin
goarch: amd64
BenchmarkSimple-8              1000000          1830 ns/op          80 B/op           6 allocs/op
BenchmarkSimpleNative-8        1000000          1320 ns/op          32 B/op           3 allocs/op
BenchmarkMetrics-8              100000         20779 ns/op        6136 B/op          70 allocs/op
BenchmarkMetrics100-8            10000        261341 ns/op      162279 B/op        1637 allocs/op
BenchmarkMetrics500-8             2000       1351223 ns/op      799228 B/op        8039 allocs/op
BenchmarkMetrics1000-8            1000       2376714 ns/op     1594208 B/op       16029 allocs/op
```
