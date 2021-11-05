package parallel

import (
	"context"
	"sync"

	"go.uber.org/cff"
)

// ExampleParallel provides a flow that deploys multiple calls to cff.Parallel
// This flow is compiled to test the CFF compiler's internal state.
func ExampleParallel(m *sync.Map, c chan<- string) error {
	sendFn := func() {
		c <- "send"
	}

	err := cff.Parallel(
		context.Background(),
		cff.Concurrency(2),
		cff.Tasks(
			func() {
				m.Store("foo", "done")
			},
			sendFn,
		),
		cff.Task(
			func(ctx context.Context) error {
				return ctx.Err()
			},
		),
	)
	if err != nil {
		return err
	}

	sendFnCtxErr := func(ctx context.Context) error {
		_, _ = ctx.Deadline()
		c <- "send"
		return nil
	}

	err = cff.Parallel(
		context.Background(),
		cff.Concurrency(2),
		cff.Tasks(
			func(_ context.Context) {
				m.Store("bar", "complete")
			},
			sendFnCtxErr,
		),
		cff.Task(
			func() {
				m.Store("bar", "finished")
			},
		),
	)
	if err != nil {
		return err
	}
	return nil
}
