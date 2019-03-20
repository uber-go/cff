// +build cff

package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"go.uber.org/cff"
)

func main() {
	h := &h{}
	ctx := context.Background()
	err := h.swallow(ctx, os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(os.Args, "is swallowed")
}

type h struct{}

func (h *h) swallow(ctx context.Context, req string) (err error) {
	err = cff.Flow(ctx,
		cff.Params(req),
		cff.Task(
			func(s string) error {
				if s == "tide pods" {
					return errors.New("can not swallow")
				}
				return nil
			},
		),
		cff.Task(func(string) {}),
	)
	return
}

func (h *h) tripleSwallow(ctx context.Context, req string) (err error) {
	err = cff.Flow(ctx,
		cff.Params(req),
		cff.Tasks(
			func(string) {},
			func(string) {},
			func(string) {},
		),
	)
	return
}
