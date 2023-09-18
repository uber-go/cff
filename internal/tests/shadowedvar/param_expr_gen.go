//go:build !cff
// +build !cff

package shadowedvar

import (
	"context"
	"runtime/debug"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/cff"
)

// ParamOrder initializes a cff.Flow to test that the order in which user
// provided expressions are evalauted matches the order in which they were
// provided to the cff.Flow.
func ParamOrder(track *orderCheck) error {
	var res string
	return func() (err error) {

		_21_3 := track.ctx()

		_22_14 := track.param1()

		_22_30 := track.param2()

		_23_15 := &res

		_24_12 := func(_ int, _ bool) (string, error) {
			return "", nil
		}
		ctx := _21_3
		var v1 int = _22_14
		var v2 bool = _22_30
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/shadowedvar/param_expr.go",
				Line:   20,
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

		// go.uber.org/cff/internal/tests/shadowedvar/param_expr.go:24:12
		var (
			v3 string
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

			v3, err = _24_12(v1, v2)

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

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_23_15) = v3 // string

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
}

// NilParam verifies that cff.Flow is compilable with a user provided nil.
// cff should compile and generate this flow even if no test function
// uses it.
func NilParam() {
	var res []int
	func() (err error) {

		_36_3 := context.Background()

		_37_14 := 1

		_37_17 := true

		_38_15 := &res

		_40_4 := func(_ int, _ bool) ([]int, error) {
			return nil, nil
		}
		ctx := _36_3
		var v1 int = _37_14
		var v2 bool = _37_17
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/internal/tests/shadowedvar/param_expr.go",
				Line:   35,
				Column: 2,
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

		// go.uber.org/cff/internal/tests/shadowedvar/param_expr.go:40:4
		var (
			v4 []int
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
					taskEmitter.TaskPanicRecovered(ctx, recovered)
					v4, err = nil, nil
				}
			}()

			defer task1.ran.Store(true)

			v4, err = _40_4(v1, v2)

			if err != nil {
				taskEmitter.TaskErrorRecovered(ctx, err)
				v4, err = nil, nil
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task1.job = sched.Enqueue(ctx, cff.Job{
			Run: task1.run,
		})
		tasks = append(tasks, task1)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_38_15) = v4 // []int

		flowEmitter.FlowSuccess(ctx)
		return nil
	}()
}

// checks the call order of parameter expressions invoked by cff.Flow.
// cff.Flow assumes all non-task parameters expressions are not concurrently
// invoked.
type orderCheck struct {
	t *testing.T
	// counter tracks order in which expressions are invoked.
	counter int
	// order expectations for parameter expressions.
	order map[string]int
}

func (c *orderCheck) ctx() context.Context {
	o, ok := c.order["ctx"]
	require.True(c.t, ok)
	assert.Equal(c.t, o, c.counter)

	c.counter++
	return context.Background()
}

func (c *orderCheck) param1() int {
	o, ok := c.order["param1"]
	require.True(c.t, ok)
	assert.Equal(c.t, o, c.counter)

	c.counter++
	return 0
}

func (c *orderCheck) param2() bool {
	o, ok := c.order["param2"]
	require.True(c.t, ok)
	assert.Equal(c.t, o, c.counter)

	c.counter++
	return true
}
