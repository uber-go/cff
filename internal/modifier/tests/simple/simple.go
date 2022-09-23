package simple

import (
	"context"

	"go.uber.org/cff"
)

type bar struct{ i int64 }

// Flow is a very simple flow with some inputs and outputs.
func Flow() (int, string, error) {
	var (
		iRes int
		sRes string
	)
	err := cff.Flow(context.Background(),
		cff.Concurrency(2),
		cff.Results(&iRes, &sRes),
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
		cff.Task(
			func(i int) (string, error) {
				if i != 0 {
					return "non-zero", nil
				}
				return "zero", nil
			},
		),
	)
	return iRes, sRes, err
}
