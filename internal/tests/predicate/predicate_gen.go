//go:build !cff
// +build !cff

package predicate

import (
	"context"
	"runtime/debug"
	"time"

	"go.uber.org/cff"
)

// Simple runs the provided function in a task flow if the provided boolean
// is true.
func Simple(f func(), pred bool) error {
	var s string
	return func() (err error) {

		_17_3 := context.Background()

		_18_15 := &s

		_20_4 := func() string {
			f()
			return "foo"
		}

		_24_18 := func() bool { return pred }
		ctx := _17_3
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/predicate/predicate.go",
				Line:   16,
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

		// go.uber.org/cff/internal/tests/predicate/predicate.go:24:4
		var p0 bool
		var p0PanicRecover interface{}
		var p0PanicStacktrace string
		_ = p0PanicStacktrace // possibly unused.
		pred1 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) error
			job *cff.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p0PanicRecover = recovered
					p0PanicStacktrace = string(debug.Stack())
				}
			}()
			p0 = _24_18()
			return nil
		}

		pred1.job = sched.Enqueue(ctx, cff.Job{
			Run: pred1.run,
		})

		// go.uber.org/cff/internal/tests/predicate/predicate.go:20:4
		var (
			v1 string
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
				var stacktrace string
				if recovered != nil {
					stacktrace = string(debug.Stack())
				}
				if recovered == nil && p0PanicRecover != nil {
					recovered = p0PanicRecover
					stacktrace = p0PanicStacktrace
				}
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: stacktrace,
					}
				}
			}()

			if !p0 {
				return nil
			}

			defer task0.ran.Store(true)

			v1 = _20_4()

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task0.job = sched.Enqueue(ctx, cff.Job{
			Run: task0.run,
			Dependencies: []*cff.ScheduledJob{
				pred1.job,
			},
		})
		tasks = append(tasks, task0)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_18_15) = v1 // string

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
}

// SimpleWithContextTask is a task flow which checks that context can be passed into Task w/out
// errors.
func SimpleWithContextTask() error {
	var s string
	return func() (err error) {

		_34_3 := context.Background()

		_35_15 := &s

		_36_14 := int64(2)

		_38_4 := func(ctx context.Context) string {
			return "foo"
		}

		_42_5 := func(int64) bool {
			return false
		}
		ctx := _34_3
		var v2 int64 = _36_14
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/predicate/predicate.go",
				Line:   33,
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

		// go.uber.org/cff/internal/tests/predicate/predicate.go:41:4
		var p0 bool
		var p0PanicRecover interface{}
		var p0PanicStacktrace string
		_ = p0PanicStacktrace // possibly unused.
		pred1 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) error
			job *cff.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p0PanicRecover = recovered
					p0PanicStacktrace = string(debug.Stack())
				}
			}()
			p0 = _42_5(v2)
			return nil
		}

		pred1.job = sched.Enqueue(ctx, cff.Job{
			Run: pred1.run,
		})

		// go.uber.org/cff/internal/tests/predicate/predicate.go:38:4
		var (
			v1 string
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
				var stacktrace string
				if recovered != nil {
					stacktrace = string(debug.Stack())
				}
				if recovered == nil && p0PanicRecover != nil {
					recovered = p0PanicRecover
					stacktrace = p0PanicStacktrace
				}
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: stacktrace,
					}
				}
			}()

			if !p0 {
				return nil
			}

			defer task1.ran.Store(true)

			v1 = _38_4(ctx)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task1.job = sched.Enqueue(ctx, cff.Job{
			Run: task1.run,
			Dependencies: []*cff.ScheduledJob{
				pred1.job,
			},
		})
		tasks = append(tasks, task1)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_35_15) = v1 // string

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
}

// SimpleWithContextPredicate is a task flow which checks that context can be passed into
// Predicate but not Task.
func SimpleWithContextPredicate() error {
	var s string
	return func() (err error) {

		_54_3 := context.Background()

		_55_15 := &s

		_56_14 := int64(2)

		_58_4 := func() string {
			return "foo"
		}

		_62_5 := func(context.Context, int64) bool {
			return false
		}
		ctx := _54_3
		var v2 int64 = _56_14
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/predicate/predicate.go",
				Line:   53,
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

		// go.uber.org/cff/internal/tests/predicate/predicate.go:61:4
		var p0 bool
		var p0PanicRecover interface{}
		var p0PanicStacktrace string
		_ = p0PanicStacktrace // possibly unused.
		pred1 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) error
			job *cff.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p0PanicRecover = recovered
					p0PanicStacktrace = string(debug.Stack())
				}
			}()
			p0 = _62_5(ctx, v2)
			return nil
		}

		pred1.job = sched.Enqueue(ctx, cff.Job{
			Run: pred1.run,
		})

		// go.uber.org/cff/internal/tests/predicate/predicate.go:58:4
		var (
			v1 string
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
				var stacktrace string
				if recovered != nil {
					stacktrace = string(debug.Stack())
				}
				if recovered == nil && p0PanicRecover != nil {
					recovered = p0PanicRecover
					stacktrace = p0PanicStacktrace
				}
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: stacktrace,
					}
				}
			}()

			if !p0 {
				return nil
			}

			defer task2.ran.Store(true)

			v1 = _58_4()

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task2.job = sched.Enqueue(ctx, cff.Job{
			Run: task2.run,
			Dependencies: []*cff.ScheduledJob{
				pred1.job,
			},
		})
		tasks = append(tasks, task2)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_55_15) = v1 // string

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
}

// SimpleWithContextTaskAndPredicate is a task flow which checks that context can be passed into
// Predicate and Task.
func SimpleWithContextTaskAndPredicate() error {
	var s string
	return func() (err error) {

		_74_3 := context.Background()

		_75_15 := &s

		_76_14 := int64(2)

		_78_4 := func(ctx context.Context) string {
			return "foo"
		}

		_82_5 := func(context.Context, int64) bool {
			return false
		}
		ctx := _74_3
		var v2 int64 = _76_14
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/predicate/predicate.go",
				Line:   73,
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

		// go.uber.org/cff/internal/tests/predicate/predicate.go:81:4
		var p0 bool
		var p0PanicRecover interface{}
		var p0PanicStacktrace string
		_ = p0PanicStacktrace // possibly unused.
		pred1 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) error
			job *cff.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p0PanicRecover = recovered
					p0PanicStacktrace = string(debug.Stack())
				}
			}()
			p0 = _82_5(ctx, v2)
			return nil
		}

		pred1.job = sched.Enqueue(ctx, cff.Job{
			Run: pred1.run,
		})

		// go.uber.org/cff/internal/tests/predicate/predicate.go:78:4
		var (
			v1 string
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
				var stacktrace string
				if recovered != nil {
					stacktrace = string(debug.Stack())
				}
				if recovered == nil && p0PanicRecover != nil {
					recovered = p0PanicRecover
					stacktrace = p0PanicStacktrace
				}
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: stacktrace,
					}
				}
			}()

			if !p0 {
				return nil
			}

			defer task3.ran.Store(true)

			v1 = _78_4(ctx)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task3.job = sched.Enqueue(ctx, cff.Job{
			Run: task3.run,
			Dependencies: []*cff.ScheduledJob{
				pred1.job,
			},
		})
		tasks = append(tasks, task3)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_75_15) = v1 // string

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
}

// ExtraDependencies is a task flow where the predicate has more dependencies
// than the task.
func ExtraDependencies() error {
	type t1 struct{}
	type t2 struct{}
	type t3 struct{}

	var out t3
	return func() (err error) {

		_98_3 := context.Background()

		_99_14 := int(42)

		_100_15 := &out

		_102_4 := func(int) t1 { return t1{} }

		_104_4 := func() t2 { return t2{} }

		_106_4 := func(t2) t3 { return t3{} }

		_108_5 := func(int, t1) bool {
			return true
		}
		ctx := _98_3
		var v3 int = _99_14
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/predicate/predicate.go",
				Line:   97,
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

		// go.uber.org/cff/internal/tests/predicate/predicate.go:102:4
		var (
			v4 t1
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
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: string(debug.Stack()),
					}
				}
			}()

			defer task4.ran.Store(true)

			v4 = _102_4(v3)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task4.job = sched.Enqueue(ctx, cff.Job{
			Run: task4.run,
		})
		tasks = append(tasks, task4)

		// go.uber.org/cff/internal/tests/predicate/predicate.go:104:4
		var (
			v5 t2
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
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: string(debug.Stack()),
					}
				}
			}()

			defer task5.ran.Store(true)

			v5 = _104_4()

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task5.job = sched.Enqueue(ctx, cff.Job{
			Run: task5.run,
		})
		tasks = append(tasks, task5)

		// go.uber.org/cff/internal/tests/predicate/predicate.go:107:4
		var p0 bool
		var p0PanicRecover interface{}
		var p0PanicStacktrace string
		_ = p0PanicStacktrace // possibly unused.
		pred1 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) error
			job *cff.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p0PanicRecover = recovered
					p0PanicStacktrace = string(debug.Stack())
				}
			}()
			p0 = _108_5(v3, v4)
			return nil
		}

		pred1.job = sched.Enqueue(ctx, cff.Job{
			Run: pred1.run,
			Dependencies: []*cff.ScheduledJob{
				task4.job,
			},
		})

		// go.uber.org/cff/internal/tests/predicate/predicate.go:106:4
		var (
			v6 t3
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
				var stacktrace string
				if recovered != nil {
					stacktrace = string(debug.Stack())
				}
				if recovered == nil && p0PanicRecover != nil {
					recovered = p0PanicRecover
					stacktrace = p0PanicStacktrace
				}
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: stacktrace,
					}
				}
			}()

			if !p0 {
				return nil
			}

			defer task6.ran.Store(true)

			v6 = _106_4(v5)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task6.job = sched.Enqueue(ctx, cff.Job{
			Run: task6.run,
			Dependencies: []*cff.ScheduledJob{
				task5.job,
				pred1.job,
			},
		})
		tasks = append(tasks, task6)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_100_15) = v6 // go.uber.org/cff/internal/tests/predicate.t3

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
}

// MultiplePredicates is a task flow which checks that the outputs of multiple
// predicates can be distinguished.
func MultiplePredicates() error {
	var s string
	var b bool
	return func() (err error) {

		_122_3 := context.Background()

		_123_15 := &s

		_123_19 := &b

		_125_4 := func() string {
			return "foo"
		}

		_128_18 := func() bool { return true }

		_131_4 := func() bool {
			return true
		}

		_134_18 := func() bool { return false }
		ctx := _122_3
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/predicate/predicate.go",
				Line:   121,
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

		// go.uber.org/cff/internal/tests/predicate/predicate.go:128:4
		var p0 bool
		var p0PanicRecover interface{}
		var p0PanicStacktrace string
		_ = p0PanicStacktrace // possibly unused.
		pred1 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) error
			job *cff.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p0PanicRecover = recovered
					p0PanicStacktrace = string(debug.Stack())
				}
			}()
			p0 = _128_18()
			return nil
		}

		pred1.job = sched.Enqueue(ctx, cff.Job{
			Run: pred1.run,
		})

		// go.uber.org/cff/internal/tests/predicate/predicate.go:125:4
		var (
			v1 string
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
				var stacktrace string
				if recovered != nil {
					stacktrace = string(debug.Stack())
				}
				if recovered == nil && p0PanicRecover != nil {
					recovered = p0PanicRecover
					stacktrace = p0PanicStacktrace
				}
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: stacktrace,
					}
				}
			}()

			if !p0 {
				return nil
			}

			defer task7.ran.Store(true)

			v1 = _125_4()

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task7.job = sched.Enqueue(ctx, cff.Job{
			Run: task7.run,
			Dependencies: []*cff.ScheduledJob{
				pred1.job,
			},
		})
		tasks = append(tasks, task7)

		// go.uber.org/cff/internal/tests/predicate/predicate.go:134:4
		var p1 bool
		var p1PanicRecover interface{}
		var p1PanicStacktrace string
		_ = p1PanicStacktrace // possibly unused.
		pred2 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) error
			job *cff.ScheduledJob
		})
		pred2.run = func(ctx context.Context) (err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p1PanicRecover = recovered
					p1PanicStacktrace = string(debug.Stack())
				}
			}()
			p1 = _134_18()
			return nil
		}

		pred2.job = sched.Enqueue(ctx, cff.Job{
			Run: pred2.run,
		})

		// go.uber.org/cff/internal/tests/predicate/predicate.go:131:4
		var (
			v7 bool
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
				var stacktrace string
				if recovered != nil {
					stacktrace = string(debug.Stack())
				}
				if recovered == nil && p1PanicRecover != nil {
					recovered = p1PanicRecover
					stacktrace = p1PanicStacktrace
				}
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: stacktrace,
					}
				}
			}()

			if !p1 {
				return nil
			}

			defer task8.ran.Store(true)

			v7 = _131_4()

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task8.job = sched.Enqueue(ctx, cff.Job{
			Run: task8.run,
			Dependencies: []*cff.ScheduledJob{
				pred2.job,
			},
		})
		tasks = append(tasks, task8)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_123_15) = v1 // string
		*(_123_19) = v7 // bool

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
}

// Panicked is a task flow that contains a task predicate that panics.
func Panicked() error {
	var s string
	return func() (err error) {

		_143_3 := context.Background()

		_144_15 := &s

		_146_4 := func(ctx context.Context) string {
			return "foo"
		}

		_150_5 := func() bool {
			panic("sad times")
			return true
		}
		ctx := _143_3
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/predicate/predicate.go",
				Line:   142,
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

		// go.uber.org/cff/internal/tests/predicate/predicate.go:149:4
		var p0 bool
		var p0PanicRecover interface{}
		var p0PanicStacktrace string
		_ = p0PanicStacktrace // possibly unused.
		pred1 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) error
			job *cff.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p0PanicRecover = recovered
					p0PanicStacktrace = string(debug.Stack())
				}
			}()
			p0 = _150_5()
			return nil
		}

		pred1.job = sched.Enqueue(ctx, cff.Job{
			Run: pred1.run,
		})

		// go.uber.org/cff/internal/tests/predicate/predicate.go:146:4
		var (
			v1 string
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
				var stacktrace string
				if recovered != nil {
					stacktrace = string(debug.Stack())
				}
				if recovered == nil && p0PanicRecover != nil {
					recovered = p0PanicRecover
					stacktrace = p0PanicStacktrace
				}
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: stacktrace,
					}
				}
			}()

			if !p0 {
				return nil
			}

			defer task9.ran.Store(true)

			v1 = _146_4(ctx)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task9.job = sched.Enqueue(ctx, cff.Job{
			Run: task9.run,
			Dependencies: []*cff.ScheduledJob{
				pred1.job,
			},
		})
		tasks = append(tasks, task9)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_144_15) = v1 // string

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
}

// PanickedWithFallback is a flow that runs a panicing task predicate with a
// fallback.
func PanickedWithFallback() (string, error) {
	var s string
	err := func() (err error) {

		_164_3 := context.Background()

		_165_15 := &s

		_167_4 := func(ctx context.Context) (string, error) {
			return "foo", nil
		}

		_171_5 := func() bool {
			panic("sad times")
			return true
		}

		_176_21 := "predicate-fallback"
		ctx := _164_3
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/predicate/predicate.go",
				Line:   163,
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

		// go.uber.org/cff/internal/tests/predicate/predicate.go:170:4
		var p0 bool
		var p0PanicRecover interface{}
		var p0PanicStacktrace string
		_ = p0PanicStacktrace // possibly unused.
		pred1 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) error
			job *cff.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p0PanicRecover = recovered
					p0PanicStacktrace = string(debug.Stack())
				}
			}()
			p0 = _171_5()
			return nil
		}

		pred1.job = sched.Enqueue(ctx, cff.Job{
			Run: pred1.run,
		})

		// go.uber.org/cff/internal/tests/predicate/predicate.go:167:4
		var (
			v1 string
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

				if recovered == nil && p0PanicRecover != nil {
					recovered = p0PanicRecover

				}
				if recovered != nil {
					taskEmitter.TaskPanicRecovered(ctx, recovered)
					v1, err = _176_21, nil
				}
			}()

			if !p0 {
				return nil
			}

			defer task10.ran.Store(true)

			v1, err = _167_4(ctx)

			if err != nil {
				taskEmitter.TaskErrorRecovered(ctx, err)
				v1, err = _176_21, nil
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task10.job = sched.Enqueue(ctx, cff.Job{
			Run: task10.run,
			Dependencies: []*cff.ScheduledJob{
				pred1.job,
			},
		})
		tasks = append(tasks, task10)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_165_15) = v1 // string

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
	return s, err
}
