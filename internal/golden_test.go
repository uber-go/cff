package internal

import (
	"flag"
	"fmt"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
)

var generateFlag = flag.Bool("generate", false, "Update generated code rather than verifying")

const (
	goldenTestDir            = "tests" // internal/tests
	goldenTestImportInternal = "go.uber.org/cff/internal"
)

type testCase struct {
	Name       string
	Dir        string
	ImportPath string
}

func discoverTestCases(t *testing.T, dir string) []*testCase {
	infos, err := ioutil.ReadDir(dir)
	require.NoErrorf(t, err, "failed to ls %q", dir)

	var tests []*testCase
	for _, info := range infos {
		if !info.IsDir() {
			continue
		}

		tests = append(tests, &testCase{
			Name:       info.Name(),
			ImportPath: filepath.Join(goldenTestImportInternal, goldenTestDir, info.Name()),
		})
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
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
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
			}, tt.ImportPath+"/...")
			require.NoError(t, err, "could not load packages")
			require.NotEmpty(t, pkgs, "didn't find any packages")

			for _, pkg := range pkgs {
				relativePath, err := filepath.Rel(tt.ImportPath, pkg.PkgPath)
				require.NoErrorf(t, err, "could not determine relative path to %q from %q",
					pkg.PkgPath, tt.ImportPath)

				wantDir := filepath.Dir(pkg.GoFiles[0])
				gotDir := filepath.Join(tempDir, relativePath)

				assert.NoErrorf(t, Process(fset, pkg, gotDir),
					"failed to process package %v", pkg.PkgPath)

				gotInfos, err := ioutil.ReadDir(gotDir)
				require.NoErrorf(t, err, "failed to ls %q", gotDir)

				for _, info := range gotInfos {
					wantFile := filepath.Join(wantDir, info.Name())
					gotFile := filepath.Join(gotDir, info.Name())

					// Copy the file over before checking its contents.
					if *generateFlag {
						if !assert.NoErrorf(t, copyFile(t, gotFile, wantFile), "could not copy generated %q for %q", info.Name(), pkg.PkgPath) {
							continue
						}
					}

					wantFileContents, err := ioutil.ReadFile(wantFile)
					require.NoErrorf(t, err, "could not read %q", wantFile)

					gotFileContents, err := ioutil.ReadFile(gotFile)
					require.NoErrorf(t, err, "could not read %q", gotFile)

					// assert.Equal on string prints a diff if the contents
					// don't match.
					assert.Equal(t, string(wantFileContents), string(gotFileContents),
						"contents of %q don't match", info.Name())
				}
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
