// +build !cff

package foo

import (
	newctx "context"
)

func run(ctx newctx.Context) error {
	var result struct{}
	return func(ctx newctx.Context, v1 string) (err error) {

		if ctx.Err() != nil {
			return ctx.Err()
		}

		var v2 struct{}
		v2 = func(string) struct{} {
			panic("don't call me")
		}(v1)

		*(&result) = v2

		return err
	}(ctx, "foo")
}
