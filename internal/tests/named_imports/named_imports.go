// +build cff

package namedimports

import (
	newctx "context"

	cffv2 "go.uber.org/cff"
)

func run(ctx newctx.Context) error {
	var result struct{}
	return cffv2.Flow(ctx,
		cffv2.Params("foo"),
		cffv2.Results(&result),
		cffv2.Task(
			func(string) struct{} {
				panic("don't call me")
			},
		),
	)
}
