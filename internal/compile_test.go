package internal

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/cff/internal/pkg"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/packages/packagestest"
	"golang.org/x/tools/go/types/typeutil"
)

func TestNoOutputTypes(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatalf("expected mustSetNoOutputProvider to panic")
		}
	}()
	providers := &typeutil.Map{}
	f := &flow{
		providers: providers,
	}
	noOutput := f.addNoOutput()
	task := &task{
		Function:   &function{},
		invokeType: noOutput,
	}
	f.mustSetNoOutputProvider(task.Function, 0)
	f.mustSetNoOutputProvider(task.Function, 0)
	assert.Equal(t, task.invokeType, f.providers.At(task.invokeType))
}

// toCompile contains a CFF source file, compiler, and Package that the source
// file can be compiled with. This assortment is a convenience for invoking the
// CFF compiler in test.
type toCompile struct {
	file     *ast.File
	compiler *compiler
	pkg      *pkg.Package
}

// setupCompilers loads a collection of CFF source files and compilers for
// the source files that can be used to invoke and test compiler behaviour.
func setupCompilers(
	t *testing.T,
	pattern string,
	modules []packagestest.Module,
) map[string]toCompile {
	exp := packagestest.Export(
		t, packagestest.GOPATH,
		modules,
	)
	fset := token.NewFileSet()

	cfg := exp.Config
	cfg.BuildFlags = []string{"-tags=cff"}
	cfg.Fset = fset
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
	pkgs, err := packages.Load(
		exp.Config,
		pattern)
	require.NoError(t, err, "could not load packages")
	require.NotEmpty(t, pkgs, "didn't find any packages")

	toCompileMap := make(map[string]toCompile)
	for _, gopkg := range pkgs {
		pkg := newPackage(gopkg)
		for i, path := range pkg.CompiledGoFiles {
			c := newCompiler(compilerOpts{
				Fset:    cfg.Fset,
				Info:    pkg.TypesInfo,
				Package: pkg.Types,
			})
			toCompileMap[path] = toCompile{
				file:     pkg.Syntax[i],
				compiler: c,
				pkg:      pkg,
			}
		}
	}
	return toCompileMap
}

func newPackage(p *packages.Package) *pkg.Package {
	return &pkg.Package{
		CompiledGoFiles: p.CompiledGoFiles,
		Syntax:          p.Syntax,
		Types:           p.Types,
		TypesInfo:       p.TypesInfo,
	}
}

// TestCompileFile_Predicate tests that the internal state of the CFF compiler
// is correct after compiling cff.Predicates.
func TestCompileFile_Predicate(t *testing.T) {
	cffModule := packagestest.Module{
		Name:  "go.uber.org/cff",
		Files: packagestest.MustCopyFileTree("./.."),
	}
	modules := []packagestest.Module{cffModule}
	setups := setupCompilers(t, filepath.Join(internalTests, "compile_tests/predicate/..."), modules)
	for _, c := range setups {
		file := c.compiler.compileFile(c.file, c.pkg)
		for _, flow := range file.Flows {
			t.Run("every predicate type is a dependency of a task function", func(t *testing.T) {
				set := new(typeutil.Map)
				for _, p := range flow.predicateTypes {
					set.Set(p, struct{}{})
				}
				for _, topoFunc := range flow.TopoFuncs {
					for _, depType := range topoFunc.Dependencies {
						set.Delete(depType)
					}
				}
				assert.Equal(t, 0, set.Len())
			})
			t.Run("a function's Dependencies include all of that function's inputs", func(t *testing.T) {
				for _, fun := range flow.Funcs {
					inSet := new(typeutil.Map)
					for _, in := range fun.inputs() {
						inSet.Set(in, struct{}{})
					}
					for _, dep := range fun.Dependencies {
						inSet.Delete(dep)
					}
					assert.Equal(t, 0, inSet.Len())
				}
			})
			t.Run("compiled flow TopoFuncs include all Tasks and Predicates", func(t *testing.T) {
				set := map[*function]struct{}{}
				for _, fun := range flow.TopoFuncs {
					set[fun] = struct{}{}
				}
				for _, task := range flow.Tasks {
					delete(set, task.Function)
				}
				for _, pred := range flow.Predicates {
					delete(set, pred.Function)
				}
				assert.Len(t, set, 0)
			})
			t.Run("after dependency resolution, predicates come before their tasks", func(t *testing.T) {
				taskSet := map[*task]struct{}{}
				for _, fun := range flow.TopoFuncs {
					if fun.Task != nil {
						taskSet[fun.Task] = struct{}{}
					}
					if fun.Predicate != nil {
						_, ok := taskSet[fun.Predicate.Task]
						assert.False(t, ok)
					}
				}
			})
			t.Run("functions must have either a task or predicate", func(t *testing.T) {
				for _, fun := range flow.TopoFuncs {
					if fun.Task != nil && fun.Predicate != nil {
						t.Error("function has a Task and Predicate set.")
					}
				}
			})
		}
	}
}
