//go:build !cff
// +build !cff

package setconcurrency

import (
	"bufio"
	"bytes"
	"context"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"go.uber.org/cff"
)

// This must be updated if scheduler.worker is renamed.
const _workerFunction = "go.uber.org/cff/scheduler.worker"

// NumWorkers runs a cff flow with the provided concurrency, and reports the
// number of workers from within the flow.
func NumWorkers(conc int) (int, error) {
	var numGoroutines int

	err := func() (err error) {

		_26_3 := context.Background()

		_27_19 := conc

		_28_15 := &numGoroutines

		_33_12 := func() (int, error) {
			return numWorkersStable(10, time.Millisecond)
		}
		ctx := _26_3
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/setconcurrency/setconcurrency.go",
				Line:   25,
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
				Concurrency: _27_19, Emitter: schedEmitter,
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

		// go.uber.org/cff/internal/tests/setconcurrency/setconcurrency.go:33:12
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
				var stacktrace string
				if recovered != nil {
					stacktrace = string(debug.Stack())
				}
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: stacktrace,
					}
				}
			}()

			defer task0.ran.Store(true)

			v1, err = _33_12()

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

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_28_15) = v1 // int

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()

	return numGoroutines, err
}

// NumWorkersNoArg runs a cff flow, and reports the
// number of workers from within the flow.
func NumWorkersNoArg() (int, error) {
	var numGoroutines int

	err := func() (err error) {

		_47_3 := context.Background()

		_48_15 := &numGoroutines

		_53_12 := func() (int, error) {
			return numWorkersStable(10, time.Millisecond)
		}
		ctx := _47_3
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/setconcurrency/setconcurrency.go",
				Line:   46,
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

		// go.uber.org/cff/internal/tests/setconcurrency/setconcurrency.go:53:12
		var (
			v1 int
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
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: stacktrace,
					}
				}
			}()

			defer task1.ran.Store(true)

			v1, err = _53_12()

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

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_48_15) = v1 // int

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()

	return numGoroutines, err
}

// numWorkersStable waits for the number of workers reported by numWorkers to
// stabilize for n ticks before reporting it.
func numWorkersStable(n int, tick time.Duration) (int, error) {
	numw, err := numWorkers()
	if err != nil {
		return 0, err
	}

	for remaining := n; remaining > 0; {
		time.Sleep(tick)
		next, err := numWorkers()
		if err != nil {
			return 0, err
		}

		if numw == next {
			remaining--
		} else {
			numw = next
			remaining = n
		}
	}

	return numw, nil
}

// numWorkers reports the number of goroutines currently running the cff
// scheduler's worker function.
func numWorkers() (int, error) {
	scanner := bufio.NewScanner(bytes.NewReader(getStack()))

	var (
		workers int
		inStack bool
	)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case !inStack && strings.HasPrefix(line, "goroutine "):
			// goroutine 42 [running]:
			inStack = true
		case inStack:
			// path/to/package.function(...)
			if strings.HasPrefix(line, _workerFunction+"(") {
				workers++
			}
		case len(line) == 0:
			inStack = false
		}
	}

	return workers, scanner.Err()
}

// getStack retrieves a stack trace for all running goroutines using
// runtime.Stack.
func getStack() []byte {
	const bufferSize = 64 * 1024 // 64kb

	// runtime.Stack reports the number of bytes actually written to the
	// buffer. If the buffer wasn't large enough, it stops writing. To
	// make sure we have the full stack, we'll double the buffer until we
	// have one large enough to hold the full stack trace.
	for size := bufferSize; ; size *= 2 {
		buf := make([]byte, size)
		if n := runtime.Stack(buf, true); n < size {
			return buf[:n]
		}
	}
}
