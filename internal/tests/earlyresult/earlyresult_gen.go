//go:build !cff
// +build !cff

package earlyresult

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/cff"
)

type foo struct{}
type bar struct{}
type baz struct{}
type qux struct{}

// EarlyResult makes sure ordering for an early cff.Results doesn't cause compiler to error.
func EarlyResult(ctx context.Context) error {
	request := int(2)
	var out *bar
	var out2 *foo
	return func() (err error) {

		_23_3 := ctx

		_24_14 := request

		_25_15 := &out

		_25_21 := &out2

		_27_4 := func(*foo) *bar {
			return &bar{}
		}

		_31_4 := func(*foo) *baz {
			return &baz{}
		}

		_35_4 := func(*bar, *baz) *qux {
			return &qux{}
		}

		_39_4 := func(int) *foo {
			return &foo{}
		}

		_43_4 := func(*qux) error {
			return nil
		}
		ctx := _23_3
		var v1 int = _24_14
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/earlyresult/earlyresult.go",
				Line:   22,
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

		// go.uber.org/cff/internal/tests/earlyresult/earlyresult.go:39:4
		var (
			v2 *foo
		)
		task3 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task3.emitter = cff.NopTaskEmitter()
		task3.run = func(ctx context.Context) (err error) {
			taskEmitter := task3.emitter
			startTime := time.Now()
			defer func() {
				if task3.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			defer task3.ran.Store(true)

			v2 = _39_4(v1)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task3.job = sched.Enqueue(ctx, cff.Job{
			Run: task3.run,
		})
		tasks = append(tasks, task3)

		// go.uber.org/cff/internal/tests/earlyresult/earlyresult.go:27:4
		var (
			v3 *bar
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
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			defer task0.ran.Store(true)

			v3 = _27_4(v2)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task0.job = sched.Enqueue(ctx, cff.Job{
			Run: task0.run,
			Dependencies: []*cff.ScheduledJob{
				task3.job,
			},
		})
		tasks = append(tasks, task0)

		// go.uber.org/cff/internal/tests/earlyresult/earlyresult.go:31:4
		var (
			v4 *baz
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
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			defer task1.ran.Store(true)

			v4 = _31_4(v2)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task1.job = sched.Enqueue(ctx, cff.Job{
			Run: task1.run,
			Dependencies: []*cff.ScheduledJob{
				task3.job,
			},
		})
		tasks = append(tasks, task1)

		// go.uber.org/cff/internal/tests/earlyresult/earlyresult.go:35:4
		var (
			v5 *qux
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
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			defer task2.ran.Store(true)

			v5 = _35_4(v3, v4)

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

		// go.uber.org/cff/internal/tests/earlyresult/earlyresult.go:43:4
		task4 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task4.emitter = cff.NopTaskEmitter()
		task4.run = func(ctx context.Context) (err error) {
			taskEmitter := task4.emitter
			startTime := time.Now()
			defer func() {
				if task4.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			defer task4.ran.Store(true)

			err = _43_4(v5)

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task4.job = sched.Enqueue(ctx, cff.Job{
			Run: task4.run,
			Dependencies: []*cff.ScheduledJob{
				task2.job,
			},
		})
		tasks = append(tasks, task4)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_25_15) = v3 // *go.uber.org/cff/internal/tests/earlyresult.bar
		*(_25_21) = v2 // *go.uber.org/cff/internal/tests/earlyresult.foo

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
}

// ConsumesResult makes sure that we can have an early cff.Results and run post-processing tasks.
func ConsumesResult() error {
	// t1 -> genService.Status_GetStatus_Args
	type t1 struct{}
	// t2 -> statusValidator.Request
	type t2 struct{}
	// t3 -> genService.StatusResponse
	type t3 struct{}
	// t4 -> statusValidator.Response
	type t4 struct{}
	// t5 -> node.NodeContext
	type t5 struct{}
	// t6 -> statusPostProcessor.Request
	type t6 struct{}
	// t7 -> statusPostProcessor.Response
	type t7 struct{}

	var v1 *t3
	var request *t1

	return func() (err error) {

		_71_18 := context.Background()

		_72_15 := &v1

		_73_14 := request

		_76_4 := func(*t1) *t2 { return &t2{} }

		_80_4 := func(*t4) *t5 { return &t5{} }

		_84_4 := func(*t2) (*t4, error) { return nil, nil }

		_89_4 := func(*t5) (*t3, error) { return nil, nil }

		_94_4 := func(*t3) *t6 { return &t6{} }

		_98_4 := func(*t6) (*t7, error) { return nil, nil }

		_101_4 := func(*t7) error {
			return nil
		}
		ctx := _71_18
		var v6 *t1 = _73_14
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/earlyresult/earlyresult.go",
				Line:   71,
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

		// go.uber.org/cff/internal/tests/earlyresult/earlyresult.go:76:4
		var (
			v7 *t2
		)
		task5 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task5.emitter = cff.NopTaskEmitter()
		task5.run = func(ctx context.Context) (err error) {
			taskEmitter := task5.emitter
			startTime := time.Now()
			defer func() {
				if task5.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			defer task5.ran.Store(true)

			v7 = _76_4(v6)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task5.job = sched.Enqueue(ctx, cff.Job{
			Run: task5.run,
		})
		tasks = append(tasks, task5)

		// go.uber.org/cff/internal/tests/earlyresult/earlyresult.go:84:4
		var (
			v8 *t4
		)
		task7 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task7.emitter = cff.NopTaskEmitter()
		task7.run = func(ctx context.Context) (err error) {
			taskEmitter := task7.emitter
			startTime := time.Now()
			defer func() {
				if task7.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			defer task7.ran.Store(true)

			v8, err = _84_4(v7)

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task7.job = sched.Enqueue(ctx, cff.Job{
			Run: task7.run,
			Dependencies: []*cff.ScheduledJob{
				task5.job,
			},
		})
		tasks = append(tasks, task7)

		// go.uber.org/cff/internal/tests/earlyresult/earlyresult.go:80:4
		var (
			v9 *t5
		)
		task6 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task6.emitter = cff.NopTaskEmitter()
		task6.run = func(ctx context.Context) (err error) {
			taskEmitter := task6.emitter
			startTime := time.Now()
			defer func() {
				if task6.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			defer task6.ran.Store(true)

			v9 = _80_4(v8)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task6.job = sched.Enqueue(ctx, cff.Job{
			Run: task6.run,
			Dependencies: []*cff.ScheduledJob{
				task7.job,
			},
		})
		tasks = append(tasks, task6)

		// go.uber.org/cff/internal/tests/earlyresult/earlyresult.go:89:4
		var (
			v10 *t3
		)
		task8 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task8.emitter = cff.NopTaskEmitter()
		task8.run = func(ctx context.Context) (err error) {
			taskEmitter := task8.emitter
			startTime := time.Now()
			defer func() {
				if task8.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			defer task8.ran.Store(true)

			v10, err = _89_4(v9)

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task8.job = sched.Enqueue(ctx, cff.Job{
			Run: task8.run,
			Dependencies: []*cff.ScheduledJob{
				task6.job,
			},
		})
		tasks = append(tasks, task8)

		// go.uber.org/cff/internal/tests/earlyresult/earlyresult.go:94:4
		var (
			v11 *t6
		)
		task9 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task9.emitter = cff.NopTaskEmitter()
		task9.run = func(ctx context.Context) (err error) {
			taskEmitter := task9.emitter
			startTime := time.Now()
			defer func() {
				if task9.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			defer task9.ran.Store(true)

			v11 = _94_4(v10)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task9.job = sched.Enqueue(ctx, cff.Job{
			Run: task9.run,
			Dependencies: []*cff.ScheduledJob{
				task8.job,
			},
		})
		tasks = append(tasks, task9)

		// go.uber.org/cff/internal/tests/earlyresult/earlyresult.go:98:4
		var (
			v12 *t7
		)
		task10 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task10.emitter = cff.NopTaskEmitter()
		task10.run = func(ctx context.Context) (err error) {
			taskEmitter := task10.emitter
			startTime := time.Now()
			defer func() {
				if task10.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			defer task10.ran.Store(true)

			v12, err = _98_4(v11)

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task10.job = sched.Enqueue(ctx, cff.Job{
			Run: task10.run,
			Dependencies: []*cff.ScheduledJob{
				task9.job,
			},
		})
		tasks = append(tasks, task10)

		// go.uber.org/cff/internal/tests/earlyresult/earlyresult.go:101:4
		task11 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task11.emitter = cff.NopTaskEmitter()
		task11.run = func(ctx context.Context) (err error) {
			taskEmitter := task11.emitter
			startTime := time.Now()
			defer func() {
				if task11.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			defer task11.ran.Store(true)

			err = _101_4(v12)

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task11.job = sched.Enqueue(ctx, cff.Job{
			Run: task11.run,
			Dependencies: []*cff.ScheduledJob{
				task10.job,
			},
		})
		tasks = append(tasks, task11)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_72_15) = v10 // *go.uber.org/cff/internal/tests/earlyresult.t3

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
}
