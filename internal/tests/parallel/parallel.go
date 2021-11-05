package parallel

import (
	"context"
	"errors"
	"sync"

	"go.uber.org/cff"
)

// TasksAndTask runs cff.Tasks and cff.Task in parallel to populate the
// provided map.
func TasksAndTask(m *sync.Map) error {
	return cff.Parallel(
		context.Background(),
		cff.Concurrency(2),
		cff.Tasks(
			func() {
				m.Store("foo", "bar")
			},
			func(_ context.Context) {
				m.Store("fiz", "buzz")
			},
		),
		cff.Task(
			func(_ context.Context) {
				m.Store("go", "lang")
			},
		),
	)
}

// TasksWithError runs a parallel cff.Tasks that errors.
func TasksWithError() error {
	return cff.Parallel(
		context.Background(),
		cff.Concurrency(2),
		cff.Tasks(
			func() error {
				return errors.New("sad times")
			},
		),
	)
}

// TasksWithPanic runs a parallel cff.Tasks that panics.
func TasksWithPanic() error {
	return cff.Parallel(
		context.Background(),
		cff.Concurrency(2),
		cff.Tasks(
			func() {
				panic("sad times")
			},
		),
	)
}

// MultipleTasks runs multiple cff.Tasks in parallel to populate the provided
// channel.
func MultipleTasks(c chan<- string) error {
	send := func(_ context.Context) error {
		c <- "send"
		return nil
	}
	return cff.Parallel(
		context.Background(),
		cff.Concurrency(2),
		cff.Tasks(
			func() {
				c <- "multiple"
			},
		),
		cff.Tasks(
			send,
		),
	)
}

// ContextErrorBefore runs a cff.Tasks function to test that the function is
// not run if the context errors before scheduler execution.
func ContextErrorBefore(ctx context.Context, src, target []int) error {
	return cff.Parallel(
		ctx,
		cff.Concurrency(2),
		cff.Tasks(
			func() {
				target[0] = src[0]
			},
		),
	)
}

// ContextErrorInFlight runs a cff.Tasks function to test that the function is
// not run if the context errors during scheduler execution.
func ContextErrorInFlight(ctx context.Context, cancel func(), src, target []int) error {
	blocker := make(chan struct{})
	return cff.Parallel(
		ctx,
		cff.Concurrency(2),
		cff.Tasks(
			// Busy both workers during context cancellation so that the
			// third function is processed after cancellation.
			func() {
				cancel()
				close(blocker)
			},
			func() {
				<-blocker
			},
			func() {
				target[0] = src[0]
			},
		),
	)
}

// TaskWithError runs a parallel cff.Task that errors.
func TaskWithError() error {
	return cff.Parallel(
		context.Background(),
		cff.Concurrency(2),
		cff.Task(
			func() error {
				return errors.New("sad times")
			},
		),
	)
}

// TaskWithPanic runs a parallel cff.task that panics.
func TaskWithPanic() error {
	return cff.Parallel(
		context.Background(),
		cff.Concurrency(2),
		cff.Task(
			func() {
				panic("sad times")
			},
		),
	)
}

// MultipleTask runs multiple cff.Task in parallel to populate the provided
// slice.
func MultipleTask(src, target []int) error {
	send := func(ctx context.Context) error {
		_, _ = ctx.Deadline()
		target[1] = src[1]
		return nil
	}
	return cff.Parallel(
		context.Background(),
		cff.Concurrency(2),
		cff.Task(
			func() {
				target[0] = src[0]
			},
		),
		cff.Task(
			send,
		),
	)
}
