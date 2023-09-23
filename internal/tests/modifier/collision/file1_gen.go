//go:build !cff
// +build !cff

package collision

import (
	"context"
	"runtime/debug"
	"time"

	"go.uber.org/cff"
)

// Flow1 is a very simple flow that returns 1
func Flow1() (int, error) {
	var i int
	err := _cffFlowfile1_15_9(context.Background(),
		_cffConcurrencyfile1_16_3(1),
		_cffResultsfile1_17_3(&i),
		_cffTaskfile1_18_3(
			func() (int, error) {
				return 1, nil
			},
		),
	)
	return i, err
}
func _cffFlowfile1_15_9(
	ctx context.Context,
	mfile116_3 func() int,
	mfile117_3 func() *int,
	mfile118_3 func() func() (int, error),
) error {
	_16_19 := mfile116_3()
	_ = _16_19 // possibly unused.
	_17_15 := mfile117_3()
	_ = _17_15 // possibly unused.
	_19_4 := mfile118_3()
	_ = _19_4 // possibly unused.

	emitter := cff.NopEmitter()

	var (
		flowInfo = &cff.FlowInfo{
			File:   "go.uber.org/cff/internal/tests/modifier/collision/file1.go",
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

	sched := cff.NewScheduler(
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

	// go.uber.org/cff/internal/tests/modifier/collision/file1.go:19:4
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
				err = &cff.PanicError{
					Value:      recovered,
					Stacktrace: debug.Stack(),
				}
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

func _cffConcurrencyfile1_16_3(c int) func() int {
	return func() int { return c }
}

func _cffResultsfile1_17_3(mfile117_15 *int) func() *int {
	return func() *int { return mfile117_15 }
}

func _cffTaskfile1_18_3(mfile119_4 func() (int, error)) func() func() (int, error) {
	return func() func() (int, error) { return mfile119_4 }
}
