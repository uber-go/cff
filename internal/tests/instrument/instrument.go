// +build cff

// Package instrument verifies that default and custom Emitter
// implementations trigger on events and emit specification in the CFF2 ERD.
// DefaultEmitter tests default emitter.
// These tests will be removed in the future as an implementation detail.
// CustomEmitter tests mocks for custom emitter.
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
	h := &DefaultEmitter{
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

// DefaultEmitter is used by other tests.
type DefaultEmitter struct {
	Scope  tally.Scope
	Logger *zap.Logger
}

// Run executes a flow to test instrumentation.
func (h *DefaultEmitter) Run(ctx context.Context, req string, fields ...zap.Field) (res uint8, err error) {
	err = cff.Flow(ctx,
		cff.Params(req),
		cff.Results(&res),
		cff.Metrics(h.Scope),
		cff.Logger(h.Logger),
		cff.WithLogFields(fields...),
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

// ExplicitListOfFields is a flow with an explicit list of log fields.
func (h *DefaultEmitter) ExplicitListOfFields(ctx context.Context, req string) (res int, err error) {
	err = cff.Flow(ctx,
		cff.Params(req),
		cff.Results(&res),
		cff.InstrumentFlow("ExplicitListOfFields"),
		cff.Metrics(h.Scope),
		cff.Logger(h.Logger),
		cff.WithLogFields(zap.String("foo", "bar"), zap.Int("baz", 42)),
		cff.Task(
			strconv.Atoi,
			cff.Instrument("Atoi"),
		),
	)
	return
}

// Do executes a flow to test instrumentation.
func (h *DefaultEmitter) Do(ctx context.Context, req string) (res int, err error) {
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
func (h *DefaultEmitter) Work(ctx context.Context, req string) (res int, err error) {
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
func (h *DefaultEmitter) T3630161(ctx context.Context) {
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
func (h *DefaultEmitter) T3795761(ctx context.Context, shouldRun bool, shouldError bool) string {
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

// These tests replicate the ones written for instrumentation to verify that
// custom Emitter will trigger similarly to default implementation.

// CustomEmitter is used by other tests.
type CustomEmitter struct {
	Scope   tally.Scope
	Logger  *zap.Logger
	Emitter cff.Emitter
}

// Run executes a flow to test instrumentation.
func (h *CustomEmitter) Run(ctx context.Context, req string) (res uint8, err error) {
	err = cff.Flow(ctx,
		cff.Params(req),
		cff.Results(&res),
		cff.Logger(h.Logger),
		cff.InstrumentFlow("AtoiRun"),
		cff.WithEmitter(h.Emitter),
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
func (h *CustomEmitter) Do(ctx context.Context, req string) (res int, err error) {
	err = cff.Flow(ctx,
		cff.Params(req),
		cff.Results(&res),
		cff.InstrumentFlow("AtoiDo"),
		cff.WithEmitter(h.Emitter),
		cff.Logger(h.Logger),
		cff.Task(
			strconv.Atoi,
			cff.Instrument("Atoi"),
		),
	)
	return
}

// Work executes a flow to test instrumentation.
func (h *CustomEmitter) Work(ctx context.Context, req string) (res int, err error) {
	err = cff.Flow(ctx,
		cff.Params(req),
		cff.Results(&res),
		cff.WithEmitter(h.Emitter),
		cff.Logger(h.Logger),
		cff.Task(
			strconv.Atoi,
			cff.Instrument("Atoi"),
		),
	)
	return
}

// T3630161 reproduces T3630161 by executing a flow that runs a task that failed,
// recovers, and then runs another task.
func (h *CustomEmitter) T3630161(ctx context.Context) {
	var s string
	_ = cff.Flow(ctx,
		cff.Results(&s),
		cff.WithEmitter(h.Emitter),
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

// T3795761 reproduces T3795761 where a task that returns no error should only
// emit skipped metric if it was not run.
func (h *CustomEmitter) T3795761(ctx context.Context, shouldRun bool,
	shouldError bool) string {
	var s string
	_ = cff.Flow(ctx,
		cff.Results(&s),
		cff.WithEmitter(h.Emitter),
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

// FlowAlwaysPanics is a flow that tests Metrics Emitter
func (h *CustomEmitter) FlowAlwaysPanics(ctx context.Context) error {
	return cff.Flow(ctx,
		cff.Logger(h.Logger),
		cff.WithEmitter(h.Emitter),
		cff.Task(func() {
			panic("always")
		},
			cff.Invoke(true),
			cff.Instrument("Atoi"),
		),
	)
}
