//go:build !cff
// +build !cff

package shadowedvar

import (
	"context"
	"fmt"
	"time"

	cff2 "go.uber.org/cff"
)

// CtxConflict introduces a variable conflict with ctx to demonstrate that
// cff does not shadow variables.
func CtxConflict(ctx string) (string, error) {
	var result string
	err := func() (err error) {

		_17_3 := context.Background()

		_18_16 := &result

		_19_13 := func() (string, error) {
			var hello string
			hello = ctx
			return hello, nil
		}
		ctx := _17_3
		emitter := cff2.NopEmitter()

		var (
			flowInfo = &cff2.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/shadowedvar/shadowedvar.go",
				Line:   16,
				Column: 9,
			}
			flowEmitter = cff2.NopFlowEmitter()

			schedInfo = &cff2.SchedulerInfo{
				Name:      flowInfo.Name,
				Directive: cff2.FlowDirective,
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

		sched := cff2.NewScheduler(
			cff2.SchedulerParams{
				Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff2.TaskEmitter
			ran     cff2.AtomicBool
			run     func(context.Context) error
			job     *cff2.ScheduledJob
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/shadowedvar/shadowedvar.go:19:13
		var (
			v1 string
		)
		task0 := new(struct {
			emitter cff2.TaskEmitter
			ran     cff2.AtomicBool
			run     func(context.Context) error
			job     *cff2.ScheduledJob
		})
		task0.emitter = cff2.NopTaskEmitter()
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

			v1, err = _19_13()

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task0.job = sched.Enqueue(ctx, cff2.Job{
			Run: task0.run,
		})
		tasks = append(tasks, task0)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_18_16) = v1 // string

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()

	return result, err
}

// CtxConflictParallel introduces a variable conflict with ctx within cff.Parallel Task
// to demonstrate that cff does not shadow variables.
func CtxConflictParallel(ctx string) (string, string, error) {
	var result1 string
	var result2 string
	err := func() (err error) {

		_35_3 := context.Background()

		_36_13 := func() {
			result1 = ctx
		}

		_39_13 := func() {
			result2 = ctx
		}
		ctx := _35_3
		emitter := cff2.NopEmitter()

		var (
			parallelInfo = &cff2.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/shadowedvar/shadowedvar.go",
				Line:   34,
				Column: 9,
			}
			directiveInfo = &cff2.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff2.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff2.NopParallelEmitter()

			schedInfo = &cff2.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff2.ParallelDirective,
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

		sched := cff2.NewScheduler(
			cff2.SchedulerParams{
				Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff2.TaskEmitter
			fn      func(context.Context) error
			ran     cff2.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/shadowedvar/shadowedvar.go:36:13
		task1 := new(struct {
			emitter cff2.TaskEmitter
			fn      func(context.Context) error
			ran     cff2.AtomicBool
		})
		task1.emitter = cff2.NopTaskEmitter()
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
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task1.ran.Store(true)

			_36_13()

			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff2.Job{
			Run: task1.fn,
		})
		tasks = append(tasks, task1)

		// go.uber.org/cff/internal/tests/shadowedvar/shadowedvar.go:39:13
		task2 := new(struct {
			emitter cff2.TaskEmitter
			fn      func(context.Context) error
			ran     cff2.AtomicBool
		})
		task2.emitter = cff2.NopTaskEmitter()
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
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task2.ran.Store(true)

			_39_13()

			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff2.Job{
			Run: task2.fn,
		})
		tasks = append(tasks, task2)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line shadowedvar.go:41*/
	}()
	return result1, result2, err
}

// CtxConflictSlice introduces a variable conflict with ctx within cff.Slice function
// to demonstrate that cff does not shadow variables.
func CtxConflictSlice(ctx string, target []string) error {
	return func() (err error) {

		_50_3 := context.Background()

		_51_20 := 2

		_53_4 := func(idx int, val string) error {
			target[idx] = ctx + val
			return nil
		}

		_57_4 := target
		ctx := _50_3
		emitter := cff2.NopEmitter()

		var (
			parallelInfo = &cff2.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/shadowedvar/shadowedvar.go",
				Line:   49,
				Column: 9,
			}
			directiveInfo = &cff2.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff2.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff2.NopParallelEmitter()

			schedInfo = &cff2.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff2.ParallelDirective,
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

		sched := cff2.NewScheduler(
			cff2.SchedulerParams{
				Concurrency: _51_20, Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff2.TaskEmitter
			fn      func(context.Context) error
			ran     cff2.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/shadowedvar/shadowedvar.go:52:3
		sliceTask3Slice := _57_4
		for idx, val := range sliceTask3Slice {
			idx := idx
			val := val
			sliceTask3 := new(struct {
				emitter cff2.TaskEmitter
				fn      func(context.Context) error
				ran     cff2.AtomicBool
			})
			sliceTask3.fn = func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()
				err = _53_4(idx, val)
				return
			}
			sched.Enqueue(ctx, cff2.Job{
				Run: sliceTask3.fn,
			})
		}

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line shadowedvar.go:58*/
	}()
}

// CtxConflictMap introduces a variable conflict with ctx within cff.Map function
// to demonstrate that cff does not shadow variables.
func CtxConflictMap(ctx int, input map[int]int) ([]int, error) {
	slice := make([]int, len(input))
	err := func() (err error) {

		_67_3 := context.Background()

		_68_20 := 2

		_70_4 := func(key int, val int) {
			slice[key] = ctx + val
		}

		_73_4 := input
		ctx := _67_3
		emitter := cff2.NopEmitter()

		var (
			parallelInfo = &cff2.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/shadowedvar/shadowedvar.go",
				Line:   66,
				Column: 9,
			}
			directiveInfo = &cff2.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff2.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff2.NopParallelEmitter()

			schedInfo = &cff2.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff2.ParallelDirective,
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

		sched := cff2.NewScheduler(
			cff2.SchedulerParams{
				Concurrency: _68_20, Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff2.TaskEmitter
			fn      func(context.Context) error
			ran     cff2.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/shadowedvar/shadowedvar.go:69:3
		for key, val := range _73_4 {
			key := key
			val := val
			mapTask4 := new(struct {
				emitter cff2.TaskEmitter
				fn      func(context.Context) error
				ran     cff2.AtomicBool
			})
			mapTask4.fn = func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()

				_70_4(key, val)
				return
			}

			sched.Enqueue(ctx, cff2.Job{
				Run: mapTask4.fn,
			})
		}

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line shadowedvar.go:74*/
	}()
	return slice, err
}

// PredicateCtxConflict runs the provided function in a task flow if the
// provided boolean is true. This tests if the cff flow works even when the ctx
// variable is shadowed.
func PredicateCtxConflict(f func(), ctx bool) error {
	var s string
	return func() (err error) {

		_85_3 := context.Background()

		_86_16 := &s

		_88_4 := func() string {
			f()
			return "foo"
		}

		_92_19 := func() bool { return ctx }
		ctx := _85_3
		emitter := cff2.NopEmitter()

		var (
			flowInfo = &cff2.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/shadowedvar/shadowedvar.go",
				Line:   84,
				Column: 9,
			}
			flowEmitter = cff2.NopFlowEmitter()

			schedInfo = &cff2.SchedulerInfo{
				Name:      flowInfo.Name,
				Directive: cff2.FlowDirective,
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

		sched := cff2.NewScheduler(
			cff2.SchedulerParams{
				Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff2.TaskEmitter
			ran     cff2.AtomicBool
			run     func(context.Context) error
			job     *cff2.ScheduledJob
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/shadowedvar/shadowedvar.go:92:4
		var p0 bool
		var p0PanicRecover interface{}
		pred1 := new(struct {
			ran cff2.AtomicBool
			run func(context.Context) error
			job *cff2.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p0PanicRecover = recovered
				}
			}()
			p0 = _92_19()
			return nil
		}

		pred1.job = sched.Enqueue(ctx, cff2.Job{
			Run: pred1.run,
		})

		// go.uber.org/cff/internal/tests/shadowedvar/shadowedvar.go:88:4
		var (
			v1 string
		)
		task5 := new(struct {
			emitter cff2.TaskEmitter
			ran     cff2.AtomicBool
			run     func(context.Context) error
			job     *cff2.ScheduledJob
		})
		task5.emitter = cff2.NopTaskEmitter()
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
				if recovered == nil && p0PanicRecover != nil {
					recovered = p0PanicRecover
				}
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			if !p0 {
				return nil
			}

			defer task5.ran.Store(true)

			v1 = _88_4()

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task5.job = sched.Enqueue(ctx, cff2.Job{
			Run: task5.run,
			Dependencies: []*cff2.ScheduledJob{
				pred1.job,
			},
		})
		tasks = append(tasks, task5)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_86_16) = v1 // string

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
}
