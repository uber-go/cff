package tests

import (
	"context"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uber-go/tally"
	"go.uber.org/cff"
	"go.uber.org/cff/internal/emittertest"
	"go.uber.org/cff/internal/tests/benchmark"
	"go.uber.org/cff/internal/tests/instrument"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

// TestInstrumentFlowEmitter verifies that new Emitter interface gets called if
// it's passed in.
func TestInstrumentFlowEmitter(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	ctx := context.Background()

	emitter := emittertest.NewMockEmitter(mockCtrl)

	taskEmitter := emittertest.NewMockTaskEmitter(mockCtrl)
	flowEmitter := emittertest.NewMockFlowEmitter(mockCtrl)
	schedEmitter := emittertest.NewMockSchedulerEmitter(mockCtrl)
	schedEmitter.EXPECT().EmitScheduler(gomock.Any()).AnyTimes()

	flowsucc := flowEmitter.EXPECT().FlowSuccess(ctx)
	flowEmitter.EXPECT().FlowDone(ctx, gomock.Any()).After(flowsucc)

	// 2 tasks.
	taskEmitter.EXPECT().TaskSuccess(ctx).Times(2)
	taskEmitter.EXPECT().TaskDone(ctx, gomock.Any()).Times(2)

	flowInfo := &cff.FlowInfo{
		Name:   "AtoiRun",
		File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
		Line:   378,
		Column: 8,
	}

	emitter.EXPECT().FlowInit(flowInfo).Return(flowEmitter)

	emitter.EXPECT().SchedulerInit(
		&cff.SchedulerInfo{
			Name:      flowInfo.Name,
			Directive: cff.FlowDirective,
			File:      flowInfo.File,
			Line:      flowInfo.Line,
			Column:    flowInfo.Column,
		}).Return(schedEmitter)

	// 2 in the tasks for loop inside defer() and twice after.
	emitter.EXPECT().TaskInit(gomock.Any(), gomock.Any()).Times(2).Return(taskEmitter)

	scope := tally.NewTestScope("", nil)
	// Logging
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	// CustomEmitter Unit
	g := &instrument.CustomEmitter{
		Logger:  logger,
		Scope:   scope,
		Emitter: emitter,
	}
	v, err := g.RunFlow(ctx, "1")

	assert.NoError(t, err)
	assert.Equal(t, uint8(1), v)
}

// TestInstrumentParallelEmitter verifies that new Emitter interface gets called if
// it's passed in.
func TestInstrumentParallelEmitter(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	ctx := context.Background()

	emitter := emittertest.NewMockEmitter(mockCtrl)

	parallelEmitter := emittertest.NewMockParallelEmitter(mockCtrl)
	taskEmitter := emittertest.NewMockTaskEmitter(mockCtrl)
	schedEmitter := emittertest.NewMockSchedulerEmitter(mockCtrl)
	schedEmitter.EXPECT().EmitScheduler(gomock.Any()).AnyTimes()

	tasksucc := taskEmitter.EXPECT().TaskSuccess(ctx)
	taskEmitter.EXPECT().TaskDone(ctx, gomock.Any()).After(tasksucc)

	parallelsucc := parallelEmitter.EXPECT().ParallelSuccess(ctx)
	parallelEmitter.EXPECT().ParallelDone(ctx, gomock.Any()).After(parallelsucc)

	parallelInfo := &cff.ParallelInfo{
		Name:   "RunParallelTasksAndTask",
		File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
		Line:   409,
		Column: 9,
	}

	emitter.EXPECT().ParallelInit(parallelInfo).Return(parallelEmitter)
	emitter.EXPECT().TaskInit(gomock.Any(), gomock.Any()).Times(1).Return(taskEmitter)

	emitter.EXPECT().SchedulerInit(
		&cff.SchedulerInfo{
			Name:      parallelInfo.Name,
			Directive: cff.ParallelDirective,
			File:      parallelInfo.File,
			Line:      parallelInfo.Line,
			Column:    parallelInfo.Column,
		}).Return(schedEmitter)

	scope := tally.NewTestScope("", nil)
	// Logging
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	// CustomEmitter Unit
	g := &instrument.CustomEmitter{
		Logger:  logger,
		Scope:   scope,
		Emitter: emitter,
	}
	assert.NoError(t, g.RunParallelTasksAndTask(ctx, "1"))
}

func TestInstrumentFlowErrorME(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	ctx := context.Background()

	emitter := emittertest.NewMockEmitter(mockCtrl)

	taskEmitter := emittertest.NewMockTaskEmitter(mockCtrl)
	flowEmitter := emittertest.NewMockFlowEmitter(mockCtrl)
	schedEmitter := emittertest.NewMockSchedulerEmitter(mockCtrl)
	schedEmitter.EXPECT().EmitScheduler(gomock.Any()).AnyTimes()

	// flowFailedEmitter := emittertest.NewMockFlowEmitter(mockCtrl)

	flowEmitter.EXPECT().FlowError(ctx, gomock.Any())
	flowEmitter.EXPECT().FlowDone(ctx, gomock.Any())

	// 2 tasks.
	taskEmitter.EXPECT().TaskError(ctx, gomock.Any())
	taskEmitter.EXPECT().TaskSkipped(ctx, gomock.Any())
	taskEmitter.EXPECT().TaskDone(ctx, gomock.Any())

	emitter.EXPECT().FlowInit(gomock.Any()).Return(flowEmitter)
	emitter.EXPECT().TaskInit(gomock.Any(), gomock.Any()).Times(2).Return(taskEmitter)
	emitter.EXPECT().SchedulerInit(gomock.Any()).AnyTimes().Return(schedEmitter)

	scope := tally.NewTestScope("", nil)
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	h := &instrument.CustomEmitter{
		Scope:   scope,
		Logger:  logger,
		Emitter: emitter,
	}
	_, err := h.RunFlow(ctx, "NaN")

	assert.Error(t, err)
}

func TestInstrumentParallelErrorME(t *testing.T) {
	t.Run("cff.Tasks error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ctx := context.Background()

		emitter := emittertest.NewMockEmitter(mockCtrl)

		parallelEmitter := emittertest.NewMockParallelEmitter(mockCtrl)
		schedEmitter := emittertest.NewMockSchedulerEmitter(mockCtrl)
		schedEmitter.EXPECT().EmitScheduler(gomock.Any()).AnyTimes()

		parallelEmitter.EXPECT().ParallelError(ctx, gomock.Any())
		parallelEmitter.EXPECT().ParallelDone(ctx, gomock.Any())

		emitter.EXPECT().ParallelInit(gomock.Any()).Return(parallelEmitter)
		emitter.EXPECT().SchedulerInit(gomock.Any()).AnyTimes().Return(schedEmitter)

		scope := tally.NewTestScope("", nil)
		core, _ := observer.New(zap.DebugLevel)
		logger := zap.New(core)

		h := &instrument.CustomEmitter{
			Scope:   scope,
			Logger:  logger,
			Emitter: emitter,
		}
		assert.Error(t, h.RunParallelTasks(ctx, "NaN"))
	})

	t.Run("cff.Task error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ctx := context.Background()

		emitter := emittertest.NewMockEmitter(mockCtrl)

		parallelEmitter := emittertest.NewMockParallelEmitter(mockCtrl)
		taskEmitter := emittertest.NewMockTaskEmitter(mockCtrl)
		schedEmitter := emittertest.NewMockSchedulerEmitter(mockCtrl)
		schedEmitter.EXPECT().EmitScheduler(gomock.Any()).AnyTimes()

		taskEmitter.EXPECT().TaskError(ctx, gomock.Any())
		taskEmitter.EXPECT().TaskDone(ctx, gomock.Any())

		parallelEmitter.EXPECT().ParallelError(ctx, gomock.Any())
		parallelEmitter.EXPECT().ParallelDone(ctx, gomock.Any())

		emitter.EXPECT().ParallelInit(gomock.Any()).Return(parallelEmitter)
		emitter.EXPECT().SchedulerInit(gomock.Any()).AnyTimes().Return(schedEmitter)
		emitter.EXPECT().TaskInit(gomock.Any(), gomock.Any()).Times(1).Return(taskEmitter)

		scope := tally.NewTestScope("", nil)
		core, _ := observer.New(zap.DebugLevel)
		logger := zap.New(core)

		h := &instrument.CustomEmitter{
			Scope:   scope,
			Logger:  logger,
			Emitter: emitter,
		}
		assert.Error(t, h.RunParallelTask(ctx, "NaN"))
	})
}

func TestInstrumentTaskButNotFlowME(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	ctx := context.Background()

	emitter := emittertest.NewMockEmitter(mockCtrl)

	taskEmitter := emittertest.NewMockTaskEmitter(mockCtrl)
	schedEmitter := emittertest.NewMockSchedulerEmitter(mockCtrl)
	schedEmitter.EXPECT().EmitScheduler(gomock.Any()).AnyTimes()

	taskEmitter.EXPECT().TaskSuccess(ctx)
	taskEmitter.EXPECT().TaskDone(ctx, gomock.Any())
	emitter.EXPECT().TaskInit(gomock.Any(), gomock.Any()).Return(taskEmitter)
	emitter.EXPECT().SchedulerInit(gomock.Any()).AnyTimes().Return(schedEmitter)

	scope := tally.NewTestScope("", nil)
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	g := &instrument.CustomEmitter{
		Scope:   scope,
		Logger:  logger,
		Emitter: emitter,
	}
	v, err := g.FlowOnlyInstrumentTask(ctx, "1")

	assert.NoError(t, err)
	assert.Equal(t, 1, v)
}

func TestInstrumentTaskButNotParallelME(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	ctx := context.Background()

	emitter := emittertest.NewMockEmitter(mockCtrl)

	taskEmitter := emittertest.NewMockTaskEmitter(mockCtrl)
	schedEmitter := emittertest.NewMockSchedulerEmitter(mockCtrl)
	schedEmitter.EXPECT().EmitScheduler(gomock.Any()).AnyTimes()

	taskEmitter.EXPECT().TaskSuccess(ctx)
	taskEmitter.EXPECT().TaskDone(ctx, gomock.Any())

	emitter.EXPECT().TaskInit(gomock.Any(), gomock.Any()).Return(taskEmitter)
	emitter.EXPECT().SchedulerInit(gomock.Any()).AnyTimes().Return(schedEmitter)

	scope := tally.NewTestScope("", nil)
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	g := &instrument.CustomEmitter{
		Scope:   scope,
		Logger:  logger,
		Emitter: emitter,
	}
	assert.NoError(t, g.ParallelOnlyInstrumentTask(ctx, "1"))
}

func TestInstrumentFlowCancelledContextME(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	scope := tally.NewTestScope("", nil)
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	flowCancelledErr := ctx.Err()

	emitter := emittertest.NewMockEmitter(mockCtrl)

	taskEmitter := emittertest.NewMockTaskEmitter(mockCtrl)
	flowEmitter := emittertest.NewMockFlowEmitter(mockCtrl)
	schedEmitter := emittertest.NewMockSchedulerEmitter(mockCtrl)
	schedEmitter.EXPECT().EmitScheduler(gomock.Any()).AnyTimes()

	flowEmitter.EXPECT().FlowError(ctx, flowCancelledErr)
	flowEmitter.EXPECT().FlowDone(ctx, gomock.Any())

	taskEmitter.EXPECT().TaskSkipped(ctx, gomock.Any()).Times(2)

	emitter.EXPECT().FlowInit(gomock.Any()).Return(flowEmitter)
	emitter.EXPECT().TaskInit(gomock.Any(), gomock.Any()).AnyTimes().Return(taskEmitter)
	emitter.EXPECT().SchedulerInit(gomock.Any()).AnyTimes().Return(schedEmitter)

	g := &instrument.CustomEmitter{
		Scope:   scope,
		Logger:  logger,
		Emitter: emitter,
	}

	_, err := g.RunFlow(ctx, "1")
	assert.Error(t, err)
}

func TestInstrumentParallelCancelledContextME(t *testing.T) {
	t.Run("cff.Tasks", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)

		scope := tally.NewTestScope("", nil)
		core, _ := observer.New(zap.DebugLevel)
		logger := zap.New(core)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		parallelCancelledError := ctx.Err()

		emitter := emittertest.NewMockEmitter(mockCtrl)

		parallelEmitter := emittertest.NewMockParallelEmitter(mockCtrl)
		schedEmitter := emittertest.NewMockSchedulerEmitter(mockCtrl)
		schedEmitter.EXPECT().EmitScheduler(gomock.Any()).AnyTimes()

		parallelEmitter.EXPECT().ParallelError(ctx, parallelCancelledError)
		parallelEmitter.EXPECT().ParallelDone(ctx, gomock.Any())

		emitter.EXPECT().ParallelInit(gomock.Any()).Return(parallelEmitter)
		emitter.EXPECT().SchedulerInit(gomock.Any()).AnyTimes().Return(schedEmitter)

		g := &instrument.CustomEmitter{
			Scope:   scope,
			Logger:  logger,
			Emitter: emitter,
		}

		assert.Error(t, g.RunParallelTasks(ctx, "1"))
	})

	t.Run("cff.Task", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)

		scope := tally.NewTestScope("", nil)
		core, _ := observer.New(zap.DebugLevel)
		logger := zap.New(core)
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		cancel()

		parallelCancelledError := ctx.Err()

		emitter := emittertest.NewMockEmitter(mockCtrl)

		taskEmitter := emittertest.NewMockTaskEmitter(mockCtrl)
		parallelEmitter := emittertest.NewMockParallelEmitter(mockCtrl)
		schedEmitter := emittertest.NewMockSchedulerEmitter(mockCtrl)
		schedEmitter.EXPECT().EmitScheduler(gomock.Any()).AnyTimes()

		taskEmitter.EXPECT().TaskSkipped(ctx, parallelCancelledError)

		parallelEmitter.EXPECT().ParallelError(ctx, parallelCancelledError)
		parallelEmitter.EXPECT().ParallelDone(ctx, gomock.Any())

		emitter.EXPECT().SchedulerInit(gomock.Any()).AnyTimes().Return(schedEmitter)
		emitter.EXPECT().ParallelInit(gomock.Any()).Return(parallelEmitter)
		emitter.EXPECT().TaskInit(gomock.Any(), gomock.Any()).Times(1).Return(taskEmitter)

		g := &instrument.CustomEmitter{
			Scope:   scope,
			Logger:  logger,
			Emitter: emitter,
		}

		assert.Error(t, g.RunParallelTask(ctx, "1"))
	})
}

func TestInstrumentFlowRecoverME(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	ctx := context.Background()

	scope := tally.NewTestScope("", nil)
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	emitter := emittertest.NewMockEmitter(mockCtrl)

	taskEmitter := emittertest.NewMockTaskEmitter(mockCtrl)
	flowEmitter := emittertest.NewMockFlowEmitter(mockCtrl)
	schedEmitter := emittertest.NewMockSchedulerEmitter(mockCtrl)
	schedEmitter.EXPECT().EmitScheduler(gomock.Any()).AnyTimes()

	flowEmitter.EXPECT().FlowSuccess(ctx)
	flowEmitter.EXPECT().FlowDone(ctx, gomock.Any())

	taskEmitter.EXPECT().TaskSuccess(ctx)
	taskEmitter.EXPECT().TaskErrorRecovered(ctx, gomock.Any())
	taskEmitter.EXPECT().TaskDone(ctx, gomock.Any()).Times(2)

	emitter.EXPECT().FlowInit(&cff.FlowInfo{
		Name:   "AtoiRun",
		File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
		Line:   378,
		Column: 8,
	}).Return(flowEmitter)
	emitter.EXPECT().TaskInit(gomock.Any(), gomock.Any()).Times(2).Return(taskEmitter)
	emitter.EXPECT().SchedulerInit(gomock.Any()).AnyTimes().Return(schedEmitter)

	g := &instrument.CustomEmitter{
		Scope:   scope,
		Logger:  logger,
		Emitter: emitter,
	}

	v, err := g.RunFlow(ctx, "300")

	assert.NoError(t, err)
	assert.Equal(t, uint8(0), v)
}

// TestT3630161ME tests against regression for T3630161
func TestT3630161ME(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	ctx := context.Background()

	emitter := emittertest.NewMockEmitter(mockCtrl)

	taskEmitter := emittertest.NewMockTaskEmitter(mockCtrl)
	flowEmitter := emittertest.NewMockFlowEmitter(mockCtrl)
	schedEmitter := emittertest.NewMockSchedulerEmitter(mockCtrl)
	schedEmitter.EXPECT().EmitScheduler(gomock.Any()).AnyTimes()

	flowEmitter.EXPECT().FlowSuccess(ctx)
	flowEmitter.EXPECT().FlowDone(ctx, gomock.Any())

	// 2 tasks.
	taskEmitter.EXPECT().TaskErrorRecovered(ctx, gomock.Any())
	taskEmitter.EXPECT().TaskDone(ctx, gomock.Any()).Times(2)
	taskEmitter.EXPECT().TaskSuccess(ctx)

	emitter.EXPECT().FlowInit(gomock.Any()).Return(flowEmitter)
	emitter.EXPECT().TaskInit(gomock.Any(), gomock.Any()).Times(2).Return(taskEmitter)
	emitter.EXPECT().SchedulerInit(gomock.Any()).AnyTimes().Return(schedEmitter)

	scope := tally.NewTestScope("", nil)
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	g := &instrument.CustomEmitter{
		Scope:   scope,
		Logger:  logger,
		Emitter: emitter,
	}

	g.T3630161(ctx)
}

// TestT3795761 tests against regression for T3795761 where a task that returns no error is not reported as
// skipped when an earlier task that it depends on returns an error.
func TestT3795761ME(t *testing.T) {
	ctx := context.Background()

	scope := tally.NewTestScope("", nil)
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	t.Run("should run error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)

		emitter := emittertest.NewMockEmitter(mockCtrl)

		taskEmitter := emittertest.NewMockTaskEmitter(mockCtrl)
		flowEmitter := emittertest.NewMockFlowEmitter(mockCtrl)
		schedEmitter := emittertest.NewMockSchedulerEmitter(mockCtrl)
		schedEmitter.EXPECT().EmitScheduler(gomock.Any()).AnyTimes()

		taskEmitter.EXPECT().TaskSuccess(ctx)
		taskEmitter.EXPECT().TaskError(ctx, gomock.Any())
		taskEmitter.EXPECT().TaskDone(ctx, gomock.Any()).Times(2)

		flowEmitter.EXPECT().FlowError(ctx, gomock.Any())
		flowEmitter.EXPECT().FlowDone(ctx, gomock.Any())

		emitter.EXPECT().FlowInit(gomock.Any()).AnyTimes().Return(flowEmitter)
		emitter.EXPECT().TaskInit(gomock.Any(), gomock.Any()).AnyTimes().Return(taskEmitter)
		emitter.EXPECT().SchedulerInit(gomock.Any()).AnyTimes().Return(schedEmitter)

		g := &instrument.CustomEmitter{
			Scope:   scope,
			Logger:  logger,
			Emitter: emitter,
		}
		g.T3795761(ctx, true, true)
	})

	t.Run("should run no error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)

		emitter := emittertest.NewMockEmitter(mockCtrl)

		taskEmitter := emittertest.NewMockTaskEmitter(mockCtrl)
		flowEmitter := emittertest.NewMockFlowEmitter(mockCtrl)
		schedEmitter := emittertest.NewMockSchedulerEmitter(mockCtrl)
		schedEmitter.EXPECT().EmitScheduler(gomock.Any()).AnyTimes()

		taskEmitter.EXPECT().TaskSuccess(ctx).Times(2)
		taskEmitter.EXPECT().TaskDone(ctx, gomock.Any()).Times(2)

		flowEmitter.EXPECT().FlowSuccess(ctx)
		flowEmitter.EXPECT().FlowDone(ctx, gomock.Any())

		emitter.EXPECT().FlowInit(gomock.Any()).AnyTimes().Return(flowEmitter)
		emitter.EXPECT().TaskInit(gomock.Any(), gomock.Any()).AnyTimes().Return(taskEmitter)
		emitter.EXPECT().SchedulerInit(gomock.Any()).AnyTimes().Return(schedEmitter)

		g := &instrument.CustomEmitter{
			Scope:   scope,
			Logger:  logger,
			Emitter: emitter,
		}
		g.T3795761(ctx, true, false)
	})

	t.Run("should not run", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)

		emitter := emittertest.NewMockEmitter(mockCtrl)

		taskEmitter := emittertest.NewMockTaskEmitter(mockCtrl)
		flowEmitter := emittertest.NewMockFlowEmitter(mockCtrl)
		schedEmitter := emittertest.NewMockSchedulerEmitter(mockCtrl)
		schedEmitter.EXPECT().EmitScheduler(gomock.Any()).AnyTimes()

		taskEmitter.EXPECT().TaskSuccess(ctx)
		taskEmitter.EXPECT().TaskSkipped(ctx, gomock.Any())
		taskEmitter.EXPECT().TaskDone(ctx, gomock.Any()).Times(1)

		flowEmitter.EXPECT().FlowSuccess(ctx)
		flowEmitter.EXPECT().FlowDone(ctx, gomock.Any())

		emitter.EXPECT().FlowInit(gomock.Any()).AnyTimes().Return(flowEmitter)
		emitter.EXPECT().TaskInit(gomock.Any(), gomock.Any()).AnyTimes().Return(taskEmitter)
		emitter.EXPECT().SchedulerInit(gomock.Any()).AnyTimes().Return(schedEmitter)

		g := &instrument.CustomEmitter{
			Scope:   scope,
			Logger:  logger,
			Emitter: emitter,
		}
		// false, false is equivalent
		g.T3795761(ctx, false, true)
	})
}

func TestFlowPanic(t *testing.T) {
	t.Run("flow task panic", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ctx := context.Background()

		emitter := emittertest.NewMockEmitter(mockCtrl)

		// No flow emitter as flow isn't instrumented.
		taskEmitter := emittertest.NewMockTaskEmitter(mockCtrl)

		taskEmitter.EXPECT().TaskPanic(ctx, gomock.Any())
		taskEmitter.EXPECT().TaskDone(ctx, gomock.Any())

		schedEmitter := emittertest.NewMockSchedulerEmitter(mockCtrl)
		schedEmitter.EXPECT().EmitScheduler(gomock.Any()).AnyTimes()

		emitter.EXPECT().TaskInit(
			&cff.TaskInfo{
				Name:   "Atoi",
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   552,
				Column: 12,
			},
			&cff.DirectiveInfo{
				Name:      "",
				Directive: cff.FlowDirective,
				File:      "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:      549,
				Column:    9,
			}).Return(taskEmitter)
		emitter.EXPECT().SchedulerInit(gomock.Any()).Return(schedEmitter)

		scope := tally.NewTestScope("", nil)
		core, _ := observer.New(zap.DebugLevel)
		logger := zap.New(core)

		g := &instrument.CustomEmitter{
			Scope:   scope,
			Logger:  logger,
			Emitter: emitter,
		}
		assert.Error(t, g.FlowAlwaysPanics(ctx))
	})

	t.Run("task predicate panic", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ctx := context.Background()

		emitter := emittertest.NewMockEmitter(mockCtrl)
		taskEmitter := emittertest.NewMockTaskEmitter(mockCtrl)
		taskEmitter.EXPECT().TaskPanic(ctx, gomock.Any())
		taskEmitter.EXPECT().TaskSkipped(ctx, gomock.Any())
		schedEmitter := emittertest.NewMockSchedulerEmitter(mockCtrl)
		schedEmitter.EXPECT().EmitScheduler(gomock.Any()).AnyTimes()

		emitter.EXPECT().TaskInit(
			&cff.TaskInfo{
				Name:   "PredicatePanics",
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   565,
				Column: 12,
			},
			&cff.DirectiveInfo{
				Name:      "",
				Directive: cff.FlowDirective,
				File:      "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:      563,
				Column:    9,
			}).Return(taskEmitter)
		emitter.EXPECT().SchedulerInit(gomock.Any()).Return(schedEmitter)

		scope := tally.NewTestScope("", nil)
		core, _ := observer.New(zap.DebugLevel)
		logger := zap.New(core)

		g := &instrument.CustomEmitter{
			Scope:   scope,
			Logger:  logger,
			Emitter: emitter,
		}
		assert.Error(t, g.PredicatePanics(ctx))
	})

	t.Run("task predicate panic with fallback", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ctx := context.Background()

		emitter := emittertest.NewMockEmitter(mockCtrl)
		taskEmitter := emittertest.NewMockTaskEmitter(mockCtrl)
		taskEmitter.EXPECT().TaskPanicRecovered(ctx, gomock.Any())
		taskEmitter.EXPECT().TaskSkipped(ctx, gomock.Any())
		schedEmitter := emittertest.NewMockSchedulerEmitter(mockCtrl)
		schedEmitter.EXPECT().EmitScheduler(gomock.Any()).AnyTimes()

		emitter.EXPECT().TaskInit(
			&cff.TaskInfo{
				Name:   "PredicatePanicsWithFallback",
				File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:   590,
				Column: 4,
			},
			&cff.DirectiveInfo{
				Name:      "",
				Directive: cff.FlowDirective,
				File:      "go.uber.org/cff/internal/tests/instrument/instrument.go",
				Line:      585,
				Column:    8,
			}).Return(taskEmitter)
		emitter.EXPECT().SchedulerInit(gomock.Any()).Return(schedEmitter)

		scope := tally.NewTestScope("", nil)
		core, _ := observer.New(zap.DebugLevel)
		logger := zap.New(core)

		g := &instrument.CustomEmitter{
			Scope:   scope,
			Logger:  logger,
			Emitter: emitter,
		}
		res, err := g.PredicatePanicsWithFallback(ctx)
		require.NoError(t, err)
		assert.Equal(t, res, "predicate-fallback")
	})
}

func TestParallelPanic(t *testing.T) {
	t.Run("cff.Tasks panic", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ctx := context.Background()

		emitter := emittertest.NewMockEmitter(mockCtrl)

		parallelEmitter := emittertest.NewMockParallelEmitter(mockCtrl)
		schedEmitter := emittertest.NewMockSchedulerEmitter(mockCtrl)
		schedEmitter.EXPECT().EmitScheduler(gomock.Any()).AnyTimes()

		parallelEmitter.EXPECT().ParallelError(gomock.Any(), gomock.Any())
		parallelEmitter.EXPECT().ParallelDone(ctx, gomock.Any()).Times(1)

		emitter.EXPECT().ParallelInit(gomock.Any()).AnyTimes().Return(parallelEmitter)
		emitter.EXPECT().SchedulerInit(gomock.Any()).Return(schedEmitter)

		scope := tally.NewTestScope("", nil)
		core, _ := observer.New(zap.DebugLevel)
		logger := zap.New(core)

		g := &instrument.CustomEmitter{
			Scope:   scope,
			Logger:  logger,
			Emitter: emitter,
		}
		assert.Error(t, g.ParallelAlwaysPanics(ctx))
	})

	t.Run("cff.Task panic", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ctx := context.Background()

		emitter := emittertest.NewMockEmitter(mockCtrl)

		taskEmitter := emittertest.NewMockTaskEmitter(mockCtrl)
		parallelEmitter := emittertest.NewMockParallelEmitter(mockCtrl)
		schedEmitter := emittertest.NewMockSchedulerEmitter(mockCtrl)
		schedEmitter.EXPECT().EmitScheduler(gomock.Any()).AnyTimes()

		parallelEmitter.EXPECT().ParallelError(gomock.Any(), gomock.Any())
		parallelEmitter.EXPECT().ParallelDone(ctx, gomock.Any()).Times(1)

		taskEmitter.EXPECT().TaskPanic(ctx, gomock.Any())
		taskEmitter.EXPECT().TaskDone(ctx, gomock.Any())

		emitter.EXPECT().SchedulerInit(gomock.Any()).Return(schedEmitter)
		emitter.EXPECT().ParallelInit(gomock.Any()).AnyTimes().Return(parallelEmitter)
		emitter.EXPECT().TaskInit(gomock.Any(), gomock.Any()).Times(1).Return(taskEmitter)

		scope := tally.NewTestScope("", nil)
		core, _ := observer.New(zap.DebugLevel)
		logger := zap.New(core)

		g := &instrument.CustomEmitter{
			Scope:   scope,
			Logger:  logger,
			Emitter: emitter,
		}
		assert.Error(t, g.ParallelTaskAlwaysPanics(ctx))
	})
}

// TestConcurrentFlow detects data races when multiple flows share the same
// emitter.
func TestConcurrentFlow(t *testing.T) {
	scope := tally.NewTestScope("", nil)
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	builder := cff.TallyEmitter(scope)

	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			benchmark.FlowMetricsMemoized1000(logger, builder)
		}()
	}
	wg.Wait()
}
