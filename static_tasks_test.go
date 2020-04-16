package cff

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
)

// Convenience type to build lists of functions.
type Schedule = [][]func(context.Context) error

// This tests the basics of the behaviors of RunStaticTasks with a single
// gororutine. This test avoids going into behaviors that require multiple
// goroutines.
func TestRunStaticTasks_SanityChecks(t *testing.T) {
	ctx := context.Background()

	t.Run("empty", func(t *testing.T) {
		require.NoError(t, RunStaticTasks(ctx, nil))
	})

	t.Run("non-empty schedule, empty list", func(t *testing.T) {
		require.NoError(t, RunStaticTasks(ctx, Schedule{{}}))
	})

	t.Run("pre-cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(ctx)
		cancel()

		err := RunStaticTasks(ctx, Schedule{
			{
				func(context.Context) error {
					t.Error("should not be called")
					return nil
				},
			},
		})

		require.Error(t, err)
		assert.True(t, errors.Is(err, context.Canceled),
			"should be context cancelled error")
	})

	t.Run("single task is called", func(t *testing.T) {
		var called atomic.Bool
		err := RunStaticTasks(ctx, Schedule{
			{
				func(context.Context) error {
					assert.False(t, called.Swap(true),
						"task called twice")
					return nil
				},
			},
		})
		require.NoError(t, err)
	})

	t.Run("single task returns error", func(t *testing.T) {
		var errSadness = errors.New("great sadness")

		err := RunStaticTasks(ctx, Schedule{
			{
				func(context.Context) error {
					return errSadness
				},
			},
		})
		require.Error(t, err)
		assert.True(t, errors.Is(err, errSadness), "must be errSadness")
	})

	t.Run("single task cancels context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(ctx)
		cancel()

		var called atomic.Bool
		err := RunStaticTasks(ctx, Schedule{
			{
				func(context.Context) error {
					assert.False(t, called.Swap(true),
						"task called twice")
					cancel()
					return nil
				},
				func(context.Context) error {
					t.Error("should not be called")
					return nil
				},
			},
		})

		require.Error(t, err)
		assert.True(t, errors.Is(err, context.Canceled),
			"should be context cancelled error")
	})
}

func TestRunStaticTasks(t *testing.T) {
	ctx := context.Background()

	t.Run("multiple tasks are run together", func(t *testing.T) {
		const N = 100

		// To tests that all tasks in a group are executed together,
		// we'll count the number of ongoing tasks and a barrier which
		// blocks them from exiting before all the others have started
		// up.
		//
		// When the running count hits N, we'll release the barrier so
		// that they can all exit safely.
		//
		// Note that we can't juse use the task index (i) because
		// there's no guarantee that the last task in the list is also
		// the last task to actually get scheduled.
		var running atomic.Int64
		barrier := make(chan struct{})

		tasks := make([]func(context.Context) error, N)
		for i := range tasks {
			tasks[i] = func(ctx context.Context) error {
				if running.Inc() == N {
					// Last task closes the barrier so
					// that all tasks can exit.
					close(barrier)
				}
				defer running.Dec()

				select {
				case <-barrier:
					return nil
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}

		// Don't wait more than 100 milliseconds for all tasks to have
		// started up.
		ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()

		err := RunStaticTasks(ctx, Schedule{tasks})
		require.NoError(t, err)
	})

	t.Run("first error in group is captured", func(t *testing.T) {
		const (
			N      = 100
			failOn = 25 // all tasks at this position or more will fail
		)

		// Similar to the previous test, we'll count the task position
		// with an atomic and start failing tasks when the count hits
		// failOn.
		//
		// The barrier applies to failing tasks only, and it ensures
		// that the first task to hit the failure condition is also
		// the first to return. Otherwise, another task could return
		// early and be the first to fail from CFF's point of view.
		var (
			running atomic.Int64
			wantErr error
			barrier = make(chan struct{})
		)
		tasks := make([]func(context.Context) error, N)
		for i := range tasks {
			tasks[i] = func(context.Context) error {
				pos := running.Inc()
				if pos < failOn {
					return nil
				}

				// The first failing task unblocks the
				// barrier. All others wait on it.
				err := fmt.Errorf("error %d", pos)
				if pos == failOn {
					defer close(barrier)
					wantErr = err
				} else {
					<-barrier
				}

				return err
			}
		}

		ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()

		err := RunStaticTasks(ctx, Schedule{tasks})
		require.Error(t, err)
		assert.Truef(t, errors.Is(err, wantErr), "must be error %d", failOn)
	})

	t.Run("single task is run after many", func(t *testing.T) {
		const N = 50

		var ran atomic.Int64

		tasks := make([]func(context.Context) error, N)
		for i := range tasks {
			tasks[i] = func(context.Context) error {
				ran.Inc()
				return nil
			}
		}

		single := func(context.Context) error {
			assert.Truef(t, ran.CAS(N, N+1),
				"must run all %d tasks before running single", N)
			return nil
		}

		err := RunStaticTasks(ctx, Schedule{
			tasks,
			{single},
		})
		require.NoError(t, err)

		assert.Equal(t, int64(N+1), ran.Load())
	})
}
