//go:build !cff
// +build !cff

package namedimports

import (
	newctx "context"
	"runtime/debug"
	"time"

	cffv2 "go.uber.org/cff"
)

func run(ctx newctx.Context) error {
	var result struct{}
	return func() (err error) {

		_14_20 := ctx

		_15_16 := "foo"

		_16_17 := &result

		_18_4 := func(string) struct{} {
			panic("don't call me")
		}
		ctx := _14_20
		var v1 string = _15_16
		emitter := cffv2.NopEmitter()

		var (
			flowInfo = &cffv2.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/named_imports/named_imports.go",
				Line:   14,
				Column: 9,
			}
			flowEmitter = cffv2.NopFlowEmitter()

			schedInfo = &cffv2.SchedulerInfo{
				Name:      flowInfo.Name,
				Directive: cffv2.FlowDirective,
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

		sched := cffv2.NewScheduler(
			cffv2.SchedulerParams{
				Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cffv2.TaskEmitter
			ran     cffv2.AtomicBool
			run     func(newctx.Context) error
			job     *cffv2.ScheduledJob
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/named_imports/named_imports.go:18:4
		var (
			v2 struct{}
		)
		task0 := new(struct {
			emitter cffv2.TaskEmitter
			ran     cffv2.AtomicBool
			run     func(newctx.Context) error
			job     *cffv2.ScheduledJob
		})
		task0.emitter = cffv2.NopTaskEmitter()
		task0.run = func(ctx newctx.Context) (err error) {
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
					err = &cffv2.PanicError{
						Value:      recovered,
						Stacktrace: debug.Stack(),
					}
				}
			}()

			defer task0.ran.Store(true)

			v2 = _18_4(v1)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task0.job = sched.Enqueue(ctx, cffv2.Job{
			Run: task0.run,
		})
		tasks = append(tasks, task0)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_16_17) = v2 // struct{}

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
}
