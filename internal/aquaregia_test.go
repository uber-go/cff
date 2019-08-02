package internal

import (
	"fmt"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/packages/packagestest"
)

const (
	aquaregiaTestDir = "failing_tests"
	internalTests    = "go.uber.org/cff/internal"
)

type errorCase struct {
	File         string
	ErrorMatches string
}

// This works by trying to transpile all of the files in the directory/package,
// accumulating all of the errors, and then checking if in a file there is an
// instance of the error we are looking for.
// Note: error accumulation is per-package so at the moment state is kept
// when running transpiler across many flows expected to fail.
func TestCodeGenerateFails(t *testing.T) {
	t.Skip("Flaky test: T3593709")
	// map [directory name] -> list of test cases
	errorCasesByDirectory := map[string][]errorCase{
		"bad-inputs": {
			{
				File:         "already-provided.go",
				ErrorMatches: "type string already provided at",
			},
			{
				File:         "cff-flow-arguments.go",
				ErrorMatches: "cff.Flow expects at least one function",
			},
			{
				File:         "cff-flow-arguments.go",
				ErrorMatches: "expected cff call but got <nil>",
			},
			{
				File:         "cff-flow-arguments.go",
				ErrorMatches: "expected cff call but got field ProvidesBad func\\(\\) go.uber.org/cff.FlowOption",
			},
			{
				File:         "cff-flow-arguments.go",
				ErrorMatches: "expected a function call, got identifier",
			},
			// ExpectsFunctionCallExpression
			{
				File:         "cff-task-arguments.go",
				ErrorMatches: "expected function, got bool",
			},
			// ExpectedFlowArgumentsSelectorExpression.
			{
				File:         "cff-task-arguments.go",
				ErrorMatches: "only cff functions may be passed as task options",
			},
			// ExpectedFlowArgumentsCallExpressions
			{
				File:         "cff-task-arguments.go",
				ErrorMatches: "expected a function call, got identifier",
			},
			// ExpectedFlowArgumentsCallExpressions
			{
				File:         "cff-task-arguments.go",
				ErrorMatches: "unknown top-level cff function \"Instrument\": only cff.Flow may be called at the top-level",
			},
			// ExpectedFlowArgumentsNotCFF
			{
				File:         "cff-task-arguments.go",
				ErrorMatches: "only cff functions may be passed as task options: found package \"go.uber.org/cff/internal/failing_tests/bad-inputs\"",
			},
			// ExpectedTasksBad
			{
				File:         "cff-task-arguments.go",
				ErrorMatches: "expected function, got untyped nil",
			},
			// ExpectedTasksBadCallExpr
			{
				File:         "cff-task-arguments.go",
				ErrorMatches: "expected cff.Task, got cff.Params; only cff.Task is allowed to be nested under cff.Tasks",
			},
			// ExpectedTasksBadCallExprNotCFF
			{
				File:         "cff-task-arguments.go",
				ErrorMatches: "expected function, got int",
			},

			{
				File:         "context-predicate.go",
				ErrorMatches: "only the first argument may be context.Context",
			},
			{
				File:         "context-task.go",
				ErrorMatches: "only the first argument may be context.Context",
			},
			{
				File:         "earlyresult.go",
				ErrorMatches: "unused output type string",
			},
			{
				File:         "earlyresult.go",
				ErrorMatches: "unused output type int32",
			},
			{
				File:         "earlyresult.go",
				ErrorMatches: "unused output type \\*go.uber.org/cff/internal/failing_tests/bad-inputs.quuz",
			},
			{
				File:         "earlyresult.go",
				ErrorMatches: "unused output type \\*go.uber.org/cff/internal/failing_tests/bad-inputs.corge",
			},
			{
				File:         "earlyresult.go",
				ErrorMatches: "unused output type \\*go.uber.org/cff/internal/failing_tests/bad-inputs.grault",
			},
			{
				File:         "error-task.go",
				ErrorMatches: "only the last result may be an error",
			},
			{
				File:         "fallback-with.go",
				ErrorMatches: "cff.FallbackWith result at position 1 of type string cannot be used as bool",
			},
			{
				File:         "fallback-with.go",
				ErrorMatches: "cff.FallbackWith result at position 2 of type bool cannot be used as string",
			},
			{
				File:         "fallback-with.go",
				ErrorMatches: "cff.FallbackWith must produce the same number of results as the task: expected 2, got 1",
			},
			{
				File:         "fallback-with.go",
				ErrorMatches: "Task must return an error for FallbackWith to be used",
			},
			{
				File:         "missing-provider.go",
				ErrorMatches: "no provider found for float64",
			},
			{
				File:         "nonpointer-result.go",
				ErrorMatches: "invalid parameter to cff.Results: expected pointer, got bool",
			},
			{
				File:         "predicate.go",
				ErrorMatches: "the function must return a single boolean result",
			},
			{
				File:         "predicate-params.go",
				ErrorMatches: "cff.Predicate expected a function but received",
			},
			{
				File:         "top-level-flow.go",
				ErrorMatches: "unknown top-level cff function \"Predicate\"",
			},
			{
				File:         "unused-inputs.go",
				ErrorMatches: "unused input type string",
			},
			{
				File:         "variadic.go",
				ErrorMatches: "variadic functions are not yet supported",
			},
		},
		"cycles": {
			{
				File:         "easy-cycle.go",
				ErrorMatches: "cycle detected: need to run \\[func\\(int64\\) string\\] to provide string",
			},
			{
				File:         "no-output.go",
				ErrorMatches: "cycle detected: need to run \\[func\\(int32\\) string\\] to provide string",
			},
		},
	}

	_ = os.Setenv("PATH", os.ExpandEnv("$TEST_SRCDIR/__main__/external/go_sdk/bin:$PATH"))
	// See if we are in Bazel environment as `go test` doesn't need GOCACHE to be set manually.
	if file, err := os.Stat(os.Getenv("TEST_TMPDIR")); err == nil {
		if file.IsDir() {
			// Go executable requires a GOCACHE to be set after go1.12.
			_ = os.Setenv("GOCACHE", filepath.Join(os.Getenv("TEST_TMPDIR"), "/cache"))
		}
	}

	// Including entire project including tests as a module. This is different from golden_test as
	// we dont need to include generated files as all of these tests are expected to fail.
	cffModule := packagestest.Module{
		Name:  "go.uber.org/cff",
		Files: packagestest.MustCopyFileTree("./.."),
	}

	for testDirectoryName, errCases := range errorCasesByDirectory {
		t.Run(fmt.Sprintf("test cases for directory %s", testDirectoryName), func(t *testing.T) {
			exp := packagestest.Export(t, packagestest.Modules, []packagestest.Module{cffModule})
			fset := token.NewFileSet()

			cfg := exp.Config
			cfg.Mode = packages.LoadSyntax
			cfg.BuildFlags = []string{"-tags=cff"}
			cfg.Fset = fset

			defer exp.Cleanup()

			// Using pattern for go test not to run _test unit tests which test generated code.

			pkgs, err := packages.Load(
				exp.Config,
				"pattern="+filepath.Join(internalTests, aquaregiaTestDir, testDirectoryName, "..."))

			require.NoError(t, err, "could not load packages")
			require.NotEmpty(t, pkgs, "didn't find any packages")

			for _, pkg := range pkgs {
				// Output path can be empty so code gets generated next to source in case of failed
				// tests.
				var errors []error
				for i := range pkg.CompiledGoFiles {
					if err := Process(fset, pkg, pkg.Syntax[i], ""); err != nil {
						errors = append(errors, err)
					}
				}
				for _, err := range errors {
					t.Logf("found error %q", err.Error())
				}
				for _, errCase := range errCases {
					found := false
					regexpError := regexp.MustCompile(fmt.Sprintf("%s.*%s", errCase.File, errCase.ErrorMatches))
					// TODO: verify exactly how many times we match the error in a file.
					for _, err := range errors {
						if ok := regexpError.MatchString(err.Error()); ok {
							found = true
							break
						}
					}

					assert.True(t, found, "expected to find error matching %q in %q", errCase.ErrorMatches, errCase.File)
				}
			}
		})
	}
}
