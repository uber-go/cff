//go:build !cff
// +build !cff

package cffintest

import (
	"context"
	"runtime/debug"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/cff"
)

func TestIsOdd(t *testing.T) {
	err := func() (err error) {

		_17_3 := context.Background()

		_18_19 := 4

		_20_4 := func() error {
			assert.Equal(t, true, isOdd(1))
			return nil
		}

		_24_4 := func() error {
			assert.Equal(t, false, isOdd(2))
			return nil
		}

		_28_4 := func() error {
			assert.Equal(t, true, isOdd(3))
			return nil
		}

		_32_4 := func() error {
			assert.Equal(t, false, isOdd(4))
			return nil
		}
		ctx := _17_3
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/cffintest/cffintest_test.go",
				Line:   16,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Concurrency: _18_19, Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/cffintest/cffintest_test.go:20:4
		task0 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task0.emitter = cff.NopTaskEmitter()
		task0.fn = func(ctx context.Context) (err error) {
			taskEmitter := task0.emitter
			startTime := time.Now()
			defer func() {
				if task0.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered == nil {
					return
				}
				taskEmitter.TaskPanic(ctx, recovered)
				err = &cff.PanicError{
					Value:      recovered,
					Stacktrace: string(debug.Stack()),
				}
			}()

			defer task0.ran.Store(true)

			err = _20_4()

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return
			}
			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task0.fn,
		})
		tasks = append(tasks, task0)

		// go.uber.org/cff/internal/tests/cffintest/cffintest_test.go:24:4
		task1 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task1.emitter = cff.NopTaskEmitter()
		task1.fn = func(ctx context.Context) (err error) {
			taskEmitter := task1.emitter
			startTime := time.Now()
			defer func() {
				if task1.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered == nil {
					return
				}
				taskEmitter.TaskPanic(ctx, recovered)
				err = &cff.PanicError{
					Value:      recovered,
					Stacktrace: string(debug.Stack()),
				}
			}()

			defer task1.ran.Store(true)

			err = _24_4()

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return
			}
			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task1.fn,
		})
		tasks = append(tasks, task1)

		// go.uber.org/cff/internal/tests/cffintest/cffintest_test.go:28:4
		task2 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task2.emitter = cff.NopTaskEmitter()
		task2.fn = func(ctx context.Context) (err error) {
			taskEmitter := task2.emitter
			startTime := time.Now()
			defer func() {
				if task2.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered == nil {
					return
				}
				taskEmitter.TaskPanic(ctx, recovered)
				err = &cff.PanicError{
					Value:      recovered,
					Stacktrace: string(debug.Stack()),
				}
			}()

			defer task2.ran.Store(true)

			err = _28_4()

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return
			}
			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task2.fn,
		})
		tasks = append(tasks, task2)

		// go.uber.org/cff/internal/tests/cffintest/cffintest_test.go:32:4
		task3 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task3.emitter = cff.NopTaskEmitter()
		task3.fn = func(ctx context.Context) (err error) {
			taskEmitter := task3.emitter
			startTime := time.Now()
			defer func() {
				if task3.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered == nil {
					return
				}
				taskEmitter.TaskPanic(ctx, recovered)
				err = &cff.PanicError{
					Value:      recovered,
					Stacktrace: string(debug.Stack()),
				}
			}()

			defer task3.ran.Store(true)

			err = _32_4()

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return
			}
			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task3.fn,
		})
		tasks = append(tasks, task3)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line cffintest_test.go:36*/
	}()
	require.NoError(t, err)
}
