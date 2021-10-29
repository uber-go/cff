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
	"go.uber.org/multierr"
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
	// File is the file where errors matching this case are found.
	File string
	// ErrorMatches is the error expected when compiling CFF code.
	ErrorMatches string
	// TestFuncs are test function names where errors matching this
	// case are found.
	TestFuncs []string
}

var codeGenerateFailCases = map[string][]errorCase{
	// map [directory name] -> list of test cases
	"bad-inputs": {
		{
			File:         "already-provided.go",
			ErrorMatches: "type string already provided at",
			TestFuncs:    []string{"AlreadyProvided"},
		},
		{
			File:         "cff-flow-arguments.go",
			ErrorMatches: "cff.Flow expects at least one function",
			TestFuncs:    []string{"ExpectsAtLeastOneArgument"},
		},
		{
			File:         "cff-flow-arguments.go",
			ErrorMatches: "expected cff call but got <nil>",
			TestFuncs:    []string{"FlowArgumentCallExpression2"},
		},
		{
			File:         "cff-flow-arguments.go",
			ErrorMatches: "expected cff call but got field ProvidesBad func\\(\\) go.uber.org/cff.Option",
			TestFuncs:    []string{"FlowArgumentNonCFF"},
		},
		{
			File:         "cff-flow-arguments.go",
			ErrorMatches: "expected a function call, got identifier",
			TestFuncs:    []string{"FlowArgumentsCallExpression"},
		},
		{
			File:         "cff-task-arguments.go",
			ErrorMatches: "expected function, got bool",
			TestFuncs:    []string{"ExpectsFunctionCallExpression"},
		},
		{
			File:         "cff-task-arguments.go",
			ErrorMatches: "only cff functions can be passed as task options",
			TestFuncs:    []string{"ExpectedFlowArgumentsSelectorExpression"},
		},
		{
			File:         "cff-task-arguments.go",
			ErrorMatches: "expected a function call, got identifier",
			TestFuncs:    []string{"ExpectedFlowArgumentsCallExpressions"},
		},
		{
			File:         "cff-task-arguments.go",
			ErrorMatches: `unexpected code generation directive "Instrument": only cff.Flow or cff.Parallel may be called at the top-level`,
			TestFuncs:    []string{"ExpectedFlowArgumentsCallExpressions", "ExpectedFlowArgumentsNotCFF"},
		},
		{
			File:         "cff-task-arguments.go",
			ErrorMatches: "only cff functions may be passed as task options: found package \"go.uber.org/cff/internal/failing_tests/bad-inputs\"",
			TestFuncs:    []string{"ExpectedFlowArgumentsNotCFF"},
		},
		{
			File:         "cff-task-arguments.go",
			ErrorMatches: "expected function, got untyped nil",
			TestFuncs:    []string{"ExpectedTasksBad"},
		},
		{
			File:         "cff-task-arguments.go",
			ErrorMatches: "expected function, got go.uber.org/cff.Option",
			TestFuncs:    []string{"ExpectedTasksBadCallExpr"},
		},
		{
			File:         "cff-task-arguments.go",
			ErrorMatches: "expected function, got int",
			TestFuncs:    []string{"ExpectedTasksBadCallExprNotCFF"},
		},
		{
			File:         "context-predicate.go",
			ErrorMatches: "only the first argument may be context.Context",
			TestFuncs:    []string{"ContextPredicate"},
		},
		{
			File:         "context-task.go",
			ErrorMatches: "only the first argument may be context.Context",
			TestFuncs:    []string{"ContextTask"},
		},
		{
			File:         "earlyresult.go",
			ErrorMatches: "unused output type string",
			TestFuncs:    []string{"EarlyResult"},
		},
		{
			File:         "earlyresult.go",
			ErrorMatches: "unused output type int32",
			TestFuncs:    []string{"EarlyResultDiamond"},
		},
		{
			File:         "earlyresult.go",
			ErrorMatches: "unused output type \\*go.uber.org/cff/internal/failing_tests/bad-inputs.quuz",
			TestFuncs:    []string{"EarlyResultMultipleFlows"},
		},
		{
			File:         "earlyresult.go",
			ErrorMatches: "unused output type \\*go.uber.org/cff/internal/failing_tests/bad-inputs.corge",
			TestFuncs:    []string{"EarlyResultMultipleFlows"},
		},
		{
			File:         "earlyresult.go",
			ErrorMatches: "unused output type \\*go.uber.org/cff/internal/failing_tests/bad-inputs.grault",
			TestFuncs:    []string{"EarlyResultMultipleFlows"},
		},
		{
			File:         "error-task.go",
			ErrorMatches: "only the last result may be an error",
			TestFuncs:    []string{"EarlyResultMultipleFlows"},
		},
		{
			File:         "fallback-with.go",
			ErrorMatches: "cff.FallbackWith result at position 1 of type string cannot be used as bool",
			TestFuncs:    []string{"FallbackWithTypeMismatch"},
		},
		{
			File:         "fallback-with.go",
			ErrorMatches: "cff.FallbackWith result at position 2 of type bool cannot be used as string",
			TestFuncs:    []string{"FallbackWithTypeMismatch"},
		},
		{
			File:         "fallback-with.go",
			ErrorMatches: "cff.FallbackWith must produce the same number of results as the task: expected 2, got 1",
			TestFuncs:    []string{"FallbackWithBadPositionalArguments"},
		},
		{
			File:         "fallback-with.go",
			ErrorMatches: "Task must return an error for FallbackWith to be used",
			TestFuncs:    []string{"FallbackWithNoError"},
		},
		{
			File:         "instrument.go",
			ErrorMatches: "cff.Instrument requires a cff\\.Emitter to be provided: use cff\\.WithEmitter",
			TestFuncs:    []string{"MissingCFFLoggerAndMetrics"},
		},

		{
			File:         "missing-provider.go",
			ErrorMatches: "no provider found for float64",
			TestFuncs:    []string{"MissingProvider"},
		},
		{
			File:         "nonpointer-result.go",
			ErrorMatches: "invalid parameter to cff.Results: expected pointer, got bool",
			TestFuncs:    []string{"ResultsNonPointer"},
		},
		{
			File:         "parallel.go",
			ErrorMatches: "the only allowed argument is a single context.Context parameter",
			TestFuncs:    []string{"ParallelInvalidParamsType", "ParallelInvalidFuncVar"},
		},
		{
			File:         "parallel.go",
			ErrorMatches: "only the first argument may be context.Context",
			TestFuncs:    []string{"ParallelInvalidParamsMultiple"},
		},
		{
			File:         "parallel.go",
			ErrorMatches: "the only allowed return value is an error",
			TestFuncs:    []string{"ParallelInvalidReturnType"},
		},
		{
			File:         "predicate.go",
			ErrorMatches: "the function must return a single boolean result",
			TestFuncs:    []string{"PredicateReturnsNonbool", "PredicateReturnsMultipleValues"},
		},
		{
			File:         "predicate-params.go",
			ErrorMatches: "cff.Predicate expected a function but received",
			TestFuncs:    []string{"PredicateParams"},
		},
		{
			File:         "unused-outputs.go",
			ErrorMatches: "unused output type bool",
			TestFuncs:    []string{"DisconnectedSubgraph"},
		},
		{
			File:         "unused-outputs.go",
			ErrorMatches: "unused output type uint32",
			TestFuncs:    []string{"DisconnectedSubgraphPredicate"},
		},
		{
			File:         "top-level-flow.go",
			ErrorMatches: "unexpected code generation directive \"Predicate\"",
			TestFuncs:    []string{"BadTopLevelFunction"},
		},
		{
			File:         "unused-inputs.go",
			ErrorMatches: "unused input type string",
			TestFuncs:    []string{"UnusedInputs"},
		},
		{
			File:         "unused-task.go",
			ErrorMatches: "cff\\.Invoke cannot be provided on a Task that produces values besides errors",
			TestFuncs:    []string{"UnsupportedInvoke"},
		},
		{
			File:         "unused-task.go",
			ErrorMatches: "task must return at least one non-error value but currently produces zero.",
			TestFuncs:    []string{"NoInvokeNoResults"},
		},
		{
			File:         "unused-task-error.go",
			ErrorMatches: "task must return at least one non-error value but currently produces zero.",
			TestFuncs:    []string{"NoInvokeWithError"},
		},
		{
			File:         "variadic.go",
			ErrorMatches: "variadic functions are not yet supported",
			TestFuncs:    []string{"Variadic", "VariadicPredicate"},
		},
	},
	"cycles": {
		{
			File:         "easy-cycle.go",
			ErrorMatches: "cycle detected: need to run \\[func\\(int64\\) string\\] to provide string",
			TestFuncs:    []string{"EasyCycle"},
		},
		{
			File:         "no-output.go",
			ErrorMatches: "cycle detected: need to run \\[func\\(int32\\) string\\] to provide string",
			TestFuncs:    []string{"EasyCycleNoOut"},
		},
	},
}

// This works by trying to transpile all of the files in the directory/package,
// accumulating all of the errors, and then checking if in a file there is an
// instance of the error we are looking for.
// Note: error accumulation is per-package so at the moment state is kept
// when running transpiler across many flows expected to fail.
func TestCodeGenerateFails(t *testing.T) {
	// fileError hashes a file and errMatch for map lookups.
	type fileError struct {
		file     string
		errMatch string
	}

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
						errors = append(errors, multierr.Errors(err)...)
					}
				}
				for _, err := range errors {
					t.Logf("found error %q", err.Error())
				}

				matchSet := make(map[fileError]struct{})
				observedErrors := make(map[fileError]int)
				for _, match := range errCases {
					fe := fileError{
						file:     match.File,
						errMatch: match.ErrorMatches,
					}
					_, ok := matchSet[fe]
					require.False(t, ok, "match case %q in file %q should be unique", match.ErrorMatches, match.File)
					matchSet[fe] = struct{}{}

					regexpError := regexp.MustCompile(fmt.Sprintf("%s.*%s", match.File, match.ErrorMatches))
					for _, err := range errors {
						// Unwrap combined errors produced by the cff
						// compiler.
						if ok := regexpError.MatchString(err.Error()); ok {
							observedErrors[fileError{
								file:     match.File,
								errMatch: match.ErrorMatches,
							}]++
						}
					}
				}
				for _, cse := range errCases {
					fe := fileError{
						file:     cse.File,
						errMatch: cse.ErrorMatches,
					}
					assert.Equal(
						t,
						len(cse.TestFuncs),
						observedErrors[fe],
						"incorrect matches for %q in file %q",
						cse.ErrorMatches,
						cse.File,
					)
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
