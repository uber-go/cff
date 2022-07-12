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

	tt := []struct {
		desc         string
		numTask      int
		numSliceTask int
		numMapTask   int
	}{
		{
			desc:         "first parallel",
			numTask:      3,
			numSliceTask: 2,
			numMapTask:   1,
		},
		{
			desc:         "second parallel",
			numTask:      3,
			numSliceTask: 1,
			numMapTask:   1,
		},
	}

	for _, setup := range setups {
		compiled := setup.compiler.compileFile(setup.file, setup.pkg)
		require.Len(t, compiled.Parallels, 2)

		t.Run("task serial incremented", func(t *testing.T) {
			for i, p := range compiled.Parallels {
				serialSet := make(map[int]struct{})
				for _, task := range p.Tasks {
					_, ok := serialSet[task.Serial]
					assert.False(t, ok)
					serialSet[task.Serial] = struct{}{}
				}
				assert.Len(t, p.Tasks, tt[i].numTask)
			}
		})

		t.Run("parallel tasks are independent", func(t *testing.T) {
			for _, p := range compiled.Parallels {
				for _, task := range p.Tasks {
					assert.Empty(t, task.Function.Dependencies)
					assert.Empty(t, task.Function.DependsOn)
				}
			}
		})

		t.Run("slice serial incremented", func(t *testing.T) {
			for i, p := range compiled.Parallels {
				serialSet := make(map[int]struct{})
				for _, task := range p.SliceTasks {
					_, ok := serialSet[task.Serial]
					assert.False(t, ok)
					serialSet[task.Serial] = struct{}{}
				}
				assert.Len(t, p.SliceTasks, tt[i].numSliceTask)
			}
		})

		t.Run("map serial incremented", func(t *testing.T) {
			for i, p := range compiled.Parallels {
				serialSet := make(map[int]struct{})
				for _, task := range p.MapTasks {
					_, ok := serialSet[task.Serial]
					assert.False(t, ok)
					serialSet[task.Serial] = struct{}{}
				}
				assert.Len(t, p.MapTasks, tt[i].numMapTask)
			}
		})
	}
}
