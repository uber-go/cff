package parallel

import (
	"context"
	"fmt"
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
		cff.Slice(
			func(ctx context.Context, idx int, s string) error {
				_ = fmt.Sprintf("%d and %q", idx, s)
				_, _ = ctx.Deadline()
				return nil
			},
			[]string{"some", "thing"},
		),

		cff.Map(
			func(ctx context.Context, key string, value string) error {
				_ = fmt.Sprintf("%q : %q", key, value)
				_, _ = ctx.Deadline()
				return nil
			},
			map[string]string{"key": "value"},
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

	someSlice := []string{"some", "slice"}
	sliceFunc := func(_ int, _ string) {}

	someMap := map[string]string{"key": "value"}
	mapFunc := func(_ string, _ string) {}

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
		cff.Slice(
			sliceFunc,
			someSlice,
		),
		cff.Map(
			mapFunc,
			someMap,
		),
	)
	if err != nil {
		return err
	}
	return nil
}
