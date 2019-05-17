// +build !cff
// @generated

package foo

import (
	newctx "context"
	"fmt"
	"sync"
)

func run(ctx newctx.Context) error {
	var result struct{}
	return func(ctx newctx.Context, v1 string) (err error) {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		var (
			once0 sync.Once
		)

		var v2 struct{}

		func() {

			defer func() {
				recovered := recover()
				if recovered != nil {
					once0.Do(func() {
						recoveredErr := fmt.Errorf("task panic: %v", recovered)

						err = recoveredErr
					})
				}
			}()

			v2 = func(string) struct{} {
				panic("don't call me")
			}(v1)

		}()

		if err != nil {

			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once0
			_ = &v2
		)

		*(&result) = v2

		return err
	}(ctx, "foo")
}
