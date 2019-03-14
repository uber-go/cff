// +build !cff

package fallbackwith

import (
	"context"
	"sync"
)

// Serial executes a flow that fails with the given error, if any and recovers
// with the given string.
func Serial(e error, r string) (string, error) {
	var s string
	err := func(ctx context.Context) (err error) {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		var (
			wg0   sync.WaitGroup
			once0 sync.Once
		)

		wg0.Add(1)
		var v1 string
		var err0 error
		go func() {
			defer wg0.Done()

			v1, err0 = func() (string, error) {
				return "foo", e
			}()
			if err0 != nil {

				v1, err0 = r, nil
			}

		}()

		wg0.Wait()
		if err != nil {

			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once0
			_ = &v1
		)

		*(&s) = v1

		return err
	}(context.Background())
	return s, err
}
