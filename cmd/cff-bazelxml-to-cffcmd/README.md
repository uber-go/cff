# cff-bazelxml-to-cffcmd

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

We first run "bazel query 'kind(cff, ...)'" to find the usages of cff
bazel rule under the current directory. We output the bazel query to
XML, and pipe the XML to the tool in this folder. The tool will look
at every rule and generate a shell command that looks like:

  bazel run //src/go.uber.org/cff/cmd/cff:cff -- --file=<file> ... <importpath>

It will then run those commands unless the dry-run flag is
set. Because we directly pass the arguments to the argv bazel command,
there is no concern about command injection. When the cff tool runs,
it will output the generated files alongside the source files. Then,
the Jenkins job should run the normal coverage tooling similar to the
normal monorepo build, but now with the generated files on disk the
coverage will map on to the generated file and count for the package.