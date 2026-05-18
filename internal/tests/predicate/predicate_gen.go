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

		_18_3 := context.Background()

		_19_15 := &s

		_21_4 := func() string {
			f()
			return "foo"
		}

		_25_18 := func() bool { return pred }
		ctx := _18_3
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/predicate/predicate.go",
				Line:   17,
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

		// go.uber.org/cff/internal/tests/predicate/predicate.go:25:4
		var p0 bool
		var p0PanicRecover interface{}
		var p0PanicStacktrace []byte
		_ = p0PanicStacktrace // possibly unused.
		pred1 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) (bool, error)
			job *cff.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (result bool, err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p0PanicRecover = recovered
					p0PanicStacktrace = debug.Stack()
				}
			}()
			p0 = _25_18()
			return p0, nil
		}

		pred1.job = sched.EnqueuePredicate(ctx, cff.PredicateJob{
			Run: pred1.run,
		})

		// go.uber.org/cff/internal/tests/predicate/predicate.go:21:4
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
				var stacktrace []byte
				if recovered != nil {
					stacktrace = debug.Stack()
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

			v1 = _21_4()

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

		*(_19_15) = v1 // string

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
}

// SimpleWithContextTask is a task flow which checks that context can be passed into Task w/out
// errors.
func SimpleWithContextTask() error {
	var s string
	return func() (err error) {

		_35_3 := context.Background()

		_36_15 := &s

		_37_14 := int64(2)

		_39_4 := func(ctx context.Context) string {
			return "foo"
		}

		_43_5 := func(int64) bool {
			return false
		}
		ctx := _35_3
		var v2 int64 = _37_14
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/predicate/predicate.go",
				Line:   34,
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

		// go.uber.org/cff/internal/tests/predicate/predicate.go:42:4
		var p0 bool
		var p0PanicRecover interface{}
		var p0PanicStacktrace []byte
		_ = p0PanicStacktrace // possibly unused.
		pred1 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) (bool, error)
			job *cff.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (result bool, err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p0PanicRecover = recovered
					p0PanicStacktrace = debug.Stack()
				}
			}()
			p0 = _43_5(v2)
			return p0, nil
		}

		pred1.job = sched.EnqueuePredicate(ctx, cff.PredicateJob{
			Run: pred1.run,
		})

		// go.uber.org/cff/internal/tests/predicate/predicate.go:39:4
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
				var stacktrace []byte
				if recovered != nil {
					stacktrace = debug.Stack()
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

			v1 = _39_4(ctx)

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

		*(_36_15) = v1 // string

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
}

// SimpleWithContextPredicate is a task flow which checks that context can be passed into
// Predicate but not Task.
func SimpleWithContextPredicate() error {
	var s string
	return func() (err error) {

		_55_3 := context.Background()

		_56_15 := &s

		_57_14 := int64(2)

		_59_4 := func() string {
			return "foo"
		}

		_63_5 := func(context.Context, int64) bool {
			return false
		}
		ctx := _55_3
		var v2 int64 = _57_14
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/predicate/predicate.go",
				Line:   54,
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

		// go.uber.org/cff/internal/tests/predicate/predicate.go:62:4
		var p0 bool
		var p0PanicRecover interface{}
		var p0PanicStacktrace []byte
		_ = p0PanicStacktrace // possibly unused.
		pred1 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) (bool, error)
			job *cff.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (result bool, err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p0PanicRecover = recovered
					p0PanicStacktrace = debug.Stack()
				}
			}()
			p0 = _63_5(ctx, v2)
			return p0, nil
		}

		pred1.job = sched.EnqueuePredicate(ctx, cff.PredicateJob{
			Run: pred1.run,
		})

		// go.uber.org/cff/internal/tests/predicate/predicate.go:59:4
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
				var stacktrace []byte
				if recovered != nil {
					stacktrace = debug.Stack()
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

			v1 = _59_4()

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

		*(_56_15) = v1 // string

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
}

// SimpleWithContextTaskAndPredicate is a task flow which checks that context can be passed into
// Predicate and Task.
func SimpleWithContextTaskAndPredicate() error {
	var s string
	return func() (err error) {

		_75_3 := context.Background()

		_76_15 := &s

		_77_14 := int64(2)

		_79_4 := func(ctx context.Context) string {
			return "foo"
		}

		_83_5 := func(context.Context, int64) bool {
			return false
		}
		ctx := _75_3
		var v2 int64 = _77_14
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/predicate/predicate.go",
				Line:   74,
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

		// go.uber.org/cff/internal/tests/predicate/predicate.go:82:4
		var p0 bool
		var p0PanicRecover interface{}
		var p0PanicStacktrace []byte
		_ = p0PanicStacktrace // possibly unused.
		pred1 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) (bool, error)
			job *cff.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (result bool, err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p0PanicRecover = recovered
					p0PanicStacktrace = debug.Stack()
				}
			}()
			p0 = _83_5(ctx, v2)
			return p0, nil
		}

		pred1.job = sched.EnqueuePredicate(ctx, cff.PredicateJob{
			Run: pred1.run,
		})

		// go.uber.org/cff/internal/tests/predicate/predicate.go:79:4
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
				var stacktrace []byte
				if recovered != nil {
					stacktrace = debug.Stack()
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

			v1 = _79_4(ctx)

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

		*(_76_15) = v1 // string

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

		_99_3 := context.Background()

		_100_14 := int(42)

		_101_15 := &out

		_103_4 := func(int) t1 { return t1{} }

		_105_4 := func() t2 { return t2{} }

		_107_4 := func(t2) t3 { return t3{} }

		_109_5 := func(int, t1) bool {
			return true
		}
		ctx := _99_3
		var v3 int = _100_14
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/predicate/predicate.go",
				Line:   98,
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

		// go.uber.org/cff/internal/tests/predicate/predicate.go:103:4
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
						Stacktrace: debug.Stack(),
					}
				}
			}()

			defer task4.ran.Store(true)

			v4 = _103_4(v3)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task4.job = sched.Enqueue(ctx, cff.Job{
			Run: task4.run,
		})
		tasks = append(tasks, task4)

		// go.uber.org/cff/internal/tests/predicate/predicate.go:105:4
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
						Stacktrace: debug.Stack(),
					}
				}
			}()

			defer task5.ran.Store(true)

			v5 = _105_4()

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task5.job = sched.Enqueue(ctx, cff.Job{
			Run: task5.run,
		})
		tasks = append(tasks, task5)

		// go.uber.org/cff/internal/tests/predicate/predicate.go:108:4
		var p0 bool
		var p0PanicRecover interface{}
		var p0PanicStacktrace []byte
		_ = p0PanicStacktrace // possibly unused.
		pred1 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) (bool, error)
			job *cff.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (result bool, err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p0PanicRecover = recovered
					p0PanicStacktrace = debug.Stack()
				}
			}()
			p0 = _109_5(v3, v4)
			return p0, nil
		}

		pred1.job = sched.EnqueuePredicate(ctx, cff.PredicateJob{
			Run: pred1.run,
			Dependencies: []*cff.ScheduledJob{
				task4.job,
			},
		})

		// go.uber.org/cff/internal/tests/predicate/predicate.go:107:4
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
				var stacktrace []byte
				if recovered != nil {
					stacktrace = debug.Stack()
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

			v6 = _107_4(v5)

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

		*(_101_15) = v6 // go.uber.org/cff/internal/tests/predicate.t3

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

		_123_3 := context.Background()

		_124_15 := &s

		_124_19 := &b

		_126_4 := func() string {
			return "foo"
		}

		_129_18 := func() bool { return true }

		_132_4 := func() bool {
			return true
		}

		_135_18 := func() bool { return false }
		ctx := _123_3
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/predicate/predicate.go",
				Line:   122,
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

		// go.uber.org/cff/internal/tests/predicate/predicate.go:129:4
		var p0 bool
		var p0PanicRecover interface{}
		var p0PanicStacktrace []byte
		_ = p0PanicStacktrace // possibly unused.
		pred1 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) (bool, error)
			job *cff.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (result bool, err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p0PanicRecover = recovered
					p0PanicStacktrace = debug.Stack()
				}
			}()
			p0 = _129_18()
			return p0, nil
		}

		pred1.job = sched.EnqueuePredicate(ctx, cff.PredicateJob{
			Run: pred1.run,
		})

		// go.uber.org/cff/internal/tests/predicate/predicate.go:126:4
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
				var stacktrace []byte
				if recovered != nil {
					stacktrace = debug.Stack()
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

			v1 = _126_4()

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

		// go.uber.org/cff/internal/tests/predicate/predicate.go:135:4
		var p1 bool
		var p1PanicRecover interface{}
		var p1PanicStacktrace []byte
		_ = p1PanicStacktrace // possibly unused.
		pred2 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) (bool, error)
			job *cff.ScheduledJob
		})
		pred2.run = func(ctx context.Context) (result bool, err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p1PanicRecover = recovered
					p1PanicStacktrace = debug.Stack()
				}
			}()
			p1 = _135_18()
			return p1, nil
		}

		pred2.job = sched.EnqueuePredicate(ctx, cff.PredicateJob{
			Run: pred2.run,
		})

		// go.uber.org/cff/internal/tests/predicate/predicate.go:132:4
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
				var stacktrace []byte
				if recovered != nil {
					stacktrace = debug.Stack()
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

			v7 = _132_4()

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

		*(_124_15) = v1 // string
		*(_124_19) = v7 // bool

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
}

// Panicked is a task flow that contains a task predicate that panics.
func Panicked() error {
	var s string
	return func() (err error) {

		_144_3 := context.Background()

		_145_15 := &s

		_147_4 := func(ctx context.Context) string {
			return "foo"
		}

		_151_5 := func() bool {
			panic("sad times")
			return true
		}
		ctx := _144_3
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/predicate/predicate.go",
				Line:   143,
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

		// go.uber.org/cff/internal/tests/predicate/predicate.go:150:4
		var p0 bool
		var p0PanicRecover interface{}
		var p0PanicStacktrace []byte
		_ = p0PanicStacktrace // possibly unused.
		pred1 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) (bool, error)
			job *cff.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (result bool, err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p0PanicRecover = recovered
					p0PanicStacktrace = debug.Stack()
				}
			}()
			p0 = _151_5()
			return p0, nil
		}

		pred1.job = sched.EnqueuePredicate(ctx, cff.PredicateJob{
			Run: pred1.run,
		})

		// go.uber.org/cff/internal/tests/predicate/predicate.go:147:4
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
				var stacktrace []byte
				if recovered != nil {
					stacktrace = debug.Stack()
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

			v1 = _147_4(ctx)

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

		*(_145_15) = v1 // string

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
}

// PanickedWithFallback is a flow that runs a panicing task predicate with a
// fallback.
func PanickedWithFallback() (string, error) {
	var s string
	err := func() (err error) {

		_165_3 := context.Background()

		_166_15 := &s

		_168_4 := func(ctx context.Context) (string, error) {
			return "foo", nil
		}

		_172_5 := func() bool {
			panic("sad times")
			return true
		}

		_177_21 := "predicate-fallback"
		ctx := _165_3
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/predicate/predicate.go",
				Line:   164,
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

		// go.uber.org/cff/internal/tests/predicate/predicate.go:171:4
		var p0 bool
		var p0PanicRecover interface{}
		var p0PanicStacktrace []byte
		_ = p0PanicStacktrace // possibly unused.
		pred1 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) (bool, error)
			job *cff.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (result bool, err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p0PanicRecover = recovered
					p0PanicStacktrace = debug.Stack()
				}
			}()
			p0 = _172_5()
			return p0, nil
		}

		pred1.job = sched.EnqueuePredicate(ctx, cff.PredicateJob{
			Run: pred1.run,
		})

		// go.uber.org/cff/internal/tests/predicate/predicate.go:168:4
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
					v1, err = _177_21, nil
				}
			}()

			if !p0 {
				return nil
			}

			defer task10.ran.Store(true)

			v1, err = _168_4(ctx)

			if err != nil {
				taskEmitter.TaskErrorRecovered(ctx, err)
				v1, err = _177_21, nil
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

		*(_166_15) = v1 // string

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
	return s, err
}

// BlockingInputs builds a flow that probes whether cff.Predicate
// short-circuits the scheduler's wait on the predicated task's input
// dependencies.
//
// Shape: a slow source feeds the predicated task; a fast source feeds
// the predicate. The predicate returns false (so the task body is
// skipped via FallbackWith). A downstream consumer measures the
// elapsed wall-clock time from flow start.
//
// If predicates short-circuit input deps:
//
//	elapsed ≈ 0 (slow path is not on the critical path)
//
// If predicates only short-circuit the task body:
//
//	elapsed ≈ slowDelay (scheduler still waits on slowOut)
func BlockingInputs(slowDelay time.Duration) (time.Duration, error) {
	type slowOut struct{}
	type fastOut struct{}
	type skippedOut struct{}

	var elapsed time.Duration
	start := time.Now()
	err := func() (err error) {

		_204_3 := context.Background()

		_205_15 := &elapsed

		_206_12 := func() slowOut {
			time.Sleep(slowDelay)
			return slowOut{}
		}

		_210_12 := func() fastOut { return fastOut{} }

		_212_4 := func(slowOut) (skippedOut, error) { return skippedOut{}, nil }

		_213_18 := func(fastOut) bool { return false }

		_214_21 := skippedOut{}

		_216_12 := func(skippedOut) time.Duration {
			return time.Since(start)
		}
		ctx := _204_3
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/predicate/predicate.go",
				Line:   203,
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

		// go.uber.org/cff/internal/tests/predicate/predicate.go:206:12
		var (
			v8 slowOut
		)
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
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: debug.Stack(),
					}
				}
			}()

			defer task11.ran.Store(true)

			v8 = _206_12()

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task11.job = sched.Enqueue(ctx, cff.Job{
			Run: task11.run,
		})
		tasks = append(tasks, task11)

		// go.uber.org/cff/internal/tests/predicate/predicate.go:210:12
		var (
			v9 fastOut
		)
		task12 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task12.emitter = cff.NopTaskEmitter()
		task12.run = func(ctx context.Context) (err error) {
			taskEmitter := task12.emitter
			startTime := time.Now()
			defer func() {
				if task12.ran.Load() {
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

			defer task12.ran.Store(true)

			v9 = _210_12()

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task12.job = sched.Enqueue(ctx, cff.Job{
			Run: task12.run,
		})
		tasks = append(tasks, task12)

		// go.uber.org/cff/internal/tests/predicate/predicate.go:213:4
		var p0 bool
		var p0PanicRecover interface{}
		var p0PanicStacktrace []byte
		_ = p0PanicStacktrace // possibly unused.
		pred1 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) (bool, error)
			job *cff.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (result bool, err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p0PanicRecover = recovered
					p0PanicStacktrace = debug.Stack()
				}
			}()
			p0 = _213_18(v9)
			return p0, nil
		}

		pred1.job = sched.EnqueuePredicate(ctx, cff.PredicateJob{
			Run: pred1.run,
			Dependencies: []*cff.ScheduledJob{
				task12.job,
			},
		})

		// go.uber.org/cff/internal/tests/predicate/predicate.go:212:4
		var (
			v10 skippedOut
		)
		task13 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task13.emitter = cff.NopTaskEmitter()
		task13.run = func(ctx context.Context) (err error) {
			taskEmitter := task13.emitter
			startTime := time.Now()
			defer func() {
				if task13.ran.Load() {
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
					v10, err = _214_21, nil
				}
			}()

			if !p0 {
				return nil
			}

			defer task13.ran.Store(true)

			v10, err = _212_4(v8)

			if err != nil {
				taskEmitter.TaskErrorRecovered(ctx, err)
				v10, err = _214_21, nil
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task13.job = sched.Enqueue(ctx, cff.Job{
			Run: task13.run,
			Dependencies: []*cff.ScheduledJob{
				task11.job,
				pred1.job,
			},
		})
		tasks = append(tasks, task13)

		// go.uber.org/cff/internal/tests/predicate/predicate.go:216:12
		var (
			v11 time.Duration
		)
		task14 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task14.emitter = cff.NopTaskEmitter()
		task14.run = func(ctx context.Context) (err error) {
			taskEmitter := task14.emitter
			startTime := time.Now()
			defer func() {
				if task14.ran.Load() {
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

			defer task14.ran.Store(true)

			v11 = _216_12(v10)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task14.job = sched.Enqueue(ctx, cff.Job{
			Run: task14.run,
			Dependencies: []*cff.ScheduledJob{
				task13.job,
			},
		})
		tasks = append(tasks, task14)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_205_15) = v11 // time.Duration

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
	return elapsed, err
}
