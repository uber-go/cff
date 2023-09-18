//go:build !cff
// +build !cff

package externalpackage

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/gofrs/uuid"
	"go.uber.org/cff"
	"go.uber.org/cff/internal/tests/externalpackage/external"
)

// NestedType is a flow that has a dep on an external struct.
func NestedType(ctx context.Context, driverUUID uuid.UUID) error {
	var ok bool
	return func() (err error) {

		_17_18 := ctx

		_18_14 := driverUUID

		_19_15 := &ok

		_21_12 := func(c uuid.UUID) bool {
			return true
		}
		ctx := _17_18
		var v1 uuid.UUID = _18_14
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/externalpackage/externalpackage.go",
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

		// go.uber.org/cff/internal/tests/externalpackage/externalpackage.go:21:12
		var (
			v2 bool
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
						Stacktrace: string(debug.Stack()),
					}
				}
			}()

			defer task0.ran.Store(true)

			v2 = _21_12(v1)

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

		*(_19_15) = v2 // bool

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
}

// ImplicitType is a flow that has a dep on an external struct where the package is not imported in this file
func ImplicitType(ctx context.Context) error {
	var ok bool
	return func() (err error) {

		_30_18 := ctx

		_31_15 := &ok

		_33_12 := external.ProvidesUUID

		_34_12 := external.NeedsUUID
		ctx := _30_18
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/externalpackage/externalpackage.go",
				Line:   30,
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

		// go.uber.org/cff/internal/tests/externalpackage/externalpackage.go:33:12
		var (
			v1 uuid.UUID
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
						Stacktrace: string(debug.Stack()),
					}
				}
			}()

			defer task1.ran.Store(true)

			v1 = _33_12()

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task1.job = sched.Enqueue(ctx, cff.Job{
			Run: task1.run,
		})
		tasks = append(tasks, task1)

		// go.uber.org/cff/internal/tests/externalpackage/externalpackage.go:34:12
		var (
			v2 bool
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
						Stacktrace: string(debug.Stack()),
					}
				}
			}()

			defer task2.ran.Store(true)

			v2 = _34_12(v1)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task2.job = sched.Enqueue(ctx, cff.Job{
			Run: task2.run,
			Dependencies: []*cff.ScheduledJob{
				task1.job,
			},
		})
		tasks = append(tasks, task2)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_31_15) = v2 // bool

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
}
