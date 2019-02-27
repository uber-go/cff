// +build cff

package recoverwith

import (
	"context"

	"go.uber.org/cff"
)

// Serial executes a flow that fails with the given error, if any and recovers
// with the given string.
func Serial(e error, r string) (string, error) {
	var s string
	err := cff.Flow(
		context.Background(),
		cff.Result(&s),
		cff.Task(func() (string, error) {
			return "foo", e
		}, cff.RecoverWith(r)),
	)
	return s, err
}
