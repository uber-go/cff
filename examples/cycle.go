// +build cff

package example

import (
	"go.uber.org/cff"
	"context"
)

type cycle struct{}

type foo struct{} // a
type bar struct{} // b
type baz struct{} // c
type moo struct{} // d

func (c *cycle) Cycle(ctx context.Context) (res *moo, err error) {
	cff.Flow(ctx,
		cff.Results(&res),

		cff.Tasks(
			// b -> a
			func(b *bar) *foo {
				return &foo{}
			},
			// a -> (c, b)
			func(a *foo) (*baz, *bar) {
				return &baz{}, &bar{}
			},
			// c -> d*
			func(c *baz) *moo {
				return &moo{}
			},
		),
	)
	return
}
