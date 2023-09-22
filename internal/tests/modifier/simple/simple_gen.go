//go:build !cff
// +build !cff

package simple

import (
	"context"
	"runtime/debug"
	"time"

	"go.uber.org/cff"
	"go.uber.org/cff/internal/tests/modifier/external"
)

type bar struct{ i int64 }

// Flow is a very simple flow with some inputs and outputs.
func Flow() (int, string, error) {
	var (
		iRes int
		sRes string
	)
	err := _cffFlowsimple_21_9(context.Background(),
		_cffConcurrencysimple_22_3(2),
		_cffResultssimple_23_3(&iRes, &sRes),
		_cffTasksimple_24_3(
			func() int64 {
				return int64(1)
			},
		),
		_cffTasksimple_29_3(
			func(i int64) (*bar, error) {
				return &bar{i}, nil
			}),
		_cffTasksimple_33_3(
			func(*bar) (int, error) {
				return 1, nil
			},
		),
		_cffTasksimple_38_3(
			func(i int) (string, error) {
				if i != 0 {
					return "non-zero", nil
				}
				return "zero", nil
			},
		),
	)
	return iRes, sRes, err
}

// ModifyVarInScope is a simple flow that has a side effect of modifying a variable
// in scope.
func ModifyVarInScope() (bool, []int, error) {
	var res bool
	slc := make([]int, 3)
	err := _cffFlowsimple_55_9(context.Background(),
		_cffConcurrencysimple_56_3(2),
		_cffResultssimple_57_3(&res),
		_cffTasksimple_58_3(
			func() int64 {
				slc[0] = 1
				return int64(1)
			},
		),
		_cffTasksimple_64_3(
			func(i int64) (*bar, error) {
				slc[1] = 2
				return &bar{i}, nil
			}),
		_cffTasksimple_69_3(
			func(*bar) (bool, error) {
				slc[2] = 3
				return true, nil
			},
		),
	)
	return res, slc, err
}

// External is a simple flow that depends on an external package.
func External() (bool, error) {
	var res bool
	err := _cffFlowsimple_82_9(context.Background(),
		_cffConcurrencysimple_83_3(2),
		_cffResultssimple_84_3(&res),
		_cffTasksimple_85_3(
			func() external.A {
				return 1
			},
		),
		_cffTasksimple_90_3(external.Run),
		_cffTasksimple_91_3(
			func(b external.B) (bool, error) {
				return bool(b), nil
			},
		),
	)
	return res, err
}

// Params is a simple cff.Flow that depends on cff.Params.
func Params() (string, external.A, error) {
	var (
		res1 string
		res2 external.A
	)
	err := _cffFlowsimple_106_9(context.Background(),
		_cffConcurrencysimple_107_3(2),
		_cffParamssimple_108_3(1, true),
		_cffResultssimple_109_3(&res1, &res2),
		_cffTasksimple_110_3(
			func(i int) int64 {
				return int64(i)
			},
		),
		_cffTasksimple_115_3(
			func(i int64) (external.A, error) {
				return external.A(i), nil
			}),
		_cffTasksimple_119_3(
			func(b bool) (string, error) {
				if b {
					return "true", nil
				}
				return "false", nil
			},
		),
	)
	return res1, res2, err
}
func _cffFlowsimple_21_9(
	ctx context.Context,
	msimple22_3 func() int,
	msimple23_3 func() (*int, *string),
	msimple24_3 func() func() int64,
	msimple29_3 func() func(i int64) (*bar, error),
	msimple33_3 func() func(*bar) (int, error),
	msimple38_3 func() func(i int) (string, error),
) error {
	_22_19 := msimple22_3()
	_ = _22_19 // possibly unused.
	_23_15, _23_22 := msimple23_3()
	_, _ = _23_15, _23_22 // possibly unused.
	_25_4 := msimple24_3()
	_ = _25_4 // possibly unused.
	_30_4 := msimple29_3()
	_ = _30_4 // possibly unused.
	_34_4 := msimple33_3()
	_ = _34_4 // possibly unused.
	_39_4 := msimple38_3()
	_ = _39_4 // possibly unused.

	emitter := cff.NopEmitter()

	var (
		flowInfo = &cff.FlowInfo{
			File:   "go.uber.org/cff/internal/tests/modifier/simple/simple.go",
			Line:   21,
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
			Concurrency: _22_19, Emitter: schedEmitter,
		},
	)

	var tasks []*struct {
		emitter cff.TaskEmitter
		ran     cff.AtomicBool
		run     func(context.Context) error
		job     *cff.ScheduledJob
	}

	// go.uber.org/cff/internal/tests/modifier/simple/simple.go:25:4
	var (
		v1 int64
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

		v1 = _25_4()
		return
	}

	task0.job = sched.Enqueue(ctx, cff.Job{
		Run: task0.run,
	})

	tasks = append(tasks, task0)

	// go.uber.org/cff/internal/tests/modifier/simple/simple.go:30:4
	var (
		v2 *bar
	)
	task1 := new(struct {
		emitter cff.TaskEmitter
		ran     cff.AtomicBool
		run     func(context.Context) error
		job     *cff.ScheduledJob
	})

	task1.run = func(ctx context.Context) (err error) {
		defer func() {
			recovered := recover()
			if recovered != nil {
				err = &cff.PanicError{
					Value:      recovered,
					Stacktrace: debug.Stack(),
				}
			}
		}()

		v2, err = _30_4(v1)
		return
	}

	task1.job = sched.Enqueue(ctx, cff.Job{
		Run: task1.run,
		Dependencies: []*cff.ScheduledJob{
			task0.job,
		},
	})

	tasks = append(tasks, task1)

	// go.uber.org/cff/internal/tests/modifier/simple/simple.go:34:4
	var (
		v3 int
	)
	task2 := new(struct {
		emitter cff.TaskEmitter
		ran     cff.AtomicBool
		run     func(context.Context) error
		job     *cff.ScheduledJob
	})

	task2.run = func(ctx context.Context) (err error) {
		defer func() {
			recovered := recover()
			if recovered != nil {
				err = &cff.PanicError{
					Value:      recovered,
					Stacktrace: debug.Stack(),
				}
			}
		}()

		v3, err = _34_4(v2)
		return
	}

	task2.job = sched.Enqueue(ctx, cff.Job{
		Run: task2.run,
		Dependencies: []*cff.ScheduledJob{
			task1.job,
		},
	})

	tasks = append(tasks, task2)

	// go.uber.org/cff/internal/tests/modifier/simple/simple.go:39:4
	var (
		v4 string
	)
	task3 := new(struct {
		emitter cff.TaskEmitter
		ran     cff.AtomicBool
		run     func(context.Context) error
		job     *cff.ScheduledJob
	})

	task3.run = func(ctx context.Context) (err error) {
		defer func() {
			recovered := recover()
			if recovered != nil {
				err = &cff.PanicError{
					Value:      recovered,
					Stacktrace: debug.Stack(),
				}
			}
		}()

		v4, err = _39_4(v3)
		return
	}

	task3.job = sched.Enqueue(ctx, cff.Job{
		Run: task3.run,
		Dependencies: []*cff.ScheduledJob{
			task2.job,
		},
	})

	tasks = append(tasks, task3)

	if err := sched.Wait(ctx); err != nil {
		flowEmitter.FlowError(ctx, err)
		return err
	}

	*(_23_15) = v3 // int

	*(_23_22) = v4 // string

	flowEmitter.FlowSuccess(ctx)
	return nil
}

func _cffConcurrencysimple_22_3(c int) func() int {
	return func() int { return c }
}

func _cffResultssimple_23_3(msimple23_15 *int, msimple23_22 *string) func() (*int, *string) {
	return func() (*int, *string) { return msimple23_15, msimple23_22 }
}

func _cffTasksimple_24_3(msimple25_4 func() int64) func() func() int64 {
	return func() func() int64 { return msimple25_4 }
}

func _cffTasksimple_29_3(msimple30_4 func(i int64) (*bar, error)) func() func(i int64) (*bar, error) {
	return func() func(i int64) (*bar, error) { return msimple30_4 }
}

func _cffTasksimple_33_3(msimple34_4 func(*bar) (int, error)) func() func(*bar) (int, error) {
	return func() func(*bar) (int, error) { return msimple34_4 }
}

func _cffTasksimple_38_3(msimple39_4 func(i int) (string, error)) func() func(i int) (string, error) {
	return func() func(i int) (string, error) { return msimple39_4 }
}

func _cffFlowsimple_55_9(
	ctx context.Context,
	msimple56_3 func() int,
	msimple57_3 func() *bool,
	msimple58_3 func() func() int64,
	msimple64_3 func() func(i int64) (*bar, error),
	msimple69_3 func() func(*bar) (bool, error),
) error {
	_56_19 := msimple56_3()
	_ = _56_19 // possibly unused.
	_57_15 := msimple57_3()
	_ = _57_15 // possibly unused.
	_59_4 := msimple58_3()
	_ = _59_4 // possibly unused.
	_65_4 := msimple64_3()
	_ = _65_4 // possibly unused.
	_70_4 := msimple69_3()
	_ = _70_4 // possibly unused.

	emitter := cff.NopEmitter()

	var (
		flowInfo = &cff.FlowInfo{
			File:   "go.uber.org/cff/internal/tests/modifier/simple/simple.go",
			Line:   55,
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
			Concurrency: _56_19, Emitter: schedEmitter,
		},
	)

	var tasks []*struct {
		emitter cff.TaskEmitter
		ran     cff.AtomicBool
		run     func(context.Context) error
		job     *cff.ScheduledJob
	}

	// go.uber.org/cff/internal/tests/modifier/simple/simple.go:59:4
	var (
		v1 int64
	)
	task4 := new(struct {
		emitter cff.TaskEmitter
		ran     cff.AtomicBool
		run     func(context.Context) error
		job     *cff.ScheduledJob
	})

	task4.run = func(ctx context.Context) (err error) {
		defer func() {
			recovered := recover()
			if recovered != nil {
				err = &cff.PanicError{
					Value:      recovered,
					Stacktrace: debug.Stack(),
				}
			}
		}()

		v1 = _59_4()
		return
	}

	task4.job = sched.Enqueue(ctx, cff.Job{
		Run: task4.run,
	})

	tasks = append(tasks, task4)

	// go.uber.org/cff/internal/tests/modifier/simple/simple.go:65:4
	var (
		v2 *bar
	)
	task5 := new(struct {
		emitter cff.TaskEmitter
		ran     cff.AtomicBool
		run     func(context.Context) error
		job     *cff.ScheduledJob
	})

	task5.run = func(ctx context.Context) (err error) {
		defer func() {
			recovered := recover()
			if recovered != nil {
				err = &cff.PanicError{
					Value:      recovered,
					Stacktrace: debug.Stack(),
				}
			}
		}()

		v2, err = _65_4(v1)
		return
	}

	task5.job = sched.Enqueue(ctx, cff.Job{
		Run: task5.run,
		Dependencies: []*cff.ScheduledJob{
			task4.job,
		},
	})

	tasks = append(tasks, task5)

	// go.uber.org/cff/internal/tests/modifier/simple/simple.go:70:4
	var (
		v5 bool
	)
	task6 := new(struct {
		emitter cff.TaskEmitter
		ran     cff.AtomicBool
		run     func(context.Context) error
		job     *cff.ScheduledJob
	})

	task6.run = func(ctx context.Context) (err error) {
		defer func() {
			recovered := recover()
			if recovered != nil {
				err = &cff.PanicError{
					Value:      recovered,
					Stacktrace: debug.Stack(),
				}
			}
		}()

		v5, err = _70_4(v2)
		return
	}

	task6.job = sched.Enqueue(ctx, cff.Job{
		Run: task6.run,
		Dependencies: []*cff.ScheduledJob{
			task5.job,
		},
	})

	tasks = append(tasks, task6)

	if err := sched.Wait(ctx); err != nil {
		flowEmitter.FlowError(ctx, err)
		return err
	}

	*(_57_15) = v5 // bool

	flowEmitter.FlowSuccess(ctx)
	return nil
}

func _cffConcurrencysimple_56_3(c int) func() int {
	return func() int { return c }
}

func _cffResultssimple_57_3(msimple57_15 *bool) func() *bool {
	return func() *bool { return msimple57_15 }
}

func _cffTasksimple_58_3(msimple59_4 func() int64) func() func() int64 {
	return func() func() int64 { return msimple59_4 }
}

func _cffTasksimple_64_3(msimple65_4 func(i int64) (*bar, error)) func() func(i int64) (*bar, error) {
	return func() func(i int64) (*bar, error) { return msimple65_4 }
}

func _cffTasksimple_69_3(msimple70_4 func(*bar) (bool, error)) func() func(*bar) (bool, error) {
	return func() func(*bar) (bool, error) { return msimple70_4 }
}

func _cffFlowsimple_82_9(
	ctx context.Context,
	msimple83_3 func() int,
	msimple84_3 func() *bool,
	msimple85_3 func() func() external.A,
	msimple90_3 func() func(a external.A) external.B,
	msimple91_3 func() func(b external.B) (bool, error),
) error {
	_83_19 := msimple83_3()
	_ = _83_19 // possibly unused.
	_84_15 := msimple84_3()
	_ = _84_15 // possibly unused.
	_86_4 := msimple85_3()
	_ = _86_4 // possibly unused.
	_90_12 := msimple90_3()
	_ = _90_12 // possibly unused.
	_92_4 := msimple91_3()
	_ = _92_4 // possibly unused.

	emitter := cff.NopEmitter()

	var (
		flowInfo = &cff.FlowInfo{
			File:   "go.uber.org/cff/internal/tests/modifier/simple/simple.go",
			Line:   82,
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
			Concurrency: _83_19, Emitter: schedEmitter,
		},
	)

	var tasks []*struct {
		emitter cff.TaskEmitter
		ran     cff.AtomicBool
		run     func(context.Context) error
		job     *cff.ScheduledJob
	}

	// go.uber.org/cff/internal/tests/modifier/simple/simple.go:86:4
	var (
		v6 external.A
	)
	task7 := new(struct {
		emitter cff.TaskEmitter
		ran     cff.AtomicBool
		run     func(context.Context) error
		job     *cff.ScheduledJob
	})

	task7.run = func(ctx context.Context) (err error) {
		defer func() {
			recovered := recover()
			if recovered != nil {
				err = &cff.PanicError{
					Value:      recovered,
					Stacktrace: debug.Stack(),
				}
			}
		}()

		v6 = _86_4()
		return
	}

	task7.job = sched.Enqueue(ctx, cff.Job{
		Run: task7.run,
	})

	tasks = append(tasks, task7)

	// go.uber.org/cff/internal/tests/modifier/simple/simple.go:90:12
	var (
		v7 external.B
	)
	task8 := new(struct {
		emitter cff.TaskEmitter
		ran     cff.AtomicBool
		run     func(context.Context) error
		job     *cff.ScheduledJob
	})

	task8.run = func(ctx context.Context) (err error) {
		defer func() {
			recovered := recover()
			if recovered != nil {
				err = &cff.PanicError{
					Value:      recovered,
					Stacktrace: debug.Stack(),
				}
			}
		}()

		v7 = _90_12(v6)
		return
	}

	task8.job = sched.Enqueue(ctx, cff.Job{
		Run: task8.run,
		Dependencies: []*cff.ScheduledJob{
			task7.job,
		},
	})

	tasks = append(tasks, task8)

	// go.uber.org/cff/internal/tests/modifier/simple/simple.go:92:4
	var (
		v5 bool
	)
	task9 := new(struct {
		emitter cff.TaskEmitter
		ran     cff.AtomicBool
		run     func(context.Context) error
		job     *cff.ScheduledJob
	})

	task9.run = func(ctx context.Context) (err error) {
		defer func() {
			recovered := recover()
			if recovered != nil {
				err = &cff.PanicError{
					Value:      recovered,
					Stacktrace: debug.Stack(),
				}
			}
		}()

		v5, err = _92_4(v7)
		return
	}

	task9.job = sched.Enqueue(ctx, cff.Job{
		Run: task9.run,
		Dependencies: []*cff.ScheduledJob{
			task8.job,
		},
	})

	tasks = append(tasks, task9)

	if err := sched.Wait(ctx); err != nil {
		flowEmitter.FlowError(ctx, err)
		return err
	}

	*(_84_15) = v5 // bool

	flowEmitter.FlowSuccess(ctx)
	return nil
}

func _cffConcurrencysimple_83_3(c int) func() int {
	return func() int { return c }
}

func _cffResultssimple_84_3(msimple84_15 *bool) func() *bool {
	return func() *bool { return msimple84_15 }
}

func _cffTasksimple_85_3(msimple86_4 func() external.A) func() func() external.A {
	return func() func() external.A { return msimple86_4 }
}

func _cffTasksimple_90_3(msimple90_12 func(a external.A) external.B) func() func(a external.A) external.B {
	return func() func(a external.A) external.B { return msimple90_12 }
}

func _cffTasksimple_91_3(msimple92_4 func(b external.B) (bool, error)) func() func(b external.B) (bool, error) {
	return func() func(b external.B) (bool, error) { return msimple92_4 }
}

func _cffFlowsimple_106_9(
	ctx context.Context,
	msimple107_3 func() int,
	msimple108_3 func() (int, bool),
	msimple109_3 func() (*string, *external.A),
	msimple110_3 func() func(i int) int64,
	msimple115_3 func() func(i int64) (external.A, error),
	msimple119_3 func() func(b bool) (string, error),
) error {
	_107_19 := msimple107_3()
	_ = _107_19 // possibly unused.
	_108_14, _108_17 := msimple108_3()
	_, _ = _108_14, _108_17 // possibly unused.
	_109_15, _109_22 := msimple109_3()
	_, _ = _109_15, _109_22 // possibly unused.
	_111_4 := msimple110_3()
	_ = _111_4 // possibly unused.
	_116_4 := msimple115_3()
	_ = _116_4 // possibly unused.
	_120_4 := msimple119_3()
	_ = _120_4 // possibly unused.

	var v3 int = _108_14
	var v5 bool = _108_17
	emitter := cff.NopEmitter()

	var (
		flowInfo = &cff.FlowInfo{
			File:   "go.uber.org/cff/internal/tests/modifier/simple/simple.go",
			Line:   106,
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
			Concurrency: _107_19, Emitter: schedEmitter,
		},
	)

	var tasks []*struct {
		emitter cff.TaskEmitter
		ran     cff.AtomicBool
		run     func(context.Context) error
		job     *cff.ScheduledJob
	}

	// go.uber.org/cff/internal/tests/modifier/simple/simple.go:111:4
	var (
		v1 int64
	)
	task10 := new(struct {
		emitter cff.TaskEmitter
		ran     cff.AtomicBool
		run     func(context.Context) error
		job     *cff.ScheduledJob
	})

	task10.run = func(ctx context.Context) (err error) {
		defer func() {
			recovered := recover()
			if recovered != nil {
				err = &cff.PanicError{
					Value:      recovered,
					Stacktrace: debug.Stack(),
				}
			}
		}()

		v1 = _111_4(v3)
		return
	}

	task10.job = sched.Enqueue(ctx, cff.Job{
		Run: task10.run,
	})

	tasks = append(tasks, task10)

	// go.uber.org/cff/internal/tests/modifier/simple/simple.go:116:4
	var (
		v6 external.A
	)
	task11 := new(struct {
		emitter cff.TaskEmitter
		ran     cff.AtomicBool
		run     func(context.Context) error
		job     *cff.ScheduledJob
	})

	task11.run = func(ctx context.Context) (err error) {
		defer func() {
			recovered := recover()
			if recovered != nil {
				err = &cff.PanicError{
					Value:      recovered,
					Stacktrace: debug.Stack(),
				}
			}
		}()

		v6, err = _116_4(v1)
		return
	}

	task11.job = sched.Enqueue(ctx, cff.Job{
		Run: task11.run,
		Dependencies: []*cff.ScheduledJob{
			task10.job,
		},
	})

	tasks = append(tasks, task11)

	// go.uber.org/cff/internal/tests/modifier/simple/simple.go:120:4
	var (
		v4 string
	)
	task12 := new(struct {
		emitter cff.TaskEmitter
		ran     cff.AtomicBool
		run     func(context.Context) error
		job     *cff.ScheduledJob
	})

	task12.run = func(ctx context.Context) (err error) {
		defer func() {
			recovered := recover()
			if recovered != nil {
				err = &cff.PanicError{
					Value:      recovered,
					Stacktrace: debug.Stack(),
				}
			}
		}()

		v4, err = _120_4(v5)
		return
	}

	task12.job = sched.Enqueue(ctx, cff.Job{
		Run: task12.run,
	})

	tasks = append(tasks, task12)

	if err := sched.Wait(ctx); err != nil {
		flowEmitter.FlowError(ctx, err)
		return err
	}

	*(_109_15) = v4 // string

	*(_109_22) = v6 // go.uber.org/cff/internal/tests/modifier/external.A

	flowEmitter.FlowSuccess(ctx)
	return nil
}

func _cffConcurrencysimple_107_3(c int) func() int {
	return func() int { return c }
}

func _cffParamssimple_108_3(msimple108_14 int, msimple108_17 bool) func() (int, bool) {
	return func() (int, bool) { return msimple108_14, msimple108_17 }
}

func _cffResultssimple_109_3(msimple109_15 *string, msimple109_22 *external.A) func() (*string, *external.A) {
	return func() (*string, *external.A) { return msimple109_15, msimple109_22 }
}

func _cffTasksimple_110_3(msimple111_4 func(i int) int64) func() func(i int) int64 {
	return func() func(i int) int64 { return msimple111_4 }
}

func _cffTasksimple_115_3(msimple116_4 func(i int64) (external.A, error)) func() func(i int64) (external.A, error) {
	return func() func(i int64) (external.A, error) { return msimple116_4 }
}

func _cffTasksimple_119_3(msimple120_4 func(b bool) (string, error)) func() func(b bool) (string, error) {
	return func() func(b bool) (string, error) { return msimple120_4 }
}
