// +build cff

package nestedparent

import (
	"context"

	"go.uber.org/cff"

	// When built under +cff build tag, this refers to the raw flow. After code
	// generation, this refers to the generated code.
	nested_child_gen "go.uber.org/cff/internal/tests/nested_child_gen"
)

// Parent is a CFF flow that uses a nested CFF flow.
func Parent(ctx context.Context, i int) (s string, err error) {
	err = cff.Flow(ctx,
		cff.Params(i),
		cff.Results(&s),

		cff.Task(nested_child_gen.Itoa),
	)

	return
}
