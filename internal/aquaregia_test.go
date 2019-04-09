package internal

import (
	"fmt"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/multierr"
	"golang.org/x/tools/go/packages"
)

const (
	aquaregiaTestDir = "failing_tests"
)

type errorCase struct {
	File         string
	ErrorMatches string
}

// Unwrap a single multierr.Error value, possibly nested, into a list of underlying errors
func unwrapMultierror(err error) []error {
	errs := multierr.Errors(err)

	if len(errs) == 1 && errs[0] == err {
		// If the underlying type is not a multierr, the multierr package will return it as-is
		return errs
	}

	var unwrappedErrors []error
	for _, err := range errs {
		unwrappedErrors = append(unwrappedErrors, unwrapMultierror(err)...)
	}

	return unwrappedErrors
}

// This works by trying to transpile all of the files in the directory/package,
// accumulating all of the errors, and then checking if in a file there is an
// instance of the error we are looking for.
// Note: error accumulation is per-package so at the moment state is kept
// when running transpiler across many flows expected to fail.
func TestCodeGenerateFails(t *testing.T) {
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
			{
				File:         "cff-task-arguments.go",
				ErrorMatches: "expected function, got bool",
			},
			{
				File:         "cff-task-arguments.go",
				ErrorMatches: "expected a function call, got identifier",
			},
			{
				File:         "cff-task-arguments.go",
				ErrorMatches: "only cff functions may be passed as task options",
			},
			{
				File:         "cff-task-arguments.go",
				ErrorMatches: "only cff functions may be passed as task options: found package \"go.uber.org/cff/internal/failing_tests/bad-inputs\"",
			},
			{
				File:         "cff-task-arguments.go",
				ErrorMatches: "expected cff.Task, got cff.Params; only cff.Task is allowed to be nested under cff.Tasks",
			},
			{
				File:         "cff-task-arguments.go",
				ErrorMatches: "expected function, got bool",
			},
			{
				File:         "cff-task-arguments.go",
				ErrorMatches: "expected function, got untyped nil",
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
				File:         "instrument-flowscope.go",
				ErrorMatches: "cff.Instrument requires a tally.Scope and \\*zap.Logger to be provided: use cff.Metrics and cff.Logger",
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
				File:         "cycle.go",
				ErrorMatches: "cycle detected",
			},
		},
	}

	for testDirectoryName, errCases := range errorCasesByDirectory {
		t.Run(fmt.Sprintf("test cases for directory %s", testDirectoryName), func(t *testing.T) {
			tempDir, err := ioutil.TempDir("", "cff-test")
			require.NoError(t, err)
			defer func() {
				assert.NoError(t, os.RemoveAll(tempDir))
			}()

			fset := token.NewFileSet()
			pkgs, err := packages.Load(&packages.Config{
				Mode:       packages.LoadSyntax,
				Fset:       fset,
				BuildFlags: []string{"-tags=cff"},
			}, filepath.Join(goldenTestImportInternal, aquaregiaTestDir, testDirectoryName, "..."))

			require.NoError(t, err, "could not load packages")
			require.NotEmpty(t, pkgs, "didn't find any packages")

			for _, pkg := range pkgs {
				errUntyped := Process(fset, pkg, tempDir)

				errorsThisPackage := unwrapMultierror(errUntyped)
				for _, err := range errorsThisPackage {
					t.Logf("found error %q", err.Error())
				}

				for _, errCase := range errCases {
					found := false
					regexpError := regexp.MustCompile(fmt.Sprintf("%s.*%s", errCase.File, errCase.ErrorMatches))
					// TODO: verify exactly how many times we match the error in a file.
					for _, err := range errorsThisPackage {
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
