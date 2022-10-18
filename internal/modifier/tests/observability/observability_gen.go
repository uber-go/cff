//go:build !cff
// +build !cff

package observability

import (
	"context"
	"fmt"
	"time"

	"github.com/uber-go/tally"
	"go.uber.org/cff"
	"go.uber.org/zap"
)

// InstrumentFlow is a cff.Flow with an instrumented flow.
func InstrumentFlow(scope tally.Scope, logger *zap.Logger) (int64, error) {
	var res int64
	err := func() (err error) {

		_17_18 := context.Background()

		_18_19 := 2

		_19_14 := 1

		_20_15 := &res

		_21_19 := cff.TallyEmitter(scope)

		_22_19 := cff.LogEmitter(logger)

		_23_22 := "Instrumented"

		_25_4 := func(i int) int64 {
			return int64(1)
		}
		ctx := _17_18
		var v1 int = _19_14
		emitter := cff.EmitterStack(_21_19, _22_19)

		var (
			flowInfo = &cff.FlowInfo{
				Name:   _23_22,
				File:   "go.uber.org/cff/internal/modifier/tests/observability/observability.go",
				Line:   17,
				Column: 9,
			}
			flowEmitter = emitter.FlowInit(flowInfo)

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

		sched := cff.BeginFlow(
			cff.SchedulerParams{
				Concurrency: _18_19, Emitter: schedEmitter,
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

		// go.uber.org/cff/internal/modifier/tests/observability/observability.go:25:4
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

			v2 = _25_4(v1)

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

		*(_20_15) = v2 // int64

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
	return res, err
}
