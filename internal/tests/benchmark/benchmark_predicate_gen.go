//go:build !cff
// +build !cff

package benchmark

import (
	"context"
	"runtime/debug"
	"time"

	"go.uber.org/cff"
)

const (
	_workDuration = 100 * time.Millisecond
	_concurrency  = 2
)

// work is a task function that takes a pre-specifed amount of time.
func work() int {
	time.Sleep(_workDuration)
	return 0
}

// pred is a predicate function that takes a pre-specifed amount of time.
func pred() bool {
	time.Sleep(_workDuration)
	return true
}

// PredicateCombined runs a cff.Flow that exercises the function of a prior
// version of cff.Predicate that ran the predicate function within the task
// that the predicate is an option for (i.e the state of the world before
// https://code.uberinternal.com/D5495165).
func PredicateCombined() float64 {
	var res float64
	func() (err error) {

		_37_3 := context.Background()

		_38_19 := _concurrency

		_39_15 := &res

		_41_4 := func() int {
			return work()
		}

		_46_4 := func(num int) (f float64) {
			if !pred() {
				return
			}
			f = float64(work())
			return
		}
		ctx := _37_3
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/benchmark/benchmark_predicate.go",
				Line:   36,
				Column: 2,
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
				Concurrency: _38_19, Emitter: schedEmitter,
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

		// go.uber.org/cff/internal/tests/benchmark/benchmark_predicate.go:41:4
		var (
			v1 int
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

			v1 = _41_4()

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task0.job = sched.Enqueue(ctx, cff.Job{
			Run: task0.run,
		})
		tasks = append(tasks, task0)

		// go.uber.org/cff/internal/tests/benchmark/benchmark_predicate.go:46:4
		var (
			v2 float64
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

			v2 = _46_4(v1)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task1.job = sched.Enqueue(ctx, cff.Job{
			Run: task1.run,
			Dependencies: []*cff.ScheduledJob{
				task0.job,
			},
		})
		tasks = append(tasks, task1)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_39_15) = v2 // float64

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
	return res
}

// PredicateSplit runs a cff.Flow to exercise the function of the current
// predicate optimization which decouples the predicate function from the task
// it is an option for (i.e the state of the world after
// https://code.uberinternal.com/D5495165).
func PredicateSplit() float64 {
	var res float64
	func() (err error) {

		_65_3 := context.Background()

		_66_19 := _concurrency

		_67_15 := &res

		_69_4 := func() int {
			return work()
		}

		_74_4 := func(num int) float64 {
			return float64(work())
		}

		_78_5 := pred
		ctx := _65_3
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/benchmark/benchmark_predicate.go",
				Line:   64,
				Column: 2,
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
				Concurrency: _66_19, Emitter: schedEmitter,
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

		// go.uber.org/cff/internal/tests/benchmark/benchmark_predicate.go:69:4
		var (
			v1 int
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

			v1 = _69_4()

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task2.job = sched.Enqueue(ctx, cff.Job{
			Run: task2.run,
		})
		tasks = append(tasks, task2)

		// go.uber.org/cff/internal/tests/benchmark/benchmark_predicate.go:77:4
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
			p0 = _78_5()
			return p0, nil
		}

		pred1.job = sched.EnqueuePredicate(ctx, cff.PredicateJob{
			Run: pred1.run,
		})

		// go.uber.org/cff/internal/tests/benchmark/benchmark_predicate.go:74:4
		var (
			v2 float64
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

			v2 = _74_4(v1)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task3.job = sched.Enqueue(ctx, cff.Job{
			Run: task3.run,
			Dependencies: []*cff.ScheduledJob{
				task2.job,
				pred1.job,
			},
		})
		tasks = append(tasks, task3)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_67_15) = v2 // float64

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
	return res
}

// PredicateFalseWithSlowDep runs a cff.Flow where a predicate returns false
// while a sibling input task sleeps for slowDelay. Returns the elapsed
// wall-clock time from flow start to consumer dispatch, so the benchmark
// can report consumer_ns/op as a custom metric.
//
// Before scheduler-aware predicates, the consumer waited the full
// slowDelay before its skip gate fired. After, the consumer dispatches
// immediately on predicate=false, even though Wait() still drains the
// slow producer.
func PredicateFalseWithSlowDep(slowDelay time.Duration) time.Duration {
	type slowOut struct{}
	type fastOut struct{}
	type skippedOut struct{}

	var elapsed time.Duration
	start := time.Now()
	func() (err error) {

		_102_3 := context.Background()

		_103_19 := _concurrency

		_104_15 := &elapsed

		_105_12 := func() slowOut {
			time.Sleep(slowDelay)
			return slowOut{}
		}

		_109_12 := func() fastOut { return fastOut{} }

		_111_4 := func(slowOut) (skippedOut, error) { return skippedOut{}, nil }

		_112_18 := func(fastOut) bool { return false }

		_113_21 := skippedOut{}

		_115_12 := func(skippedOut) time.Duration {
			return time.Since(start)
		}
		ctx := _102_3
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/benchmark/benchmark_predicate.go",
				Line:   101,
				Column: 2,
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
				Concurrency: _103_19, Emitter: schedEmitter,
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

		// go.uber.org/cff/internal/tests/benchmark/benchmark_predicate.go:105:12
		var (
			v3 slowOut
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

			v3 = _105_12()

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task4.job = sched.Enqueue(ctx, cff.Job{
			Run: task4.run,
		})
		tasks = append(tasks, task4)

		// go.uber.org/cff/internal/tests/benchmark/benchmark_predicate.go:109:12
		var (
			v4 fastOut
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

			v4 = _109_12()

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task5.job = sched.Enqueue(ctx, cff.Job{
			Run: task5.run,
		})
		tasks = append(tasks, task5)

		// go.uber.org/cff/internal/tests/benchmark/benchmark_predicate.go:112:4
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
			p0 = _112_18(v4)
			return p0, nil
		}

		pred1.job = sched.EnqueuePredicate(ctx, cff.PredicateJob{
			Run: pred1.run,
			Dependencies: []*cff.ScheduledJob{
				task5.job,
			},
		})

		// go.uber.org/cff/internal/tests/benchmark/benchmark_predicate.go:111:4
		var (
			v5 skippedOut
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

				if recovered == nil && p0PanicRecover != nil {
					recovered = p0PanicRecover
				}
				if recovered != nil {
					taskEmitter.TaskPanicRecovered(ctx, recovered)
					v5, err = _113_21, nil
				}
			}()

			if !p0 {
				return nil
			}

			defer task6.ran.Store(true)

			v5, err = _111_4(v3)

			if err != nil {
				taskEmitter.TaskErrorRecovered(ctx, err)
				v5, err = _113_21, nil
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task6.job = sched.Enqueue(ctx, cff.Job{
			Run: task6.run,
			Dependencies: []*cff.ScheduledJob{
				task4.job,
				pred1.job,
			},
		})
		tasks = append(tasks, task6)

		// go.uber.org/cff/internal/tests/benchmark/benchmark_predicate.go:115:12
		var (
			v6 time.Duration
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
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: debug.Stack(),
					}
				}
			}()

			defer task7.ran.Store(true)

			v6 = _115_12(v5)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task7.job = sched.Enqueue(ctx, cff.Job{
			Run: task7.run,
			Dependencies: []*cff.ScheduledJob{
				task6.job,
			},
		})
		tasks = append(tasks, task7)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_104_15) = v6 // time.Duration

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
	return elapsed
}
