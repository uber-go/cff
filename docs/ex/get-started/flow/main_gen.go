// region directive
//go:build !cff

package main

// endregion directive

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"time"

	"go.uber.org/cff"
)

// region fake-client
var uber UberAPI = new(fakeUberClient)

// endregion fake-client

// region resp-decl
type Response struct {
	Rider    string
	Driver   string
	HomeCity string
}

// endregion resp-decl

// region main
func main() {
	// endregion main
	// region resp-var
	var res *Response
	// endregion resp-var
	// region ctx
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// endregion ctx
	// region flow-start
	err := func() (err error) {

		_42_18 := ctx

		_49_14 := 12

		_51_15 := &res

		_53_12 := func(tripID int) (*Trip, error) {

			return uber.TripByID(tripID)
		}

		_59_12 := func(trip *Trip) (*Driver, error) {
			return uber.DriverByID(trip.DriverID)
		}

		_64_12 := func(trip *Trip) (*Rider, error) {
			return uber.RiderByID(trip.RiderID)
		}

		_69_12 := func(rider *Rider) (*Location, error) {
			return uber.LocationByID(rider.HomeID)
		}

		_75_12 := func(r *Rider, d *Driver, home *Location) *Response {
			return &Response{
				Driver:   d.Name,
				Rider:    r.Name,
				HomeCity: home.City,
			}
		}
		ctx := _42_18
		var v1 int = _49_14
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/docs/ex/get-started/flow/main.go",
				Line:   42,
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

		// go.uber.org/cff/docs/ex/get-started/flow/main.go:53:12
		var (
			v2 *Trip
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

			v2, err = _53_12(v1)

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

		// go.uber.org/cff/docs/ex/get-started/flow/main.go:59:12
		var (
			v3 *Driver
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

			v3, err = _59_12(v2)

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

		// go.uber.org/cff/docs/ex/get-started/flow/main.go:64:12
		var (
			v4 *Rider
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
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: stacktrace,
					}
				}
			}()

			defer task2.ran.Store(true)

			v4, err = _64_12(v2)

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task2.job = sched.Enqueue(ctx, cff.Job{
			Run: task2.run,
			Dependencies: []*cff.ScheduledJob{
				task0.job,
			},
		})
		tasks = append(tasks, task2)

		// go.uber.org/cff/docs/ex/get-started/flow/main.go:69:12
		var (
			v5 *Location
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
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: stacktrace,
					}
				}
			}()

			defer task3.ran.Store(true)

			v5, err = _69_12(v4)

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task3.job = sched.Enqueue(ctx, cff.Job{
			Run: task3.run,
			Dependencies: []*cff.ScheduledJob{
				task2.job,
			},
		})
		tasks = append(tasks, task3)

		// go.uber.org/cff/docs/ex/get-started/flow/main.go:75:12
		var (
			v6 *Response
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

			defer task4.ran.Store(true)

			v6 = _75_12(v4, v3, v5)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task4.job = sched.Enqueue(ctx, cff.Job{
			Run: task4.run,
			Dependencies: []*cff.ScheduledJob{
				task2.job,
				task1.job,
				task3.job,
			},
		})
		tasks = append(tasks, task4)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_51_15) = v6 // *go.uber.org/cff/docs/ex/get-started/flow.Response

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
	// region tail
	if err != nil {
		log.Fatal(err)
	}
	// endregion error

	fmt.Println(res.Driver, "drove", res.Rider, "who lives in", res.HomeCity)
	// endregion tail
}
