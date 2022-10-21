//go:build failing
// +build failing

package badinputs

import (
	"context"
	"strconv"

	"go.uber.org/cff"
)

func FlowWithoutTag() error {
	var message string
	return cff.Flow(context.Background(),
		cff.Params(1),
		cff.Results(&message),
		cff.Task(
			func(i int) string {
				return strconv.Itoa(i)
			},
		),
	)
}

func ParallelWithoutTag() error {
	return cff.Parallel(context.Background(),
		cff.Slice(
			func(int, string) {},
			[]string{"foo", "bar"},
		),
	)
}
