//go:build !cff
// +build !cff

package importcollision

import (
	"context"
	_template "html/template"
	"runtime/debug"
	__template "text/template"
	"time"

	packagewithdash "go.uber.org/cff/internal/tests/importcollision/package-with-dash"
	"go.uber.org/cff/internal/tests/importcollision/template"

	cff2 "go.uber.org/cff"
)

// Flow tests a flow that requires code generation to resolve multiple imports with
// the same base path.
func Flow() (string, error) {
	var result string
	err := func() (err error) {

		_19_3 := context.Background()

		_20_16 := &result

		_21_13 := GetHTMLTemplate

		_22_13 := GetTextTemplate

		_23_13 := GetFoo

		_24_13 := GetResult

		_25_13 := template.GetError
		ctx := _19_3
		emitter := cff2.NopEmitter()

		var (
			flowInfo = &cff2.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/importcollision/import_collision.go",
				Line:   18,
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

		// go.uber.org/cff/internal/tests/importcollision/import_collision.go:21:13
		var (
			v1 *_template.Template
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
					err = &cff2.PanicError{
						Value:      recovered,
						Stacktrace: string(debug.Stack()),
					}
				}
			}()

			defer task0.ran.Store(true)

			v1 = _21_13()

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task0.job = sched.Enqueue(ctx, cff2.Job{
			Run: task0.run,
		})
		tasks = append(tasks, task0)

		// go.uber.org/cff/internal/tests/importcollision/import_collision.go:22:13
		var (
			v2 *__template.Template
		)
		task1 := new(struct {
			emitter cff2.TaskEmitter
			ran     cff2.AtomicBool
			run     func(context.Context) error
			job     *cff2.ScheduledJob
		})
		task1.emitter = cff2.NopTaskEmitter()
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
					err = &cff2.PanicError{
						Value:      recovered,
						Stacktrace: string(debug.Stack()),
					}
				}
			}()

			defer task1.ran.Store(true)

			v2 = _22_13()

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task1.job = sched.Enqueue(ctx, cff2.Job{
			Run: task1.run,
		})
		tasks = append(tasks, task1)

		// go.uber.org/cff/internal/tests/importcollision/import_collision.go:23:13
		var (
			v3 packagewithdash.Foo
		)
		task2 := new(struct {
			emitter cff2.TaskEmitter
			ran     cff2.AtomicBool
			run     func(context.Context) error
			job     *cff2.ScheduledJob
		})
		task2.emitter = cff2.NopTaskEmitter()
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
					err = &cff2.PanicError{
						Value:      recovered,
						Stacktrace: string(debug.Stack()),
					}
				}
			}()

			defer task2.ran.Store(true)

			v3 = _23_13()

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task2.job = sched.Enqueue(ctx, cff2.Job{
			Run: task2.run,
		})
		tasks = append(tasks, task2)

		// go.uber.org/cff/internal/tests/importcollision/import_collision.go:24:13
		var (
			v4 string
		)
		task3 := new(struct {
			emitter cff2.TaskEmitter
			ran     cff2.AtomicBool
			run     func(context.Context) error
			job     *cff2.ScheduledJob
		})
		task3.emitter = cff2.NopTaskEmitter()
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
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = &cff2.PanicError{
						Value:      recovered,
						Stacktrace: string(debug.Stack()),
					}
				}
			}()

			defer task3.ran.Store(true)

			v4 = _24_13(v1, v2, v3)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task3.job = sched.Enqueue(ctx, cff2.Job{
			Run: task3.run,
			Dependencies: []*cff2.ScheduledJob{
				task0.job,
				task1.job,
				task2.job,
			},
		})
		tasks = append(tasks, task3)

		// go.uber.org/cff/internal/tests/importcollision/import_collision.go:25:13
		task4 := new(struct {
			emitter cff2.TaskEmitter
			ran     cff2.AtomicBool
			run     func(context.Context) error
			job     *cff2.ScheduledJob
		})
		task4.emitter = cff2.NopTaskEmitter()
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
					err = &cff2.PanicError{
						Value:      recovered,
						Stacktrace: string(debug.Stack()),
					}
				}
			}()

			defer task4.ran.Store(true)

			err = _25_13()

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task4.job = sched.Enqueue(ctx, cff2.Job{
			Run: task4.run,
		})
		tasks = append(tasks, task4)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_20_16) = v4 // string

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()

	return result, err
}
