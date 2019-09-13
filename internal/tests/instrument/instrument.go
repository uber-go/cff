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

// T3630161 reproduces T3630161 by executing a flow that runs a task that failed, recovers, and then runs another task.
func (h *H) T3630161(ctx context.Context) {
	var s string
	_ = cff.Flow(ctx,
		cff.Results(&s),
		cff.Metrics(h.Scope),
		cff.Logger(h.Logger),
		cff.InstrumentFlow("T3630161"),

		cff.Task(
			func() (string, error) {
				return "", errors.New("always errors")
			},
			cff.Instrument("Err"),
			cff.FallbackWith("fallback value"),
		),

		cff.Task(
			func(s string) error {
				return nil
			},
			cff.Instrument("End"),
			cff.Invoke(true),
		),
	)
	return
}

// T3795761 reproduces T3795761 where a task that returns no error should only emit skipped metric if it was not run
func (h *H) T3795761(ctx context.Context, shouldRun bool, shouldError bool) string {
	var s string
	_ = cff.Flow(ctx,
		cff.Results(&s),
		cff.Metrics(h.Scope),
		cff.Logger(h.Logger),
		cff.InstrumentFlow("T3795761"),

		cff.Task(
			func() int {
				return 0
			},
			cff.Instrument("ProvidesInt"),
		),

		cff.Task(
			func(s int) (string, error) {
				if shouldError {
					return "", errors.New("err")
				}

				return "ok", nil
			},
			cff.Predicate(func() bool { return shouldRun }),
			cff.Instrument("NeedsInt"),
		),
	)
	return s
}
