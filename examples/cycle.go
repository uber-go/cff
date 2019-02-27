// +build cff

package example

import (
	"go.uber.org/cff"
	"context"
)

type cycle struct{}

type Foo struct{} // a
type Bar struct{} // b
type Baz struct{} // c
type Moo struct{} // d

func (c *cycle) Cycle(ctx context.Context) (res *Moo, err error) {
	cff.Flow(ctx,
		cff.Result(&res),

		cff.Tasks(
			// b -> a
			func(b *Bar) *Foo {
				return &Foo{}
			},
			// a -> (c, b)
			func(a *Foo) (*Baz, *Bar) {
				return &Baz{}, &Bar{}
			},
			// c -> d*
			func(c *Baz) *Moo {
				return &Moo{}
			},
		),
	)
	return
}
