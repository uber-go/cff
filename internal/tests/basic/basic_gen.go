//go:build !cff
// +build !cff

package basic

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"go.uber.org/cff"
)

type foo struct{ i int }
type bar struct{ i int64 }

// SimpleFlow is a very simple flow with some inputs and outputs.
func SimpleFlow() (string, error) {
	var message string
	err := func() (err error) {

		_20_18 := context.Background()

		_21_14 := 1

		_22_15 := &message

		_24_4 := func(i int) int64 {
			return int64(i)
		}

		_29_4 := func(i int) (*foo, error) {
			return &foo{i}, nil
		}

		_33_4 := func(i int64) (*bar, error) {
			return &bar{i}, nil
		}

		_37_4 := func(*foo, *bar) (string, error) {
			return "hello world", nil
		}
		ctx := _20_18
		var v1 int = _21_14
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/basic/basic.go",
				Line:   20,
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

		// go.uber.org/cff/internal/tests/basic/basic.go:24:4
		var (
			v2 int64
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

			v2 = _24_4(v1)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task0.job = sched.Enqueue(ctx, cff.Job{
			Run: task0.run,
		})
		tasks = append(tasks, task0)

		// go.uber.org/cff/internal/tests/basic/basic.go:29:4
		var (
			v3 *foo
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

			v3, err = _29_4(v1)

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

		// go.uber.org/cff/internal/tests/basic/basic.go:33:4
		var (
			v4 *bar
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

			v4, err = _33_4(v2)

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task2.job = sched.Enqueue(ctx, cff.Job{
			Run: task2.run,
			Dependencies: []*cff.ScheduledJob{
				task0.job,
			},
		})
		tasks = append(tasks, task2)

		// go.uber.org/cff/internal/tests/basic/basic.go:37:4
		var (
			v5 string
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

			v5, err = _37_4(v3, v4)

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task3.job = sched.Enqueue(ctx, cff.Job{
			Run: task3.run,
			Dependencies: []*cff.ScheduledJob{
				task1.job,
				task2.job,
			},
		})
		tasks = append(tasks, task3)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_22_15) = v5 // string

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
	return message, err
}

// NoParamsFlow is a flow that does not accept any parameters.
func NoParamsFlow(ctx context.Context) (io.Reader, error) {
	var r io.Reader
	err := func() (err error) {

		_47_18 := ctx

		_48_15 := &r

		_50_4 := func() *bytes.Buffer {
			return bytes.NewBufferString("hello world")
		}

		_54_4 := func(b *bytes.Buffer) io.Reader { return b }
		ctx := _47_18
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/basic/basic.go",
				Line:   47,
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

		// go.uber.org/cff/internal/tests/basic/basic.go:50:4
		var (
			v6 *bytes.Buffer
		)
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

			v6 = _50_4()

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task4.job = sched.Enqueue(ctx, cff.Job{
			Run: task4.run,
		})
		tasks = append(tasks, task4)

		// go.uber.org/cff/internal/tests/basic/basic.go:54:4
		var (
			v7 io.Reader
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

			v7 = _54_4(v6)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task5.job = sched.Enqueue(ctx, cff.Job{
			Run: task5.run,
			Dependencies: []*cff.ScheduledJob{
				task4.job,
			},
		})
		tasks = append(tasks, task5)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_48_15) = v7 // io.Reader

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
	return r, err
}

// SerialFailableFlow runs the provided function in-order using a flow.
func SerialFailableFlow(ctx context.Context, f1, f2 func() error) error {
	type t1 struct{}
	type t2 struct{}
	type t3 struct{}

	// We use the types to make f2 depend on f1.
	var out t3
	return func() (err error) {

		_69_3 := ctx

		_70_15 := &out

		_72_4 := func() (t1, error) {
			return t1{}, f1()
		}

		_76_4 := func(t1) (t2, error) {
			return t2{}, f2()
		}

		_80_4 := func(t2) t3 {
			return t3{}
		}
		ctx := _69_3
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/basic/basic.go",
				Line:   68,
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

		// go.uber.org/cff/internal/tests/basic/basic.go:72:4
		var (
			v8 t1
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

			v8, err = _72_4()

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task6.job = sched.Enqueue(ctx, cff.Job{
			Run: task6.run,
		})
		tasks = append(tasks, task6)

		// go.uber.org/cff/internal/tests/basic/basic.go:76:4
		var (
			v9 t2
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

			v9, err = _76_4(v8)

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
				task6.job,
			},
		})
		tasks = append(tasks, task7)

		// go.uber.org/cff/internal/tests/basic/basic.go:80:4
		var (
			v10 t3
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

			v10 = _80_4(v9)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task8.job = sched.Enqueue(ctx, cff.Job{
			Run: task8.run,
			Dependencies: []*cff.ScheduledJob{
				task7.job,
			},
		})
		tasks = append(tasks, task8)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_70_15) = v10 // go.uber.org/cff/internal/tests/basic.t3

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
}

// ProduceMultiple has a task which produces multiple values.
func ProduceMultiple() error {
	type t1 struct{}
	type t2 struct{}
	type t3 struct{}
	type t4 struct{}

	//   t1
	//   /\
	//  v   v
	// t2   t3
	//   \ /
	//    v
	//   t4

	var out t4
	return func() (err error) {

		_104_3 := context.Background()

		_105_14 := t1{}

		_106_15 := &out

		_108_4 := func(t1) (t2, t3) {
			return t2{}, t3{}
		}

		_112_4 := func(t2, t3) t4 {
			return t4{}
		}
		ctx := _104_3
		var v11 t1 = _105_14
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/basic/basic.go",
				Line:   103,
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

		// go.uber.org/cff/internal/tests/basic/basic.go:108:4
		var (
			v12 t2
			v13 t3
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

			v12, v13 = _108_4(v11)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task9.job = sched.Enqueue(ctx, cff.Job{
			Run: task9.run,
		})
		tasks = append(tasks, task9)

		// go.uber.org/cff/internal/tests/basic/basic.go:112:4
		var (
			v14 t4
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

			v14 = _112_4(v12, v13)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task10.job = sched.Enqueue(ctx, cff.Job{
			Run: task10.run,
			Dependencies: []*cff.ScheduledJob{
				task9.job,
				task9.job,
			},
		})
		tasks = append(tasks, task10)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_106_15) = v14 // go.uber.org/cff/internal/tests/basic.t4

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
}
