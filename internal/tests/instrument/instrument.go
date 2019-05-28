// +build cff

package instrument

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"go.uber.org/cff"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

func main() {
	scope := tally.NoopScope
	logger := zap.NewNop()
	h := &H{
		Scope:  scope,
		Logger: logger,
	}
	ctx := context.Background()
	res, err := h.Run(ctx, os.Args[1])
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d\n", res)
}

// H is used by other tests.
type H struct {
	Scope  tally.Scope
	Logger *zap.Logger
}

// Run executes a flow to test instrumentation.
func (h *H) Run(ctx context.Context, req string) (res uint8, err error) {
	err = cff.Flow(ctx,
		cff.Params(req),
		cff.Results(&res),
		cff.Metrics(h.Scope),
		cff.Logger(h.Logger),
		cff.InstrumentFlow("AtoiRun"),

		cff.Task(
			strconv.Atoi,
			cff.Instrument("Atoi"),
		),

		cff.Task(
			func(i int) (uint8, error) {
				if i > -1 && i < 256 {
					return uint8(i), nil
				}
				return 0, errors.New("int can not fit into 8 bits")
			},
			cff.FallbackWith(uint8(0)),
			cff.Instrument("uint8"),
		),
	)
	return
}

// Do executes a flow to test instrumentation.
func (h *H) Do(ctx context.Context, req string) (res int, err error) {
	err = cff.Flow(ctx,
		cff.Params(req),
		cff.Results(&res),
		cff.InstrumentFlow("AtoiDo"),
		cff.Metrics(h.Scope),
		cff.Logger(h.Logger),
		cff.Task(
			strconv.Atoi,
			cff.Instrument("Atoi"),
		),
	)
	return
}

// Work executes a flow to test instrumentation.
func (h *H) Work(ctx context.Context, req string) (res int, err error) {
	err = cff.Flow(ctx,
		cff.Params(req),
		cff.Results(&res),
		cff.Metrics(h.Scope),
		cff.Logger(h.Logger),
		cff.Task(
			strconv.Atoi,
			cff.Instrument("Atoi"),
		),
	)
	return
}
