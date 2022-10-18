//go:build !cff
// +build !cff

// Package instrument verifies that default and custom Emitter
// implementations trigger on events.
// DefaultEmitter tests default emitter.
// These tests will be removed in the future as an implementation detail.
// CustomEmitter tests mocks for custom emitter.
package instrument

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/uber-go/tally"
	"go.uber.org/atomic"
	"go.uber.org/cff"
	"go.uber.org/zap"
)

func main() {
	scope := tally.NoopScope
	logger := zap.NewNop()
	h := &DefaultEmitter{
		Scope:  scope,
		Logger: logger,
	}
	ctx := context.Background()
	res, err := h.RunFlow(ctx, os.Args[1])
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d\n", res)
}

// DefaultEmitter is used by other tests.
type DefaultEmitter struct {
	Scope  tally.Scope
	Logger *zap.Logger
}

// RunFlow executes a flow to test instrumentation.
func (h *DefaultEmitter) RunFlow(ctx context.Context, req string) (res uint8, err error) {
	err = func() (err error) {

		_47_17 := ctx

		_48_14 := req

		_49_15 := &res

		_50_19 := cff.TallyEmitter(h.Scope)

		_51_19 := cff.LogEmitter(h.Logger)

		_52_22 := "AtoiRun"

		_55_4 := strconv.Atoi

		_56_19 := "Atoi"

		_60_4 := func(i int) (uint8, error) {
			if i > -1 && i < 256 {
				return uint8(i), nil
			}
			return 0, errors.New("int can not fit into 8 bits")
		}

		_66_21 := uint8(0)

		_67_19 := "uint8"
		ctx := _47_17
		var v1 string = _48_14
		emitter := cff.EmitterStack(_50_19, _51_19)

		var (
			flowInfo = &cff.FlowInfo{
				Name:   _52_22,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   47,
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

		// go.uber.org/cff/internal/tests/instrument/instrument.go:55:4
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
				Name:   _56_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   55,
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

			v2, err = _55_4(v1)

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

		// go.uber.org/cff/internal/tests/instrument/instrument.go:60:4
		var (
			v3 uint8
		)
		task1 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task1.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _67_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   60,
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
					taskEmitter.TaskPanicRecovered(ctx, recovered)
					v3, err = _66_21, nil
				}
			}()

			defer task1.ran.Store(true)

			v3, err = _60_4(v2)

			if err != nil {
				taskEmitter.TaskErrorRecovered(ctx, err)
				v3, err = _66_21, nil
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

		*(_49_15) = v3 // uint8

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
	return
}

// RunParallelTasksAndTask executes parallel cff.Tasks and cff.Task with
// directive-level instrumentation.
func (h *DefaultEmitter) RunParallelTasksAndTask(ctx context.Context, req string) error {
	fn := func() error {
		_, err := strconv.Atoi(req)
		return err
	}
	return func() (err error) {

		_80_22 := ctx

		_81_19 := cff.TallyEmitter(h.Scope)

		_82_19 := cff.LogEmitter(h.Logger)

		_83_26 := "RunParallelTasksAndTask"

		_84_13 := fn

		_86_4 := fn

		_87_19 := "Atoi"
		ctx := _80_22
		emitter := cff.EmitterStack(_81_19, _82_19)

		var (
			parallelInfo = &cff.ParallelInfo{
				Name:   _83_26,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   80,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = emitter.ParallelInit(parallelInfo)

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.BeginFlow(
			cff.SchedulerParams{
				Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/instrument/instrument.go:84:13
		task2 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task2.emitter = cff.NopTaskEmitter()
		task2.fn = func(ctx context.Context) (err error) {
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
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task2.ran.Store(true)

			err = _84_13()

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return
			}
			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task2.fn,
		})
		tasks = append(tasks, task2)

		// go.uber.org/cff/internal/tests/instrument/instrument.go:86:4
		task3 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task3.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _87_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   86,
				Column: 4,
			},
			directiveInfo,
		)
		task3.fn = func(ctx context.Context) (err error) {
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
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task3.ran.Store(true)

			err = _86_4()

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return
			}
			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task3.fn,
		})
		tasks = append(tasks, task3)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line instrument.go:88*/
	}()
}

// RunParallelTasks executes a parallel cff.Task with directive-level
// instrumentation.
func (h *DefaultEmitter) RunParallelTasks(ctx context.Context, req string) error {
	return func() (err error) {

		_95_22 := ctx

		_96_19 := cff.TallyEmitter(h.Scope)

		_97_19 := cff.LogEmitter(h.Logger)

		_98_26 := "RunParallelTasks"

		_100_4 := func() error {
			_, err := strconv.Atoi(req)
			return err
		}
		ctx := _95_22
		emitter := cff.EmitterStack(_96_19, _97_19)

		var (
			parallelInfo = &cff.ParallelInfo{
				Name:   _98_26,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   95,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = emitter.ParallelInit(parallelInfo)

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.BeginFlow(
			cff.SchedulerParams{
				Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/instrument/instrument.go:100:4
		task4 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task4.emitter = cff.NopTaskEmitter()
		task4.fn = func(ctx context.Context) (err error) {
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
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task4.ran.Store(true)

			err = _100_4()

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return
			}
			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task4.fn,
		})
		tasks = append(tasks, task4)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line instrument.go:104*/
	}()
}

// RunParallelTask executes a parallel with directive-level and task level
// instrumentation.
func (h *DefaultEmitter) RunParallelTask(ctx context.Context, req string) error {
	return func() (err error) {

		_111_22 := ctx

		_112_19 := cff.TallyEmitter(h.Scope)

		_113_19 := cff.LogEmitter(h.Logger)

		_114_26 := "RunParallelTask"

		_116_4 := func() error {
			_, err := strconv.Atoi(req)
			return err
		}

		_120_19 := "Atoi"
		ctx := _111_22
		emitter := cff.EmitterStack(_112_19, _113_19)

		var (
			parallelInfo = &cff.ParallelInfo{
				Name:   _114_26,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   111,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = emitter.ParallelInit(parallelInfo)

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.BeginFlow(
			cff.SchedulerParams{
				Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/instrument/instrument.go:116:4
		task5 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task5.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _120_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   116,
				Column: 4,
			},
			directiveInfo,
		)
		task5.fn = func(ctx context.Context) (err error) {
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
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task5.ran.Store(true)

			err = _116_4()

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return
			}
			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task5.fn,
		})
		tasks = append(tasks, task5)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line instrument.go:121*/
	}()
}

// ExplicitListOfFields is a flow with an explicit list of log fields.
func (h *DefaultEmitter) ExplicitListOfFields(ctx context.Context, req string) (res int, err error) {
	err = func() (err error) {

		_127_17 := ctx

		_128_14 := req

		_129_15 := &res

		_130_22 := "ExplicitListOfFields"

		_131_19 := cff.TallyEmitter(h.Scope)

		_132_19 := cff.LogEmitter(h.Logger)

		_134_4 := strconv.Atoi

		_135_19 := "Atoi"
		ctx := _127_17
		var v1 string = _128_14
		emitter := cff.EmitterStack(_131_19, _132_19)

		var (
			flowInfo = &cff.FlowInfo{
				Name:   _130_22,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   127,
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

		// go.uber.org/cff/internal/tests/instrument/instrument.go:134:4
		var (
			v2 int
		)
		task6 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task6.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _135_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   134,
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
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			defer task6.ran.Store(true)

			v2, err = _134_4(v1)

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task6.job = sched.Enqueue(ctx, cff.Job{
			Run: task6.run,
		})
		tasks = append(tasks, task6)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_129_15) = v2 // int

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
	return
}

// InstrumentFlowAndTask executes a flow to test instrumentation.
func (h *DefaultEmitter) InstrumentFlowAndTask(ctx context.Context, req string) (res int, err error) {
	err = func() (err error) {

		_143_17 := ctx

		_144_14 := req

		_145_15 := &res

		_146_22 := "AtoiDo"

		_147_19 := cff.TallyEmitter(h.Scope)

		_148_19 := cff.LogEmitter(h.Logger)

		_150_4 := strconv.Atoi

		_151_19 := "Atoi"
		ctx := _143_17
		var v1 string = _144_14
		emitter := cff.EmitterStack(_147_19, _148_19)

		var (
			flowInfo = &cff.FlowInfo{
				Name:   _146_22,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   143,
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

		// go.uber.org/cff/internal/tests/instrument/instrument.go:150:4
		var (
			v2 int
		)
		task7 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task7.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _151_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   150,
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
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			defer task7.ran.Store(true)

			v2, err = _150_4(v1)

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task7.job = sched.Enqueue(ctx, cff.Job{
			Run: task7.run,
		})
		tasks = append(tasks, task7)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_145_15) = v2 // int

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
	return
}

// FlowOnlyInstrumentTask executes a flow that only instruments a task, but
// not the flow directive.
func (h *DefaultEmitter) FlowOnlyInstrumentTask(ctx context.Context, req string) (res int, err error) {
	err = func() (err error) {

		_160_17 := ctx

		_161_14 := req

		_162_15 := &res

		_163_19 := cff.TallyEmitter(h.Scope)

		_164_19 := cff.LogEmitter(h.Logger)

		_166_4 := strconv.Atoi

		_167_19 := "Atoi"
		ctx := _160_17
		var v1 string = _161_14
		emitter := cff.EmitterStack(_163_19, _164_19)

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   160,
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

		// go.uber.org/cff/internal/tests/instrument/instrument.go:166:4
		var (
			v2 int
		)
		task8 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task8.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _167_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   166,
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
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			defer task8.ran.Store(true)

			v2, err = _166_4(v1)

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task8.job = sched.Enqueue(ctx, cff.Job{
			Run: task8.run,
		})
		tasks = append(tasks, task8)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_162_15) = v2 // int

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
	return
}

// ParallelOnlyInstrumentTask executes a parallel that only instruments the
// cff.Task.
func (h *DefaultEmitter) ParallelOnlyInstrumentTask(ctx context.Context, req string) error {
	return func() (err error) {

		_176_22 := ctx

		_178_19 := cff.TallyEmitter(h.Scope)

		_179_19 := cff.LogEmitter(h.Logger)

		_181_4 := func() error {
			_, err := strconv.Atoi(req)
			return err
		}

		_185_19 := "Atoi"
		ctx := _176_22
		emitter := cff.EmitterStack(_178_19, _179_19)

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   176,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.BeginFlow(
			cff.SchedulerParams{
				Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/instrument/instrument.go:181:4
		task9 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task9.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _185_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   181,
				Column: 4,
			},
			directiveInfo,
		)
		task9.fn = func(ctx context.Context) (err error) {
			taskEmitter := task9.emitter
			startTime := time.Now()
			defer func() {
				if task9.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task9.ran.Store(true)

			err = _181_4()

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return
			}
			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task9.fn,
		})
		tasks = append(tasks, task9)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line instrument.go:186*/
	}()
}

// T3630161 reproduces T3630161 by executing a flow that runs a task that failed, recovers, and then runs another task.
func (h *DefaultEmitter) T3630161(ctx context.Context) {
	var s string
	_ = func() (err error) {

		_193_15 := ctx

		_194_15 := &s

		_195_19 := cff.TallyEmitter(h.Scope)

		_196_19 := cff.LogEmitter(h.Logger)

		_197_22 := "T3630161"

		_200_4 := func() (string, error) {
			return "", errors.New("always errors")
		}

		_203_19 := "Err"

		_204_21 := "fallback value"

		_208_4 := func(s string) error {
			return nil
		}

		_211_19 := "End"
		ctx := _193_15
		emitter := cff.EmitterStack(_195_19, _196_19)

		var (
			flowInfo = &cff.FlowInfo{
				Name:   _197_22,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   193,
				Column: 6,
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

		// go.uber.org/cff/internal/tests/instrument/instrument.go:200:4
		var (
			v1 string
		)
		task10 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task10.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _203_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   200,
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
				if recovered != nil {
					taskEmitter.TaskPanicRecovered(ctx, recovered)
					v1, err = _204_21, nil
				}
			}()

			defer task10.ran.Store(true)

			v1, err = _200_4()

			if err != nil {
				taskEmitter.TaskErrorRecovered(ctx, err)
				v1, err = _204_21, nil
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task10.job = sched.Enqueue(ctx, cff.Job{
			Run: task10.run,
		})
		tasks = append(tasks, task10)

		// go.uber.org/cff/internal/tests/instrument/instrument.go:208:4
		task11 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task11.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _211_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   208,
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
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			defer task11.ran.Store(true)

			err = _208_4(v1)

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task11.job = sched.Enqueue(ctx, cff.Job{
			Run: task11.run,
			Dependencies: []*cff.ScheduledJob{
				task10.job,
			},
		})
		tasks = append(tasks, task11)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_194_15) = v1 // string

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
	return
}

// T3795761 reproduces T3795761 where a task that returns no error should only emit skipped metric if it was not run
func (h *DefaultEmitter) T3795761(ctx context.Context, shouldRun bool, shouldError bool) string {
	var s string
	_ = func() (err error) {

		_221_15 := ctx

		_222_15 := &s

		_223_19 := cff.TallyEmitter(h.Scope)

		_224_19 := cff.LogEmitter(h.Logger)

		_225_22 := "T3795761"

		_228_4 := func() int {
			return 0
		}

		_231_19 := "ProvidesInt"

		_235_4 := func(s int) (string, error) {
			if shouldError {
				return "", errors.New("err")
			}

			return "ok", nil
		}

		_242_18 := func() bool { return shouldRun }

		_243_19 := "NeedsInt"
		ctx := _221_15
		emitter := cff.EmitterStack(_223_19, _224_19)

		var (
			flowInfo = &cff.FlowInfo{
				Name:   _225_22,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   221,
				Column: 6,
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

		// go.uber.org/cff/internal/tests/instrument/instrument.go:228:4
		var (
			v2 int
		)
		task12 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task12.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _231_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   228,
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
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			defer task12.ran.Store(true)

			v2 = _228_4()

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task12.job = sched.Enqueue(ctx, cff.Job{
			Run: task12.run,
		})
		tasks = append(tasks, task12)

		// go.uber.org/cff/internal/tests/instrument/instrument.go:242:4
		var p0 bool
		var p0PanicRecover interface{}
		pred1 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) error
			job *cff.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p0PanicRecover = recovered
				}
			}()
			p0 = _242_18()
			return nil
		}

		pred1.job = sched.Enqueue(ctx, cff.Job{
			Run: pred1.run,
		})

		// go.uber.org/cff/internal/tests/instrument/instrument.go:235:4
		var (
			v1 string
		)
		task13 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task13.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _243_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   235,
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
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			if !p0 {
				return nil
			}

			defer task13.ran.Store(true)

			v1, err = _235_4(v2)

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task13.job = sched.Enqueue(ctx, cff.Job{
			Run: task13.run,
			Dependencies: []*cff.ScheduledJob{
				task12.job,
				pred1.job,
			},
		})
		tasks = append(tasks, task13)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_222_15) = v1 // string

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
	return s
}

// TaskLatencySkipped guards against regressino of T6278905 where task
// latency metrics are emitted when a task is skipped due to predicate.
func (h *DefaultEmitter) TaskLatencySkipped(ctx context.Context, shouldRun bool) {
	var s string
	_ = func() (err error) {

		_253_15 := ctx

		_254_15 := &s

		_255_19 := cff.TallyEmitter(h.Scope)

		_256_22 := "TaskLatencySkipped"

		_259_4 := func() string {
			return "ok"
		}

		_262_18 := func() bool { return shouldRun }

		_263_19 := "Task"
		ctx := _253_15
		emitter := cff.EmitterStack(_255_19)

		var (
			flowInfo = &cff.FlowInfo{
				Name:   _256_22,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   253,
				Column: 6,
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

		// go.uber.org/cff/internal/tests/instrument/instrument.go:262:4
		var p0 bool
		var p0PanicRecover interface{}
		pred1 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) error
			job *cff.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p0PanicRecover = recovered
				}
			}()
			p0 = _262_18()
			return nil
		}

		pred1.job = sched.Enqueue(ctx, cff.Job{
			Run: pred1.run,
		})

		// go.uber.org/cff/internal/tests/instrument/instrument.go:259:4
		var (
			v1 string
		)
		task14 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task14.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _263_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   259,
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
				if recovered == nil && p0PanicRecover != nil {
					recovered = p0PanicRecover
				}
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			if !p0 {
				return nil
			}

			defer task14.ran.Store(true)

			v1 = _259_4()

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task14.job = sched.Enqueue(ctx, cff.Job{
			Run: task14.run,
			Dependencies: []*cff.ScheduledJob{
				pred1.job,
			},
		})
		tasks = append(tasks, task14)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_254_15) = v1 // string

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
	return
}

// FlowAlwaysPanics tests a flow with a task that always panics.
func (h *DefaultEmitter) FlowAlwaysPanics(ctx context.Context) {
	var s string
	_ = func() (err error) {

		_272_15 := ctx

		_273_15 := &s

		_274_19 := cff.TallyEmitter(h.Scope)

		_275_22 := "Flow"

		_278_4 := func() string {
			panic("panic value")
		}

		_281_19 := "Task"
		ctx := _272_15
		emitter := cff.EmitterStack(_274_19)

		var (
			flowInfo = &cff.FlowInfo{
				Name:   _275_22,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   272,
				Column: 6,
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

		// go.uber.org/cff/internal/tests/instrument/instrument.go:278:4
		var (
			v1 string
		)
		task15 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task15.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _281_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   278,
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
		task15.run = func(ctx context.Context) (err error) {
			taskEmitter := task15.emitter
			startTime := time.Now()
			defer func() {
				if task15.ran.Load() {
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

			defer task15.ran.Store(true)

			v1 = _278_4()

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task15.job = sched.Enqueue(ctx, cff.Job{
			Run: task15.run,
		})
		tasks = append(tasks, task15)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_273_15) = v1 // string

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
	return
}

// PredicatePanics is a flow that runs a panicing task predicate.
func (h *DefaultEmitter) PredicatePanics(ctx context.Context) error {
	return func() (err error) {

		_289_18 := ctx

		_290_19 := cff.TallyEmitter(h.Scope)

		_291_22 := "Flow"

		_292_12 := func() {}

		_294_5 := func() bool {
			panic("sad times")
			return true
		}

		_300_19 := "PredicatePanics"
		ctx := _289_18
		emitter := cff.EmitterStack(_290_19)

		var (
			flowInfo = &cff.FlowInfo{
				Name:   _291_22,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   289,
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

		// go.uber.org/cff/internal/tests/instrument/instrument.go:293:4
		var p0 bool
		var p0PanicRecover interface{}
		pred1 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) error
			job *cff.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p0PanicRecover = recovered
				}
			}()
			p0 = _294_5()
			return nil
		}

		pred1.job = sched.Enqueue(ctx, cff.Job{
			Run: pred1.run,
		})

		// go.uber.org/cff/internal/tests/instrument/instrument.go:292:12
		task16 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task16.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _300_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   292,
				Column: 12,
			},
			&cff.DirectiveInfo{
				Name:      flowInfo.Name,
				Directive: cff.FlowDirective,
				File:      flowInfo.File,
				Line:      flowInfo.Line,
				Column:    flowInfo.Column,
			},
		)
		task16.run = func(ctx context.Context) (err error) {
			taskEmitter := task16.emitter
			startTime := time.Now()
			defer func() {
				if task16.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered == nil && p0PanicRecover != nil {
					recovered = p0PanicRecover
				}
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			if !p0 {
				return nil
			}

			defer task16.ran.Store(true)

			_292_12()

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task16.job = sched.Enqueue(ctx, cff.Job{
			Run: task16.run,
			Dependencies: []*cff.ScheduledJob{
				pred1.job,
			},
		})
		tasks = append(tasks, task16)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
}

// PredicatePanicsWithFallback is a flow that runs a panicing task predicate
// with a fallback.
func (h *DefaultEmitter) PredicatePanicsWithFallback(ctx context.Context) (string, error) {
	var (
		s   string
		err error
	)
	err = func() (err error) {

		_313_3 := ctx

		_314_15 := &s

		_315_19 := cff.TallyEmitter(h.Scope)

		_316_22 := "Flow"

		_318_4 := func(context.Context) (string, error) {
			return "value", nil
		}

		_322_5 := func() bool {
			panic("sad times")
			return true
		}

		_327_21 := "predicate-fallback"

		_328_19 := "PredicatePanicsWithFallback"
		ctx := _313_3
		emitter := cff.EmitterStack(_315_19)

		var (
			flowInfo = &cff.FlowInfo{
				Name:   _316_22,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   312,
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

		// go.uber.org/cff/internal/tests/instrument/instrument.go:321:4
		var p0 bool
		var p0PanicRecover interface{}
		pred1 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) error
			job *cff.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p0PanicRecover = recovered
				}
			}()
			p0 = _322_5()
			return nil
		}

		pred1.job = sched.Enqueue(ctx, cff.Job{
			Run: pred1.run,
		})

		// go.uber.org/cff/internal/tests/instrument/instrument.go:318:4
		var (
			v1 string
		)
		task17 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task17.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _328_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   318,
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
		task17.run = func(ctx context.Context) (err error) {
			taskEmitter := task17.emitter
			startTime := time.Now()
			defer func() {
				if task17.ran.Load() {
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
					v1, err = _327_21, nil
				}
			}()

			if !p0 {
				return nil
			}

			defer task17.ran.Store(true)

			v1, err = _318_4(ctx)

			if err != nil {
				taskEmitter.TaskErrorRecovered(ctx, err)
				v1, err = _327_21, nil
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task17.job = sched.Enqueue(ctx, cff.Job{
			Run: task17.run,
			Dependencies: []*cff.ScheduledJob{
				pred1.job,
			},
		})
		tasks = append(tasks, task17)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_314_15) = v1 // string

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
	return s, err
}

// ParallelAlwaysPanics tests a task which always panics.
func (h *DefaultEmitter) ParallelAlwaysPanics(ctx context.Context) {
	_ = func() (err error) {

		_336_19 := ctx

		_337_19 := cff.TallyEmitter(h.Scope)

		_338_26 := "Parallel"

		_340_4 := func() {
			panic("panic value")
		}

		_343_19 := "Trouble"
		ctx := _336_19
		emitter := cff.EmitterStack(_337_19)

		var (
			parallelInfo = &cff.ParallelInfo{
				Name:   _338_26,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   336,
				Column: 6,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = emitter.ParallelInit(parallelInfo)

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.BeginFlow(
			cff.SchedulerParams{
				Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/instrument/instrument.go:340:4
		task18 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task18.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _343_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   340,
				Column: 4,
			},
			directiveInfo,
		)
		task18.fn = func(ctx context.Context) (err error) {
			taskEmitter := task18.emitter
			startTime := time.Now()
			defer func() {
				if task18.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task18.ran.Store(true)

			_340_4()

			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task18.fn,
		})
		tasks = append(tasks, task18)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line instrument.go:344*/
	}()
	return
}

// ParallelTaskAlwaysPanics executes an instrumented cff.Task that always
// panics.
func (h *DefaultEmitter) ParallelTaskAlwaysPanics(ctx context.Context) {
	_ = func() (err error) {

		_352_19 := ctx

		_353_19 := cff.TallyEmitter(h.Scope)

		_354_26 := "Parallel"

		_356_4 := func() {
			panic("panic value")
		}

		_359_19 := "AlwaysPanic"
		ctx := _352_19
		emitter := cff.EmitterStack(_353_19)

		var (
			parallelInfo = &cff.ParallelInfo{
				Name:   _354_26,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   352,
				Column: 6,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = emitter.ParallelInit(parallelInfo)

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.BeginFlow(
			cff.SchedulerParams{
				Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/instrument/instrument.go:356:4
		task19 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task19.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _359_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   356,
				Column: 4,
			},
			directiveInfo,
		)
		task19.fn = func(ctx context.Context) (err error) {
			taskEmitter := task19.emitter
			startTime := time.Now()
			defer func() {
				if task19.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task19.ran.Store(true)

			_356_4()

			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task19.fn,
		})
		tasks = append(tasks, task19)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line instrument.go:360*/
	}()
	return
}

// These tests replicate the ones written for instrumentation to verify that
// custom Emitter will trigger similarly to default implementation.

// CustomEmitter is used by other tests.
type CustomEmitter struct {
	Scope   tally.Scope
	Logger  *zap.Logger
	Emitter cff.Emitter
}

// RunFlow executes a flow that instruments the top-level flow and tasks,
// of which one task can error.
func (h *CustomEmitter) RunFlow(ctx context.Context, req string) (res uint8, err error) {
	err = func() (err error) {

		_378_17 := ctx

		_379_14 := req

		_380_15 := &res

		_381_19 := cff.LogEmitter(h.Logger)

		_382_22 := "AtoiRun"

		_383_19 := h.Emitter

		_385_4 := strconv.Atoi

		_386_19 := "Atoi"

		_389_4 := func(i int) (uint8, error) {
			if i > -1 && i < 256 {
				return uint8(i), nil
			}
			return 0, errors.New("int can not fit into 8 bits")
		}

		_395_21 := uint8(0)

		_396_19 := "uint8"
		ctx := _378_17
		var v1 string = _379_14
		emitter := cff.EmitterStack(_381_19, _383_19)

		var (
			flowInfo = &cff.FlowInfo{
				Name:   _382_22,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   378,
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

		// go.uber.org/cff/internal/tests/instrument/instrument.go:385:4
		var (
			v2 int
		)
		task20 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task20.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _386_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   385,
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
		task20.run = func(ctx context.Context) (err error) {
			taskEmitter := task20.emitter
			startTime := time.Now()
			defer func() {
				if task20.ran.Load() {
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

			defer task20.ran.Store(true)

			v2, err = _385_4(v1)

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task20.job = sched.Enqueue(ctx, cff.Job{
			Run: task20.run,
		})
		tasks = append(tasks, task20)

		// go.uber.org/cff/internal/tests/instrument/instrument.go:389:4
		var (
			v3 uint8
		)
		task21 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task21.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _396_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   389,
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
		task21.run = func(ctx context.Context) (err error) {
			taskEmitter := task21.emitter
			startTime := time.Now()
			defer func() {
				if task21.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanicRecovered(ctx, recovered)
					v3, err = _395_21, nil
				}
			}()

			defer task21.ran.Store(true)

			v3, err = _389_4(v2)

			if err != nil {
				taskEmitter.TaskErrorRecovered(ctx, err)
				v3, err = _395_21, nil
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task21.job = sched.Enqueue(ctx, cff.Job{
			Run: task21.run,
			Dependencies: []*cff.ScheduledJob{
				task20.job,
			},
		})
		tasks = append(tasks, task21)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_380_15) = v3 // uint8

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
	return
}

// RunParallelTasksAndTask executes parallel cff.Tasks and cff.Task with
// directive-level instrumentation.
func (h *CustomEmitter) RunParallelTasksAndTask(ctx context.Context, req string) error {
	fn := func() error {
		_, err := strconv.Atoi(req)
		return err
	}
	return func() (err error) {

		_409_22 := ctx

		_410_19 := cff.TallyEmitter(h.Scope)

		_411_19 := h.Emitter

		_412_26 := "RunParallelTasksAndTask"

		_413_13 := fn

		_415_4 := fn

		_416_19 := "Atoi"
		ctx := _409_22
		emitter := cff.EmitterStack(_410_19, _411_19)

		var (
			parallelInfo = &cff.ParallelInfo{
				Name:   _412_26,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   409,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = emitter.ParallelInit(parallelInfo)

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.BeginFlow(
			cff.SchedulerParams{
				Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/instrument/instrument.go:413:13
		task22 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task22.emitter = cff.NopTaskEmitter()
		task22.fn = func(ctx context.Context) (err error) {
			taskEmitter := task22.emitter
			startTime := time.Now()
			defer func() {
				if task22.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task22.ran.Store(true)

			err = _413_13()

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return
			}
			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task22.fn,
		})
		tasks = append(tasks, task22)

		// go.uber.org/cff/internal/tests/instrument/instrument.go:415:4
		task23 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task23.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _416_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   415,
				Column: 4,
			},
			directiveInfo,
		)
		task23.fn = func(ctx context.Context) (err error) {
			taskEmitter := task23.emitter
			startTime := time.Now()
			defer func() {
				if task23.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task23.ran.Store(true)

			err = _415_4()

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return
			}
			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task23.fn,
		})
		tasks = append(tasks, task23)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line instrument.go:417*/
	}()
}

// RunParallelTasks executes parallel cff.Tasks with directive-level instrumentation.
func (h *CustomEmitter) RunParallelTasks(ctx context.Context, req string) error {
	return func() (err error) {

		_423_22 := ctx

		_424_19 := cff.LogEmitter(h.Logger)

		_425_19 := h.Emitter

		_426_26 := "RunParallelTasks"

		_428_4 := func() error {
			_, err := strconv.Atoi(req)
			return err
		}
		ctx := _423_22
		emitter := cff.EmitterStack(_424_19, _425_19)

		var (
			parallelInfo = &cff.ParallelInfo{
				Name:   _426_26,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   423,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = emitter.ParallelInit(parallelInfo)

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.BeginFlow(
			cff.SchedulerParams{
				Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/instrument/instrument.go:428:4
		task24 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task24.emitter = cff.NopTaskEmitter()
		task24.fn = func(ctx context.Context) (err error) {
			taskEmitter := task24.emitter
			startTime := time.Now()
			defer func() {
				if task24.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task24.ran.Store(true)

			err = _428_4()

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return
			}
			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task24.fn,
		})
		tasks = append(tasks, task24)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line instrument.go:432*/
	}()
}

// RunParallelTask executes a parallel to test instrumentation.
func (h *CustomEmitter) RunParallelTask(ctx context.Context, req string) error {
	return func() (err error) {

		_438_22 := ctx

		_439_19 := cff.TallyEmitter(h.Scope)

		_440_19 := h.Emitter

		_441_26 := "RunParallelTask"

		_443_4 := func() error {
			_, err := strconv.Atoi(req)
			return err
		}

		_447_19 := "Atoi"
		ctx := _438_22
		emitter := cff.EmitterStack(_439_19, _440_19)

		var (
			parallelInfo = &cff.ParallelInfo{
				Name:   _441_26,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   438,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = emitter.ParallelInit(parallelInfo)

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.BeginFlow(
			cff.SchedulerParams{
				Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/instrument/instrument.go:443:4
		task25 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task25.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _447_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   443,
				Column: 4,
			},
			directiveInfo,
		)
		task25.fn = func(ctx context.Context) (err error) {
			taskEmitter := task25.emitter
			startTime := time.Now()
			defer func() {
				if task25.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task25.ran.Store(true)

			err = _443_4()

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return
			}
			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task25.fn,
		})
		tasks = append(tasks, task25)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line instrument.go:448*/
	}()
}

// FlowOnlyInstrumentTask executes a flow that instruments the directive-level flow and
// the task.
func (h *CustomEmitter) FlowOnlyInstrumentTask(ctx context.Context, req string) (res int, err error) {
	err = func() (err error) {

		_455_17 := ctx

		_456_14 := req

		_457_15 := &res

		_458_19 := h.Emitter

		_459_19 := cff.LogEmitter(h.Logger)

		_461_4 := strconv.Atoi

		_462_19 := "Atoi"
		ctx := _455_17
		var v1 string = _456_14
		emitter := cff.EmitterStack(_458_19, _459_19)

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   455,
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

		// go.uber.org/cff/internal/tests/instrument/instrument.go:461:4
		var (
			v2 int
		)
		task26 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task26.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _462_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   461,
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
		task26.run = func(ctx context.Context) (err error) {
			taskEmitter := task26.emitter
			startTime := time.Now()
			defer func() {
				if task26.ran.Load() {
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

			defer task26.ran.Store(true)

			v2, err = _461_4(v1)

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task26.job = sched.Enqueue(ctx, cff.Job{
			Run: task26.run,
		})
		tasks = append(tasks, task26)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_457_15) = v2 // int

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
	return
}

// ParallelOnlyInstrumentTask executes a parallel that only instruments the
// individual task.
func (h *CustomEmitter) ParallelOnlyInstrumentTask(ctx context.Context, req string) error {
	return func() (err error) {

		_471_22 := ctx

		_472_19 := h.Emitter

		_473_19 := cff.LogEmitter(h.Logger)

		_475_4 := func() error {
			_, err := strconv.Atoi(req)
			return err
		}

		_479_19 := "Atoi"
		ctx := _471_22
		emitter := cff.EmitterStack(_472_19, _473_19)

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   471,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.BeginFlow(
			cff.SchedulerParams{
				Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/instrument/instrument.go:475:4
		task27 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task27.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _479_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   475,
				Column: 4,
			},
			directiveInfo,
		)
		task27.fn = func(ctx context.Context) (err error) {
			taskEmitter := task27.emitter
			startTime := time.Now()
			defer func() {
				if task27.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task27.ran.Store(true)

			err = _475_4()

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return
			}
			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task27.fn,
		})
		tasks = append(tasks, task27)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line instrument.go:480*/
	}()
}

// T3630161 reproduces T3630161 by executing a flow that runs a task that failed,
// recovers, and then runs another task.
func (h *CustomEmitter) T3630161(ctx context.Context) {
	var s string
	_ = func() (err error) {

		_488_15 := ctx

		_489_15 := &s

		_490_19 := h.Emitter

		_491_19 := cff.LogEmitter(h.Logger)

		_492_22 := "T3630161"

		_495_4 := func() (string, error) {
			return "", errors.New("always errors")
		}

		_498_19 := "Err"

		_499_21 := "fallback value"

		_503_4 := func(s string) error {
			return nil
		}

		_506_19 := "End"
		ctx := _488_15
		emitter := cff.EmitterStack(_490_19, _491_19)

		var (
			flowInfo = &cff.FlowInfo{
				Name:   _492_22,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   488,
				Column: 6,
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

		// go.uber.org/cff/internal/tests/instrument/instrument.go:495:4
		var (
			v1 string
		)
		task28 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task28.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _498_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   495,
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
		task28.run = func(ctx context.Context) (err error) {
			taskEmitter := task28.emitter
			startTime := time.Now()
			defer func() {
				if task28.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanicRecovered(ctx, recovered)
					v1, err = _499_21, nil
				}
			}()

			defer task28.ran.Store(true)

			v1, err = _495_4()

			if err != nil {
				taskEmitter.TaskErrorRecovered(ctx, err)
				v1, err = _499_21, nil
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task28.job = sched.Enqueue(ctx, cff.Job{
			Run: task28.run,
		})
		tasks = append(tasks, task28)

		// go.uber.org/cff/internal/tests/instrument/instrument.go:503:4
		task29 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task29.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _506_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   503,
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
		task29.run = func(ctx context.Context) (err error) {
			taskEmitter := task29.emitter
			startTime := time.Now()
			defer func() {
				if task29.ran.Load() {
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

			defer task29.ran.Store(true)

			err = _503_4(v1)

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task29.job = sched.Enqueue(ctx, cff.Job{
			Run: task29.run,
			Dependencies: []*cff.ScheduledJob{
				task28.job,
			},
		})
		tasks = append(tasks, task29)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_489_15) = v1 // string

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
	return
}

// T3795761 reproduces T3795761 where a task that returns no error should only
// emit skipped metric if it was not run.
func (h *CustomEmitter) T3795761(ctx context.Context, shouldRun bool,
	shouldError bool,
) string {
	var s string
	_ = func() (err error) {

		_519_15 := ctx

		_520_15 := &s

		_521_19 := h.Emitter

		_522_19 := cff.LogEmitter(h.Logger)

		_523_22 := "T3795761"

		_526_4 := func() int {
			return 0
		}

		_529_19 := "ProvidesInt"

		_533_4 := func(s int) (string, error) {
			if shouldError {
				return "", errors.New("err")
			}

			return "ok", nil
		}

		_540_18 := func() bool { return shouldRun }

		_541_19 := "NeedsInt"
		ctx := _519_15
		emitter := cff.EmitterStack(_521_19, _522_19)

		var (
			flowInfo = &cff.FlowInfo{
				Name:   _523_22,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   519,
				Column: 6,
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

		// go.uber.org/cff/internal/tests/instrument/instrument.go:526:4
		var (
			v2 int
		)
		task30 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task30.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _529_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   526,
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
		task30.run = func(ctx context.Context) (err error) {
			taskEmitter := task30.emitter
			startTime := time.Now()
			defer func() {
				if task30.ran.Load() {
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

			defer task30.ran.Store(true)

			v2 = _526_4()

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task30.job = sched.Enqueue(ctx, cff.Job{
			Run: task30.run,
		})
		tasks = append(tasks, task30)

		// go.uber.org/cff/internal/tests/instrument/instrument.go:540:4
		var p0 bool
		var p0PanicRecover interface{}
		pred1 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) error
			job *cff.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p0PanicRecover = recovered
				}
			}()
			p0 = _540_18()
			return nil
		}

		pred1.job = sched.Enqueue(ctx, cff.Job{
			Run: pred1.run,
		})

		// go.uber.org/cff/internal/tests/instrument/instrument.go:533:4
		var (
			v1 string
		)
		task31 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task31.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _541_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   533,
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
		task31.run = func(ctx context.Context) (err error) {
			taskEmitter := task31.emitter
			startTime := time.Now()
			defer func() {
				if task31.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered == nil && p0PanicRecover != nil {
					recovered = p0PanicRecover
				}
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			if !p0 {
				return nil
			}

			defer task31.ran.Store(true)

			v1, err = _533_4(v2)

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task31.job = sched.Enqueue(ctx, cff.Job{
			Run: task31.run,
			Dependencies: []*cff.ScheduledJob{
				task30.job,
				pred1.job,
			},
		})
		tasks = append(tasks, task31)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_520_15) = v1 // string

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
	return s
}

// FlowAlwaysPanics is a flow that tests Metrics Emitter.
func (h *CustomEmitter) FlowAlwaysPanics(ctx context.Context) error {
	return func() (err error) {

		_549_18 := ctx

		_550_19 := cff.LogEmitter(h.Logger)

		_551_19 := h.Emitter

		_552_12 := func() {
			panic("always")
		}

		_556_19 := "Atoi"
		ctx := _549_18
		emitter := cff.EmitterStack(_550_19, _551_19)

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   549,
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

		// go.uber.org/cff/internal/tests/instrument/instrument.go:552:12
		task32 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task32.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _556_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   552,
				Column: 12,
			},
			&cff.DirectiveInfo{
				Name:      flowInfo.Name,
				Directive: cff.FlowDirective,
				File:      flowInfo.File,
				Line:      flowInfo.Line,
				Column:    flowInfo.Column,
			},
		)
		task32.run = func(ctx context.Context) (err error) {
			taskEmitter := task32.emitter
			startTime := time.Now()
			defer func() {
				if task32.ran.Load() {
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

			defer task32.ran.Store(true)

			_552_12()

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task32.job = sched.Enqueue(ctx, cff.Job{
			Run: task32.run,
		})
		tasks = append(tasks, task32)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
}

// PredicatePanics is a flow that runs a panicing task predicate.
func (h *CustomEmitter) PredicatePanics(ctx context.Context) error {
	return func() (err error) {

		_563_18 := ctx

		_564_19 := h.Emitter

		_565_12 := func() {}

		_567_5 := func() bool {
			panic("sad times")
			return true
		}

		_573_19 := "PredicatePanics"
		ctx := _563_18
		emitter := cff.EmitterStack(_564_19)

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   563,
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

		// go.uber.org/cff/internal/tests/instrument/instrument.go:566:4
		var p0 bool
		var p0PanicRecover interface{}
		pred1 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) error
			job *cff.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p0PanicRecover = recovered
				}
			}()
			p0 = _567_5()
			return nil
		}

		pred1.job = sched.Enqueue(ctx, cff.Job{
			Run: pred1.run,
		})

		// go.uber.org/cff/internal/tests/instrument/instrument.go:565:12
		task33 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task33.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _573_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   565,
				Column: 12,
			},
			&cff.DirectiveInfo{
				Name:      flowInfo.Name,
				Directive: cff.FlowDirective,
				File:      flowInfo.File,
				Line:      flowInfo.Line,
				Column:    flowInfo.Column,
			},
		)
		task33.run = func(ctx context.Context) (err error) {
			taskEmitter := task33.emitter
			startTime := time.Now()
			defer func() {
				if task33.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered == nil && p0PanicRecover != nil {
					recovered = p0PanicRecover
				}
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			if !p0 {
				return nil
			}

			defer task33.ran.Store(true)

			_565_12()

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task33.job = sched.Enqueue(ctx, cff.Job{
			Run: task33.run,
			Dependencies: []*cff.ScheduledJob{
				pred1.job,
			},
		})
		tasks = append(tasks, task33)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
}

// PredicatePanicsWithFallback is a flow that runs a panicing task predicate
// with a fallback.
func (h *CustomEmitter) PredicatePanicsWithFallback(ctx context.Context) (string, error) {
	var (
		s   string
		err error
	)
	err = func() (err error) {

		_586_3 := ctx

		_587_15 := &s

		_588_19 := h.Emitter

		_590_4 := func(context.Context) (string, error) {
			return "value", nil
		}

		_594_5 := func() bool {
			panic("sad times")
			return true
		}

		_599_21 := "predicate-fallback"

		_600_19 := "PredicatePanicsWithFallback"
		ctx := _586_3
		emitter := cff.EmitterStack(_588_19)

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   585,
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

		// go.uber.org/cff/internal/tests/instrument/instrument.go:593:4
		var p0 bool
		var p0PanicRecover interface{}
		pred1 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) error
			job *cff.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p0PanicRecover = recovered
				}
			}()
			p0 = _594_5()
			return nil
		}

		pred1.job = sched.Enqueue(ctx, cff.Job{
			Run: pred1.run,
		})

		// go.uber.org/cff/internal/tests/instrument/instrument.go:590:4
		var (
			v1 string
		)
		task34 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task34.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _600_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   590,
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
		task34.run = func(ctx context.Context) (err error) {
			taskEmitter := task34.emitter
			startTime := time.Now()
			defer func() {
				if task34.ran.Load() {
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
					v1, err = _599_21, nil
				}
			}()

			if !p0 {
				return nil
			}

			defer task34.ran.Store(true)

			v1, err = _590_4(ctx)

			if err != nil {
				taskEmitter.TaskErrorRecovered(ctx, err)
				v1, err = _599_21, nil
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task34.job = sched.Enqueue(ctx, cff.Job{
			Run: task34.run,
			Dependencies: []*cff.ScheduledJob{
				pred1.job,
			},
		})
		tasks = append(tasks, task34)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_587_15) = v1 // string

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
	return s, err
}

// ParallelAlwaysPanics executes a directive-level instrument parallel with a
// cff.Tasks that panics.
func (h *CustomEmitter) ParallelAlwaysPanics(ctx context.Context) error {
	return func() (err error) {

		_609_22 := ctx

		_610_19 := cff.LogEmitter(h.Logger)

		_611_19 := h.Emitter

		_612_26 := "AlwaysPanic"

		_614_4 := func() {
			panic("always")
		}
		ctx := _609_22
		emitter := cff.EmitterStack(_610_19, _611_19)

		var (
			parallelInfo = &cff.ParallelInfo{
				Name:   _612_26,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   609,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = emitter.ParallelInit(parallelInfo)

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.BeginFlow(
			cff.SchedulerParams{
				Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/instrument/instrument.go:614:4
		task35 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task35.emitter = cff.NopTaskEmitter()
		task35.fn = func(ctx context.Context) (err error) {
			taskEmitter := task35.emitter
			startTime := time.Now()
			defer func() {
				if task35.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task35.ran.Store(true)

			_614_4()

			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task35.fn,
		})
		tasks = append(tasks, task35)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line instrument.go:617*/
	}()
}

// ParallelTaskAlwaysPanics executes an instrument cff.Task that always
// panics.
func (h *CustomEmitter) ParallelTaskAlwaysPanics(ctx context.Context) error {
	return func() (err error) {

		_624_22 := ctx

		_625_19 := cff.LogEmitter(h.Logger)

		_626_19 := h.Emitter

		_627_26 := "AlwaysPanic"

		_629_4 := func() {
			panic("always")
		}

		_632_19 := "Panic"
		ctx := _624_22
		emitter := cff.EmitterStack(_625_19, _626_19)

		var (
			parallelInfo = &cff.ParallelInfo{
				Name:   _627_26,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   624,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = emitter.ParallelInit(parallelInfo)

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.BeginFlow(
			cff.SchedulerParams{
				Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/instrument/instrument.go:629:4
		task36 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task36.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _632_19,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   629,
				Column: 4,
			},
			directiveInfo,
		)
		task36.fn = func(ctx context.Context) (err error) {
			taskEmitter := task36.emitter
			startTime := time.Now()
			defer func() {
				if task36.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task36.ran.Store(true)

			_629_4()

			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task36.fn,
		})
		tasks = append(tasks, task36)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line instrument.go:633*/
	}()
}

// FlowWithTwoEmitters is a flow that uses WithEmitter multiple times.
func FlowWithTwoEmitters(ctx context.Context, e1, e2 cff.Emitter, req string) (res int, err error) {
	err = func() (err error) {

		_639_17 := ctx

		_640_14 := req

		_641_15 := &res

		_642_19 := e1

		_643_19 := e2

		_644_22 := "AtoiDo"

		_645_12 := strconv.Atoi

		_645_41 := "Atoi"
		ctx := _639_17
		var v1 string = _640_14
		emitter := cff.EmitterStack(_642_19, _643_19)

		var (
			flowInfo = &cff.FlowInfo{
				Name:   _644_22,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   639,
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

		// go.uber.org/cff/internal/tests/instrument/instrument.go:645:12
		var (
			v2 int
		)
		task37 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task37.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _645_41,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   645,
				Column: 12,
			},
			&cff.DirectiveInfo{
				Name:      flowInfo.Name,
				Directive: cff.FlowDirective,
				File:      flowInfo.File,
				Line:      flowInfo.Line,
				Column:    flowInfo.Column,
			},
		)
		task37.run = func(ctx context.Context) (err error) {
			taskEmitter := task37.emitter
			startTime := time.Now()
			defer func() {
				if task37.ran.Load() {
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

			defer task37.ran.Store(true)

			v2, err = _645_12(v1)

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task37.job = sched.Enqueue(ctx, cff.Job{
			Run: task37.run,
		})
		tasks = append(tasks, task37)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_641_15) = v2 // int

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
	return
}

// ParallelWithTwoEmitters is a flow that uses WithEmitter multiple types.
func ParallelWithTwoEmitters(ctx context.Context, e1, e2 cff.Emitter, req string) (res int, err error) {
	var a atomic.Int64

	err = func() (err error) {

		_654_21 := ctx

		_655_19 := e1

		_656_19 := e2

		_657_26 := "AtoiDo"

		_659_4 := func() error {
			v, err := strconv.Atoi(req)
			a.Store(int64(v))
			return err
		}
		ctx := _654_21
		emitter := cff.EmitterStack(_655_19, _656_19)

		var (
			parallelInfo = &cff.ParallelInfo{
				Name:   _657_26,
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   654,
				Column: 8,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = emitter.ParallelInit(parallelInfo)

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.BeginFlow(
			cff.SchedulerParams{
				Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/instrument/instrument.go:659:4
		task38 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task38.emitter = cff.NopTaskEmitter()
		task38.fn = func(ctx context.Context) (err error) {
			taskEmitter := task38.emitter
			startTime := time.Now()
			defer func() {
				if task38.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task38.ran.Store(true)

			err = _659_4()

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return
			}
			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task38.fn,
		})
		tasks = append(tasks, task38)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line instrument.go:664*/
	}()
	res = int(a.Load())
	return
}
