// +build !cff

package foo

import (
	newctx "context"
	"sync"
)

func run(ctx newctx.Context) error {
	var result struct{}
	return func(ctx newctx.Context, v1 string) (err error) {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		var (
			wg0   sync.WaitGroup
			once0 sync.Once
		)

		wg0.Add(1)
		var v2 struct{}
		go func() {
			defer wg0.Done()

			v2 = func(string) struct{} {
				panic("don't call me")
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

		*(&result) = v2

		return err
	}(ctx, "foo")
}
