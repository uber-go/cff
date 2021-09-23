package example

import (
	"context"

	"go.uber.org/cff"
)

type cycle struct{}

type foo struct{} // a
type bar struct{} // b
type baz struct{} // c
type moo struct{} // d

func (c *cycle) Cycle(ctx context.Context) (res *moo, err error) {
	cff.Flow(ctx,
		cff.Results(&res),
		cff.Task(
			// b -> a
			func(b *bar) *foo {
				return &foo{}
			},
		),
		cff.Task(
			// a -> (c, b)
			func(a *foo) (*baz, *bar) {
				return &baz{}, &bar{}
			},
		),
		cff.Task(
			// c -> d*
			func(c *baz) *moo {
				return &moo{}
			},
		),
	)
	return
}
