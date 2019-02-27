// +build cff

package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"go.uber.org/cff"
	"github.com/uber-go/tally"
)

func main() {
	scope := tally.NoopScope
	h := &h{scope: scope}
	ctx := context.Background()
	res, err := h.run(ctx, os.Args[1])
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d\n", res)
}

type h struct {
	scope tally.Scope
}

func (h *h) run(ctx context.Context, req string) (res int, err error) {
	err = cff.Flow(ctx,
		cff.Provide(req),
		cff.Result(&res),
		cff.Scope(h.scope),
		cff.InstrumentFlow("AtoiRun"),

		cff.Task(
			strconv.Atoi,
			cff.Instrument("Atoi"),
		),
	)
	return
}
