package internal

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages/packagestest"
)

func TestCompileParallel(t *testing.T) {
	cffModule := packagestest.Module{
		Name:  "go.uber.org/cff",
		Files: packagestest.MustCopyFileTree("./.."),
	}
	modules := []packagestest.Module{cffModule}
	setups := setupCompilers(t, filepath.Join(internalTests, "compile_tests/parallel"), modules)
	require.Len(t, setups, 1)

	for _, setup := range setups {
		compiled := setup.compiler.compileFile(setup.file, setup.pkg)
		require.Len(t, compiled.Parallels, 2)
		set := make(map[int]struct{})
		for _, p := range compiled.Parallels {
			for _, task := range p.Tasks {
				t.Run("serial incremented", func(t *testing.T) {
					_, ok := set[task.Serial]
					assert.False(t, ok)
					set[task.Serial] = struct{}{}
				})
				t.Run("parallel tasks are independent", func(t *testing.T) {
					assert.Empty(t, task.Function.Dependencies)
					assert.Empty(t, task.Function.DependsOn)
				})
			}
		}
	}
}
