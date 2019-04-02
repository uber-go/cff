// +build !cff
// @generated

package nestedchild

import (
	"context"
	"fmt"
	"strconv"
	"sync"
)

// Itoa is a flow that is simply used by another flow.
func Itoa(ctx context.Context, i int) (s string, err error) {
	err = func(ctx context.Context, v1 int) (err error) {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		var (
			wg0   sync.WaitGroup
			once0 sync.Once
		)

		wg0.Add(1)

		var v2 string

		go func() {
			defer wg0.Done()

			defer func() {
				recovered := recover()
				if recovered != nil {
					once0.Do(func() {
						recoveredErr := fmt.Errorf("task panic: %v", recovered)

						err = recoveredErr
					})
				}
			}()

			v2 = func(i int) string {
				return strconv.Itoa(i)
			}(v1)

		}()

		wg0.Wait()
		if err != nil {

			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once0
			_ = &v2
		)

		*(&s) = v2

		return err
	}(ctx, i)

	return
}
