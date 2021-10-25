package internal

import (
	"embed"
	"fmt"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"code.uber.internal/devexp/bazel/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/packages/packagestest"
)

const (
	aquaregiaTestDir = "failing_tests"
	internalTests    = "go.uber.org/cff/internal"
)

//go:embed package_fixtures/*
var pkgFixtures embed.FS

type errorCase struct {
	File         string
	ErrorMatches string
}

var codeGenerateFailCases = map[string][]errorCase{
	// map [directory name] -> list of test cases
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
			ErrorMatches: "expected cff call but got field ProvidesBad func\\(\\) go.uber.org/cff.Option",
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
			ErrorMatches: `unexpected code generation directive "Instrument": only cff.Flow or cff.Parallel may be called at the top-level`,
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
			ErrorMatches: "expected function, got go.uber.org/cff.Option",
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
		// EarlyResultMultipleFlows
		{
			File:         "earlyresult.go",
			ErrorMatches: "unused output type \\*go.uber.org/cff/internal/failing_tests/bad-inputs.quuz",
		},
		// EarlyResultMultipleFlows
		{
			File:         "earlyresult.go",
			ErrorMatches: "unused output type \\*go.uber.org/cff/internal/failing_tests/bad-inputs.corge",
		},
		// EarlyResultMultipleFlows
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
		// MissingCFFMetrics
		{
			File:         "instrument.go",
			ErrorMatches: "cff.Instrument requires a cff\\.Emitter to be provided: use cff\\.WithEmitter",
		},
		// MissingCFFLogger and MissingCFFLoggerME
		{
			File:         "instrument.go",
			ErrorMatches: "cff.Instrument requires a cff\\.Emitter to be provided: use cff\\.WithEmitter",
		},
		// MissingCFFLoggerAndMetrics
		{
			File:         "instrument.go",
			ErrorMatches: "cff.Instrument requires a cff\\.Emitter to be provided: use cff\\.WithEmitter",
		},

		{
			File:         "missing-provider.go",
			ErrorMatches: "no provider found for float64",
		},
		{
			File:         "nonpointer-result.go",
			ErrorMatches: "invalid parameter to cff.Results: expected pointer, got bool",
		},
		// ParallelInvalidParamsType
		{
			File:         "parallel.go",
			ErrorMatches: "the only allowed argument is a single context.Context parameter",
		},
		// ParallelInvalidParamsMultiple
		{
			File:         "parallel.go",
			ErrorMatches: "the only allowed argument is a single context.Context parameter",
		},
		// ParallelInvalidReturnType
		{
			File:         "parallel.go",
			ErrorMatches: "the only allowed return value is an error",
		},
		// ParallelInvalidFuncVar
		{
			File:         "parallel.go",
			ErrorMatches: "the only allowed argument is a single context.Context parameter",
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
			File:         "unused-outputs.go",
			ErrorMatches: "unused output type bool",
		},
		{
			File:         "unused-outputs.go",
			ErrorMatches: "unused output type uint32",
		},
		{
			File:         "top-level-flow.go",
			ErrorMatches: "unexpected code generation directive \"Predicate\"",
		},
		{
			File:         "unused-inputs.go",
			ErrorMatches: "unused input type string",
		},
		// UnsupportedInvoke
		{
			File:         "unused-task.go",
			ErrorMatches: "cff\\.Invoke cannot be provided on a Task that produces values besides errors",
		},
		// NoInvokeNoResults
		{
			File:         "unused-task.go",
			ErrorMatches: "task must return at least one non-error value but currently produces zero.",
		},
		// NoInvokeWithError
		{
			File:         "unused-task-error.go",
			ErrorMatches: "task must return at least one non-error value but currently produces zero.",
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

// This works by trying to transpile all of the files in the directory/package,
// accumulating all of the errors, and then checking if in a file there is an
// instance of the error we are looking for.
// Note: error accumulation is per-package so at the moment state is kept
// when running transpiler across many flows expected to fail.
func TestCodeGenerateFails(t *testing.T) {
	for testDirectoryName, errCases := range codeGenerateFailCases {
		t.Run(fmt.Sprintf("test cases for directory %s", testDirectoryName), func(t *testing.T) {
			fset := token.NewFileSet()
			pkgs, err := loadAquaregiaPackages(
				&loadParams{
					pattern: "pattern=" + filepath.Join(internalTests, aquaregiaTestDir, testDirectoryName, "..."),
					fset:    fset,
					t:       t,
				},
			)
			require.NoError(t, err, "could not load packages")
			require.NotEmpty(t, pkgs, "didn't find any packages")

			processor := Processor{Fset: fset}

			for _, gopkg := range pkgs {
				pkg := newPackage(gopkg)
				// Output path can be empty so code gets generated next to source in case of failed
				// tests.
				var errors []error
				for i := range pkg.CompiledGoFiles {
					if err := processor.Process(pkg, pkg.Syntax[i], ""); err != nil {
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

type loadParams struct {
	pattern string
	fset    *token.FileSet
	t       *testing.T
}

func loadAquaregiaPackages(p *loadParams) ([]*packages.Package, error) {
	_ = os.Setenv("PATH", os.ExpandEnv("$TEST_SRCDIR/__main__/external/go_sdk/bin:$PATH"))
	// See if we are in Bazel environment as `go test` doesn't need GOCACHE to be set manually.
	if file, err := os.Stat(os.Getenv("TEST_TMPDIR")); err == nil {
		if file.IsDir() {
			// Go executable requires a GOCACHE to be set after go1.12.
			_ = os.Setenv("GOCACHE", filepath.Join(os.Getenv("TEST_TMPDIR"), "/cache"))
		}
	}

	// packages.Load expects to find the import when parsing files.
	cffModule := packagestest.Module{
		Name:  "go.uber.org/cff",
		Files: packagestest.MustCopyFileTree("./.."),
	}
	tallyContent, err := pkgFixtures.ReadFile("package_fixtures/tally.fixture")
	require.NoError(p.t, err)

	tallyModule := packagestest.Module{
		Name: "github.com/uber-go/tally",
		Files: map[string]interface{}{
			// This needs to be a valid Go file.
			"tally.go": string(tallyContent),
		},
		Overlay: nil,
	}

	zapContent, err := pkgFixtures.ReadFile("package_fixtures/zap.fixture")
	require.NoError(p.t, err)

	zapModule := packagestest.Module{
		Name: "go.uber.org/zap",
		Files: map[string]interface{}{
			// This needs to be a valid Go file.
			"zap.go": string(zapContent),
		},
		Overlay: nil,
	}

	observerContent, err := pkgFixtures.ReadFile("package_fixtures/observer.fixture")
	require.NoError(p.t, err)

	zapTestObserverModule := packagestest.Module{
		Name: "go.uber.org/zap/zaptest/observer",
		Files: map[string]interface{}{
			// This needs to be a valid Go file.
			"observer.go": string(observerContent),
		},
		Overlay: nil,
	}

	exp := packagestest.Export(
		p.t,
		packagestest.Modules,
		[]packagestest.Module{
			cffModule,
			tallyModule,
			zapModule,
			zapTestObserverModule,
		},
	)

	cfg := exp.Config
	cfg.BuildFlags = []string{"-tags=cff"}
	cfg.Fset = p.fset
	cfg.Tests = false
	cfg.Mode = packages.NeedName |
		packages.NeedFiles |
		packages.NeedCompiledGoFiles |
		packages.NeedImports |
		packages.NeedDeps |
		packages.NeedTypes |
		packages.NeedSyntax |
		packages.NeedTypesInfo |
		packages.NeedTypesSizes
	defer exp.Cleanup()
	pkgs, err := packages.Load(exp.Config, p.pattern)
	return pkgs, err
}

func TestCasesGoCompile(t *testing.T) {
	for testDirectoryName := range codeGenerateFailCases {
		t.Run(fmt.Sprintf("test cases for directory %s", testDirectoryName), func(t *testing.T) {
			fset := token.NewFileSet()
			pkgs, err := loadAquaregiaPackages(
				&loadParams{
					pattern: "pattern=" + filepath.Join(internalTests, aquaregiaTestDir, testDirectoryName, "..."),
					fset:    fset,
					t:       t,
				},
			)
			require.NoError(t, err, "could not load packages")
			require.NotEmpty(t, pkgs, "didn't find any packages")

			for _, gopkg := range pkgs {
				assert.Len(t, gopkg.Errors, 0, "unexpected errors while loading packages")
			}
		})
	}
}

// Tests requiring Go SDK in runtime need testutil.RunWithGoSDK due to
// https://github.com/bazelbuild/rules_go/issues/2370.
func TestMain(m *testing.M) {
	testutil.RunWithGoSDK(m)
}
