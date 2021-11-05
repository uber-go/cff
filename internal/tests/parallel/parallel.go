package parallel

import (
	"context"
	"errors"
	"sync"

	"go.uber.org/cff"
)

// Simple runs tasks in parallel to populate the provided map.
func Simple(m *sync.Map) error {
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
	)
}

// SimpleWithError runs tasks in parallel to populate the provided channel.
func SimpleWithError(c chan<- string) error {
	return cff.Parallel(
		context.Background(),
		cff.Concurrency(2),
		cff.Tasks(
			func() error {
				return errors.New("sad times")
			},
			func(_ context.Context) {
				c <- "work"
			},
		),
	)
}

// SimpleWithPanic runs a parallel task that panics.
func SimpleWithPanic() error {
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
