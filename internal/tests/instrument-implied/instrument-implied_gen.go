//go:build !cff
// +build !cff

package instrumentimplied

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/uber-go/tally"
	"go.uber.org/cff"
	"go.uber.org/zap"
)

// Parent has already generated a sibling of this file,
// but we'll do it again with a diferent flag.

//go:generate cff -auto-instrument ./...

// H is used by some tests
type H struct {
	Scope  tally.Scope
	Logger *zap.Logger
}

// ImpliedName is a flow with a task with no instrument call but the flow is instrumented
func (h *H) ImpliedName(ctx context.Context, req string) (res int, err error) {
	var unsigned uint

	err = func() (err error) {

		_30_17 := ctx

		_31_14 := req

		_32_15 := &res

		_32_21 := &unsigned

		_33_15 := &unsigned

		_34_19 := cff.TallyEmitter(h.Scope)

		_35_19 := cff.LogEmitter(h.Logger)

		_36_22 := "ImpliedName"

		_38_4 := strconv.Atoi

		_42_4 := func(i int) (uint, error) {
			return uint(i), nil
		}
		ctx := _30_17
		var v1 string = _31_14
		emitter := cff.EmitterStack(_34_19, _35_19)

		var (
			flowInfo = &cff.FlowInfo{
				Name:   _36_22,
				File:   "go.uber.org/cff/internal/tests/instrument-implied/instrument-implied.go",
				Line:   30,
				Column: 8,
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

		// go.uber.org/cff/internal/tests/instrument-implied/instrument-implied.go:38:4
		var (
			v2 int
		)
		task0 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task0.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   "instrument-implied.go.38",
				File:   "go.uber.org/cff/internal/tests/instrument-implied/instrument-implied.go",
				Line:   38,
				Column: 4,
			},
			&cff.DirectiveInfo{
				Name:      flowInfo.Name,
				Directive: cff.FlowDirective,
				File:      flowInfo.File,
				Line:      flowInfo.Line,
				Column:    flowInfo.Column,
			},
		)
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

			v2, err = _38_4(v1)

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

		// go.uber.org/cff/internal/tests/instrument-implied/instrument-implied.go:42:4
		var (
			v3 uint
		)
		task1 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task1.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   "instrument-implied.go.42",
				File:   "go.uber.org/cff/internal/tests/instrument-implied/instrument-implied.go",
				Line:   42,
				Column: 4,
			},
			&cff.DirectiveInfo{
				Name:      flowInfo.Name,
				Directive: cff.FlowDirective,
				File:      flowInfo.File,
				Line:      flowInfo.Line,
				Column:    flowInfo.Column,
			},
		)
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

			v3, err = _42_4(v2)

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
			Dependencies: []*cff.ScheduledJob{
				task0.job,
			},
		})
		tasks = append(tasks, task1)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_32_15) = v2 // int
		*(_32_21) = v3 // uint
		*(_33_15) = v3 // uint

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
	return
}
