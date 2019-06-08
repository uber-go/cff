// +build cff

package sandwich

import (
	"context"

	"go.uber.org/cff"
)

func bFlow() (s string, err error) {
	err = cff.Flow(
		context.Background(),
		cff.Results(&s),
		cff.Task(
			aFunc,
		),
	)

	return s, err
}
