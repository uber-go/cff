// +build !cff

package recoverwith

import (
	"context"
)

// Serial executes a flow that fails with the given error, if any and recovers
// with the given string.
func Serial(e error, r string) (string, error) {
	var s string
	err := func(ctx context.Context) (err error) {

		if ctx.Err() != nil {
			return ctx.Err()
		}

		var v1 string
		v1, err = func() (string, error) {
			return "foo", e
		}()
		if err != nil {

			v1, err = r, nil
		}

		*(&s) = v1

		return err
	}(context.Background())
	return s, err
}
