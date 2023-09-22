//go:build !cff
// +build !cff

package insidegeneric

import (
	"context"
	"runtime/debug"
	"time"

	"go.uber.org/cff"
)

// Producer is a function that produces a value of the given type.
type Producer[T any] func(context.Context) (T, error)

// JoinTwo combines two producers producing different values using the provided
// function.
func JoinTwo[A, B, C any](
	pa Producer[A],
	pb Producer[B],
	fn func(A, B) C,
) (C, error) {
	var c C
	err := func() (err error) {

		_23_18 := context.Background()

		_24_15 := &c

		_25_12 := pa

		_26_12 := pb

		_27_12 := fn
		ctx := _23_18
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/insidegeneric/producer.go",
				Line:   23,
				Column: 9,
			}
			flowEmitter = cff.NopFlowEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      flowInfo.Name,
				Directive: cff.FlowDirective,
				File:      flowInfo.File,
				Line:      flowInfo.Line,
				Column:    flowInfo.Column,
			}

			// possibly unused
			_ = flowInfo
		)

		startTime := time.Now()
		defer func() { flowEmitter.FlowDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/insidegeneric/producer.go:25:12
		var (
			v1 A
		)
		task0 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task0.emitter = cff.NopTaskEmitter()
		task0.run = func(ctx context.Context) (err error) {
			taskEmitter := task0.emitter
			startTime := time.Now()
			defer func() {
				if task0.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: debug.Stack(),
					}
				}
			}()

			defer task0.ran.Store(true)

			v1, err = _25_12(ctx)

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task0.job = sched.Enqueue(ctx, cff.Job{
			Run: task0.run,
		})
		tasks = append(tasks, task0)

		// go.uber.org/cff/internal/tests/insidegeneric/producer.go:26:12
		var (
			v2 B
		)
		task1 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task1.emitter = cff.NopTaskEmitter()
		task1.run = func(ctx context.Context) (err error) {
			taskEmitter := task1.emitter
			startTime := time.Now()
			defer func() {
				if task1.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: debug.Stack(),
					}
				}
			}()

			defer task1.ran.Store(true)

			v2, err = _26_12(ctx)

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task1.job = sched.Enqueue(ctx, cff.Job{
			Run: task1.run,
		})
		tasks = append(tasks, task1)

		// go.uber.org/cff/internal/tests/insidegeneric/producer.go:27:12
		var (
			v3 C
		)
		task2 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task2.emitter = cff.NopTaskEmitter()
		task2.run = func(ctx context.Context) (err error) {
			taskEmitter := task2.emitter
			startTime := time.Now()
			defer func() {
				if task2.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: debug.Stack(),
					}
				}
			}()

			defer task2.ran.Store(true)

			v3 = _27_12(v1, v2)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task2.job = sched.Enqueue(ctx, cff.Job{
			Run: task2.run,
			Dependencies: []*cff.ScheduledJob{
				task0.job,
				task1.job,
			},
		})
		tasks = append(tasks, task2)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_24_15) = v3 // C

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
	return c, err
}

// JoinMany runs the given producers and returns a slice of their results
// in-order.
func JoinMany[T any](producers ...Producer[T]) ([]T, error) {
	results := make([]T, len(producers))
	err := func() (err error) {

		_36_22 := context.Background()

		_38_4 := func(ctx context.Context, idx int, fn Producer[T]) error {
			v, err := fn(ctx)
			results[idx] = v
			return err
		}

		_43_4 := producers
		ctx := _36_22
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/insidegeneric/producer.go",
				Line:   36,
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
				Emitter: schedEmitter,
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

		// go.uber.org/cff/internal/tests/insidegeneric/producer.go:37:3
		sliceTask3Slice := _43_4
		for idx, val := range sliceTask3Slice {
			idx := idx
			val := val
			sliceTask3 := new(struct {
				emitter cff.TaskEmitter
				fn      func(context.Context) error
				ran     cff.AtomicBool
			})
			sliceTask3.fn = func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = &cff.PanicError{
							Value:      recovered,
							Stacktrace: debug.Stack(),
						}
					}
				}()
				err = _38_4(ctx, idx, val)
				return
			}
			sched.Enqueue(ctx, cff.Job{
				Run: sliceTask3.fn,
			})
		}

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line producer.go:44*/
	}()
	return results, err
}
