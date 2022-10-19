//go:build cff
// +build cff

package collision

import (
	"context"

	"go.uber.org/cff"
)

// Flow2 is a very simple flow that returns 2
func Flow2() (int, error) {
	var i int
	err := cff.Flow(context.Background(),
		cff.Concurrency(1),
		cff.Results(&i),
		cff.Task(
			func() (int, error) {
				return 2, nil
			},
		),
	)
	return i, err
}
