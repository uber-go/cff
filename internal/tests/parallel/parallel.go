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
	send := func(ctx context.Context) error {
		if err := ctx.Err(); err != nil {
			return err
		}
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
