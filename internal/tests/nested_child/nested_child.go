//go:build cff
// +build cff

package nestedchild

import (
	"context"
	"strconv"

	"go.uber.org/cff"
)

// Itoa is a flow that is simply used by another flow.
func Itoa(ctx context.Context, i int) (s string, err error) {
	err = cff.Flow(ctx,
		cff.Params(i),
		cff.Results(&s),

		cff.Task(func(i int) string {
			return strconv.Itoa(i)
		}),
	)

	return
}
