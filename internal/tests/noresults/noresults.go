// +build cff

package noresults

import (
	"context"
	"errors"
	"fmt"
	"os"

	"go.uber.org/cff"
)

func main() {
	h := &H{}
	ctx := context.Background()
	err := h.Swallow(ctx, os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(os.Args, "is swallowed")
}

// H is exported for tests.
type H struct{}

// Swallow tests that error is not swallowed..
func (h *H) Swallow(ctx context.Context, req string) (err error) {
	err = cff.Flow(ctx,
		cff.Params(req),
		cff.Task(
			func(s string) error {
				if s == "tide pods" {
					return errors.New("can not Swallow")
				}
				return nil
			},
		),
		cff.Task(func(string) {}),
	)
	return
}

// TripleSwallow tests that no error is returned and flow runs.
func (h *H) TripleSwallow(ctx context.Context, req string) (err error) {
	err = cff.Flow(ctx,
		cff.Params(req),
		cff.Task(
			func(string) {},
		),
		cff.Task(
			func(string) {}),
		cff.Task(
			func(string) {},
		),
	)
	return
}
