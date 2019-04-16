package internal

import (
	"fmt"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/types/typeutil"
)

func validateFlowCycles(f *flow, fset *token.FileSet) error {
	// Whether we've checked the subtree starting at this type for cycles
	// already.
	var visited typeutil.Map // map[types.Type]struct{}
	return findFlowCycles(f, &visited, fset)
}

// taskCyclePathEntry is an entry in the path as we walked the graph to detect cycles
type taskCyclePathEntry struct {
	// Tasks is the list of tasks in the order we visited them
	Task *task
	// Types is the type for which we visited the corresponding task in the Tasks list.
	Type types.Type
}

func prettyPrintTaskCycle(path []taskCyclePathEntry) string {
	str := fmt.Sprintf("need to run [%v] to provide %v (output)", path[0].Task.Sig, path[0].Type)
	for _, item := range path[1:] {
		str += fmt.Sprintf("\n\tneed to run [%v] to provide %v", item.Task.Sig, item.Type)
	}
	return str
}

func findFlowCycles(f *flow, visited *typeutil.Map, fset *token.FileSet) error {
	// If a flow has no Results eg a heatpipe flow, we need to do a DFS across all tasks. We are not
	// concerned with performance as this only happens once during compilation phase.
	for _, t := range f.Tasks {
		for _, dep := range t.Dependencies {
			if err := findFlowCyclesForType(f, nil /* path */, dep, visited, fset); err != nil {
				return err
			}
		}
	}
	return nil
}

func findFlowCyclesForType(f *flow, path []taskCyclePathEntry, t types.Type, visited *typeutil.Map,
	fset *token.FileSet) error {
	taskIdx, ok := f.providers.At(t).(int)
	if !ok {
		// This can happen if either cff.Params provides the type, but it can't introduce a cycle,
		// or we failed task validation and a provider is missing, so then we can't detect a cycle
		// until provider for this task exists again.
		return nil
	}

	task := f.Tasks[taskIdx]

	entry := taskCyclePathEntry{Task: task, Type: t}

	if len(path) > 0 {
		for _, p := range path {
			if types.Identical(p.Type, t) {
				return fmt.Errorf(
					"%v: cycle detected: %v",
					fset.Position(task.Node.Pos()),
					prettyPrintTaskCycle(append(path, entry)))
			}
		}
	}

	if visited.At(t) != nil {
		// Already checked this type for cycles.
		return nil
	}

	for _, dep := range task.Dependencies {
		if err := findFlowCyclesForType(f, append(path, entry), dep, visited, fset); err != nil {
			return err
		}
	}

	visited.Set(t, struct{}{})
	return nil
}
