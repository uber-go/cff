# cff-bazel-aquery

WARNING: CAVEAT EMPTOR. Please read carefully how this tool works
before using it.

This tool provides a workaround for T3187853. It is run in Jenkins to
get the generated code from CFFv2 to be present on disk so that the
coverage tool will map the coverage information on to the generated
code, thereby increasing the code coverage of the packages that use
CFFv2.

## Usage

```shell
./write_cff_gen_files.sh
```

```shell
./write_cff_gen_files.sh --dry-run
```

## Implementation

We first run `bazel query 'kind(cff, ...)'` to find the usages of cff
bazel rule under the current directory. We run bazel aquery for each
rule to find the declared output files of the cff rules and pipe the
output of bazel aquery to this tool.  This tool will copy the
generated code from the bazel-out folder and put it alongside the
original source code.

Then, the Jenkins job should run the normal coverage tooling similar
to the normal monorepo build, but now with the generated files on disk
the coverage will map on to the generated file and count for the
package.

### Protobuf IDL and aquery fixture

The protobuf IDLs are taken from the bazel repository as of v0.29.0 1e3109275a732edc51f8f1a4cac8a066b713ffff

The fixtures in `testdata` folder are the output from bazel aquery
from the testdata/example folder and can be recreated with the
following command:

```shell
bazel aquery //src/go.uber.org/cff/cmd/cff-bazel-aquery/testdata/example:cff --output=proto | gzip > aquery.dat.gz
```

You can inspect the output of the fixture by running:

```shell
cat src/go.uber.org/cff/cmd/cff-bazel-aquery/testdata/aquery.dat.gz | gunzip | protoc --decode analysis.ActionGraphContainer src/go.uber.org/cff/internal/proto/analysis.proto
```

Note the above command must be run from the monorepo root because the
include paths of the .proto file are relative to monorepo root.
