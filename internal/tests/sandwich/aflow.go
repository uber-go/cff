//go:build cff
// +build cff

package sandwich

import (
	"context"

	"go.uber.org/cff"
)

func aFlow() (s string, err error) {
	err = cff.Flow(
		context.Background(),
		cff.Results(&s),
		cff.Task(
			aFunc,
		),
	)

	return s, err
}
