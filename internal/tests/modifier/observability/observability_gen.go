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
	err := _cffFlowobservability_17_9(context.Background(),
		_cffConcurrencyobservability_18_3(2),
		_cffParamsobservability_19_3(1),
		_cffResultsobservability_20_3(&res),
		_cffWithEmitterobservability_21_3(cff.TallyEmitter(scope)),
		_cffWithEmitterobservability_22_3(cff.LogEmitter(logger)),
		_cffInstrumentFlowobservability_23_3("Instrumented"),
		_cffTaskobservability_24_3(
			func(i int) int64 {
				return int64(1)
			},
		),
	)
	return res, err
}
func _cffFlowobservability_17_9(
	ctx context.Context,
	mobservability18_3 func() int,
	mobservability19_3 func() int,
	mobservability20_3 func() *int64,
	mobservability21_3 func() cff.Emitter,
	mobservability22_3 func() cff.Emitter,
	mobservability23_3 func() string,
	mobservability24_3 func() func(i int) int64,
) error {
	_18_19 := mobservability18_3()
	_ = _18_19 // possibly unused.
	_19_14 := mobservability19_3()
	_ = _19_14 // possibly unused.
	_20_15 := mobservability20_3()
	_ = _20_15 // possibly unused.
	_21_19 := mobservability21_3()
	_ = _21_19 // possibly unused.
	_22_19 := mobservability22_3()
	_ = _22_19 // possibly unused.
	_23_22 := mobservability23_3()
	_ = _23_22 // possibly unused.
	_25_4 := mobservability24_3()
	_ = _25_4 // possibly unused.

	var v1 int = _19_14
	emitter := cff.EmitterStack(_21_19, _22_19)

	var (
		flowInfo = &cff.FlowInfo{
			Name:   _23_22,
			File:   "go.uber.org/cff/internal/tests/modifier/observability/observability.go",
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

	// go.uber.org/cff/internal/tests/modifier/observability/observability.go:25:4
	var (
		v2 int64
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

		v2 = _25_4(v1)
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
}

func _cffConcurrencyobservability_18_3(c int) func() int {
	return func() int { return c }
}

func _cffParamsobservability_19_3(mobservability19_14 int) func() int {
	return func() int { return mobservability19_14 }
}

func _cffResultsobservability_20_3(mobservability20_15 *int64) func() *int64 {
	return func() *int64 { return mobservability20_15 }
}

func _cffWithEmitterobservability_21_3(mobservability21_19 cff.Emitter) func() cff.Emitter {
	return func() cff.Emitter { return mobservability21_19 }
}

func _cffWithEmitterobservability_22_3(mobservability22_19 cff.Emitter) func() cff.Emitter {
	return func() cff.Emitter { return mobservability22_19 }
}

func _cffInstrumentFlowobservability_23_3(mobservability23_22 string) func() string {
	return func() string { return mobservability23_22 }
}

func _cffTaskobservability_24_3(mobservability25_4 func(i int) int64) func() func(i int) int64 {
	return func() func(i int) int64 { return mobservability25_4 }
}
