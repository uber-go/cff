package internal

import (
	"flag"
	"fmt"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/packages/packagestest"
)

var generateFlag = flag.Bool("generate", false, "Update generated code rather than verifying")

const (
	goldenTestDir            = "tests" // internal/tests
	goldenTestImportInternal = "go.uber.org/cff/internal"
)

type testCase struct {
	Name        string
	Path        string //  <$WORKSPACE>/code.uber.internal/mar/../basic.go
	FileName    string // basic.go
	GenFilePath string
}

func discoverTestCases(t *testing.T, dir string) []*testCase {
	// Adding a go executable to the PATH for Bazel.
	_ = os.Setenv("PATH", os.ExpandEnv("$TEST_SRCDIR/__main__/external/go_sdk/bin:$PATH"))
	// See if we are in Bazel environment as `go test` doesn't need GOCACHE to be set manually.
	if file, err := os.Stat(os.Getenv("TEST_TMPDIR")); err == nil {
		if file.IsDir() {
			// Go executable requires a GOCACHE to be set after go1.12.
			_ = os.Setenv("GOCACHE", filepath.Join(os.Getenv("TEST_TMPDIR"), "/cache"))
		}
	}

	infos, err := ioutil.ReadDir(dir)
	require.NoErrorf(t, err, "failed to ls %q", dir)
	// TODO(T2982585): Remove this deny list.
	denyList := map[string]bool{
		"instrument": true,
		"panic":      true,
	}
	var tests []*testCase
	for _, info := range infos {
		if !info.IsDir() {
			continue
		}
		if _, ok := denyList[info.Name()]; ok {
			continue
		}
		subinfos, err := ioutil.ReadDir(filepath.Join(dir, info.Name()))
		require.NoErrorf(t, err, "failed to ls %q", dir)
		for _, subinfo := range subinfos {
			if strings.TrimSuffix(filepath.Base(subinfo.Name()), ".go") == info.Name() {

				tests = append(tests, &testCase{
					Name:        info.Name(),
					Path:        filepath.Join(goldenTestDir, filepath.Join(info.Name(), info.Name())+".go"),
					FileName:    info.Name() + ".go",
					GenFilePath: filepath.Join(goldenTestDir, filepath.Join(info.Name(), info.Name())+"_gen.go"),
				})
			}
		}
	}

	return tests
}

func TestCodeIsUpToDate(t *testing.T) {
	defer func() {
		if t.Failed() {
			t.Log("Try re-running with --generate")
		}
	}()
	tests := discoverTestCases(t, goldenTestDir)

	// Need as we have cff import inside our source files.
	cffModule := packagestest.Module{
		Name:  "go.uber.org/cff",
		Files: packagestest.MustCopyFileTree("./.."),
	}
	// TODO: need zap and tally dependencies to make tests work here too.
	// Making a requirement that we only allow one generated test file per package.
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Create temp directory to output, something in /var/folders/.../cff-test{DIGITS}
			tempDir, err := ioutil.TempDir("", "cff-test")
			require.NoError(t, err)
			defer func() {
				assert.NoError(t, os.RemoveAll(tempDir))
			}()
			// Try and represent our Go module.
			mod := packagestest.Module{
				// Note: module's name doesn't really matter since we aren't importing it.
				Name: "go.uber.org/cff/internal/tests/" + tt.Name,
				Files: map[string]interface{}{
					tt.FileName:         packagestest.Copy(tt.Path),
					tt.Name + "_gen.go": packagestest.Copy(tt.GenFilePath),
				},
			}
			assert.NotNil(t, mod)
			// Order matters: we want our test package to be primarymod and include
			// references to cffModule.
			exp := packagestest.Export(t, packagestest.Modules, []packagestest.Module{mod, cffModule})
			fset := token.NewFileSet()

			cfg := exp.Config

			cfg.Mode = packages.LoadSyntax
			cfg.BuildFlags = []string{"-tags=cff"}
			cfg.Fset = fset

			defer exp.Cleanup()

			// Using pattern for go test not to run _test unit tests which test generated code.
			pkgs, err := packages.Load(exp.Config, "pattern="+tt.FileName)
			require.NoError(t, err, "could not load packages")
			require.NotEmpty(t, pkgs, "didn't find any packages")
			for _, pkg := range pkgs {

				for _, err := range pkg.Errors {
					require.Empty(t, err.Msg)
				}
				// Full path to source files discovered in the package.
				src := pkg.GoFiles[0]
				// Directory containing the entire package.
				wantDir := filepath.Dir(src)
				// Looking for our pregenerated file inside the directory.
				newName := strings.TrimSuffix(filepath.Base(src), ".go") + "_gen.go"
				// File that we output by running through Process.
				outputFile := filepath.Join(tempDir, newName)
				// The core of tests: run through our transpiler.
				assert.NoErrorf(t, Process(fset, pkg, outputFile),
					"failed to process package %v", pkg.PkgPath)
				// Our "golden" file we are comparing against.
				wantFile := filepath.Join(wantDir, filepath.Base(outputFile))

				if *generateFlag {
					if !assert.NoErrorf(t, copyFile(t, outputFile, wantFile), "could not copy generated %q for %q", filepath.Base(src), pkg.PkgPath) {
						continue
					}
				}

				wantFileContents, err := ioutil.ReadFile(wantFile)
				require.NoErrorf(t, err, "could not read %q", wantFile)

				gotFileContents, err := ioutil.ReadFile(outputFile)
				require.NoErrorf(t, err, "could not read %q", outputFile)

				// assert.Equal on string prints a diff if the contents
				// don't match.
				assert.Equal(t, string(wantFileContents), string(gotFileContents),
					"contents of %q don't match", newName)
			}
		})
	}
}

func copyFile(t *testing.T, src, dst string) error {
	srcF, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("could not open %q", src)
	}
	defer func() {
		assert.NoError(t, srcF.Close())
	}()

	dstF, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("could not open %q", dst)
	}
	defer func() {
		assert.NoError(t, dstF.Close())
	}()

	if _, err := io.Copy(dstF, srcF); err != nil {
		return fmt.Errorf("could not copy %q to %q", src, dst)
	}

	return nil
}
