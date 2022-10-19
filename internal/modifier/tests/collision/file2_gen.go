//go:build !cff
// +build !cff

package collision

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/cff"
)

// Flow2 is a very simple flow that returns 2
func Flow2() (int, error) {
	var i int
	err := _cffFlowfile2_15_9(context.Background(),
		_cffConcurrencyfile2_16_3(1),
		_cffResultsfile2_17_3(&i),
		_cffTaskfile2_18_3(
			func() (int, error) {
				return 2, nil
			},
		),
	)
	return i, err
}
func _cffFlowfile2_15_9(
	ctx context.Context,
	mfile216_3 func() int,
	mfile217_3 func() *int,
	mfile218_3 func() func() (int, error),
) error {
	_16_19 := mfile216_3()
	_ = _16_19 // possibly unused.
	_17_15 := mfile217_3()
	_ = _17_15 // possibly unused.
	_19_4 := mfile218_3()
	_ = _19_4 // possibly unused.

	emitter := cff.NopEmitter()

	var (
		flowInfo = &cff.FlowInfo{
			File:   "go.uber.org/cff/internal/modifier/tests/collision/file2.go",
			Line:   15,
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
			Concurrency: _16_19, Emitter: schedEmitter,
		},
	)

	var tasks []*struct {
		emitter cff.TaskEmitter
		ran     cff.AtomicBool
		run     func(context.Context) error
		job     *cff.ScheduledJob
	}

	// go.uber.org/cff/internal/modifier/tests/collision/file2.go:19:4
	var (
		v1 int
	)
	task0 := new(struct {
		emitter cff.TaskEmitter
		ran     cff.AtomicBool
		run     func(context.Context) error
		job     *cff.ScheduledJob
	})

	task0.run = func(ctx context.Context) (err error) {
		defer func() {
			recovered := recover()
			if recovered != nil {
				err = fmt.Errorf("task panic: %v", recovered)
			}
		}()

		v1, err = _19_4()
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

	*(_17_15) = v1 // int

	flowEmitter.FlowSuccess(ctx)
	return nil
}

func _cffConcurrencyfile2_16_3(c int) func() int {
	return func() int { return c }
}

func _cffResultsfile2_17_3(mfile217_15 *int) func() *int {
	return func() *int { return mfile217_15 }
}

func _cffTaskfile2_18_3(mfile219_4 func() (int, error)) func() func() (int, error) {
	return func() func() (int, error) { return mfile219_4 }
}
