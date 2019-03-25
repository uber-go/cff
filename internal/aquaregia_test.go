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

func TestCodeGenerateFails(t *testing.T) {
	// map [directory name] -> list of test cases
	errorCasesByDirectory := map[string][]errorCase{
		"basic": {
			{
				File:         "basic.go",
				ErrorMatches: "no provider found for float64",
			},
		},
	}

	for testDirectoryName, errCases := range errorCasesByDirectory {
		t.Run(fmt.Sprintf("test cases for directory %q", testDirectoryName), func(t *testing.T) {
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
