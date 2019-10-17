// +build cff

package badinputs

import (
	"context"

	"go.uber.org/cff"
)

// DisconnectedSubgraph does not use cff.Invoke on a disconnected subgraph.
func DisconnectedSubgraph() {
	var s string

	unused1 := func(b int) float64 { return 0 }
	unused2 := func(f float64) bool { return true }

	cff.Flow(
		context.Background(),
		cff.Results(&s),

		cff.Task(unused1),
		cff.Task(unused2),

		cff.Task(func() string {
			return ""
		}),
	)
}

// DisconnectedSubgraphPredicate does not use cff.Invoke on a disconnected subgraph that has a dependency using Predicate.
func DisconnectedSubgraphPredicate() {
	var s string

	unused1 := func(b int) float64 { return 0 }
	unused2pred := func(f float64) bool { return true }
	unused2 := func() uint32 { return uint32(0) }

	cff.Flow(
		context.Background(),
		cff.Results(&s),

		cff.Task(unused1),
		cff.Task(unused2, cff.Predicate(unused2pred)),

		cff.Task(func() string {
			return ""
		}),
	)
}
