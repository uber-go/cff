package simple

import (
	"context"

	"go.uber.org/cff"
)

type bar struct{ i int64 }

// Flow is a very simple flow with some inputs and outputs.
func Flow() (int, error) {
	var cnt int
	err := cff.Flow(context.Background(),
		cff.Concurrency(2),
		cff.Results(&cnt),
		cff.Task(
			func() int64 {
				return int64(1)
			},
		),
		cff.Task(
			func(i int64) (*bar, error) {
				return &bar{i}, nil
			}),
		cff.Task(
			func(*bar) (int, error) {
				return 1, nil
			},
		),
	)
	return cnt, err
}
