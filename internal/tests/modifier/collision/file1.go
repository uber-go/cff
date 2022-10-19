//go:build cff
// +build cff

package collision

import (
	"context"

	"go.uber.org/cff"
)

// Flow1 is a very simple flow that returns 1
func Flow1() (int, error) {
	var i int
	err := cff.Flow(context.Background(),
		cff.Concurrency(1),
		cff.Results(&i),
		cff.Task(
			func() (int, error) {
				return 1, nil
			},
		),
	)
	return i, err
}
