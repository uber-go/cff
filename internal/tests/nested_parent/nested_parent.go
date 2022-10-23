//go:build cff
// +build cff

package nestedparent

import (
	"context"

	"go.uber.org/cff"

	// When built under +cff build tag, this refers to the raw flow. After code
	// generation, this refers to the generated code.
	nestedchild "go.uber.org/cff/internal/tests/nested_child"
)

// Parent is a cff flow that uses a nested cff flow.
func Parent(ctx context.Context, i int) (s string, err error) {
	err = cff.Flow(ctx,
		cff.Params(i),
		cff.Results(&s),
		cff.Task(nestedchild.Itoa),
	)
	return
}
