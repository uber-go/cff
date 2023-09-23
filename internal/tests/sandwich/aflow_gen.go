//go:build !cff
// +build !cff

package sandwich

import (
	"context"
	"runtime/debug"
	"time"

	"go.uber.org/cff"
)

func aFlow() (s string, err error) {
	err = func() (err error) {

		_14_3 := context.Background()

		_15_15 := &s

		_17_4 := aFunc
		ctx := _14_3
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/sandwich/aflow.go",
				Line:   13,
				Column: 8,
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

		// go.uber.org/cff/internal/tests/sandwich/aflow.go:17:4
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
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: debug.Stack(),
					}
				}
			}()

			defer task0.ran.Store(true)

			v1 = _17_4()

			taskEmitter.TaskSuccess(ctx)

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

		*(_15_15) = v1 // string

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()

	return s, err
}
