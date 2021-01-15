package cff_test

import (
	"context"
	"sync"
	"testing"

	"go.uber.org/cff"
	"go.uber.org/cff/internal/tests/benchmark"
	"go.uber.org/cff/internal/tests/instrument"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

// TestInstrumentEmitter verifies that new Emitter interface gets called if
// it's passed in.
func TestInstrumentEmitter(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	ctx := context.Background()

	emitter := cff.NewMockEmitter(mockCtrl)

	taskEmitter := cff.NewMockTaskEmitter(mockCtrl)
	flowEmitter := cff.NewMockFlowEmitter(mockCtrl)
	schedEmitter := cff.NewMockSchedulerEmitter(mockCtrl)

	flowsucc := flowEmitter.EXPECT().FlowSuccess(ctx)
	flowEmitter.EXPECT().FlowDone(ctx, gomock.Any()).After(flowsucc)

	// 2 tasks.
	taskEmitter.EXPECT().TaskSuccess(ctx).Times(2)
	taskEmitter.EXPECT().TaskDone(ctx, gomock.Any()).Times(2)

	flowInfo := &cff.FlowInfo{
		Name:   "AtoiRun",
		File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
		Line:   227,
		Column: 8,
	}

	emitter.EXPECT().FlowInit(flowInfo).Return(flowEmitter)

	emitter.EXPECT().SchedulerInit(
		&cff.SchedulerInfo{
			FlowInfo: flowInfo,
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
	v, err := g.Run(ctx, "1")

	assert.NoError(t, err)
	assert.Equal(t, uint8(1), v)
}

func TestInstrumentErrorME(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	ctx := context.Background()

	emitter := cff.NewMockEmitter(mockCtrl)

	taskEmitter := cff.NewMockTaskEmitter(mockCtrl)
	flowEmitter := cff.NewMockFlowEmitter(mockCtrl)
	schedEmitter := cff.NewMockSchedulerEmitter(mockCtrl)

	// flowFailedEmitter := cff.NewMockFlowEmitter(mockCtrl)

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
	_, err := h.Run(ctx, "NaN")

	assert.Error(t, err)
}

func TestInstrumentTaskButNotFlowME(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	ctx := context.Background()

	emitter := cff.NewMockEmitter(mockCtrl)

	taskEmitter := cff.NewMockTaskEmitter(mockCtrl)
	schedEmitter := cff.NewMockSchedulerEmitter(mockCtrl)

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
	v, err := g.Work(ctx, "1")

	assert.NoError(t, err)
	assert.Equal(t, 1, v)
}

func TestInstrumentCancelledContextME(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	scope := tally.NewTestScope("", nil)
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()

	flowCancelledErr := ctx.Err()

	emitter := cff.NewMockEmitter(mockCtrl)

	taskEmitter := cff.NewMockTaskEmitter(mockCtrl)
	flowEmitter := cff.NewMockFlowEmitter(mockCtrl)
	schedEmitter := cff.NewMockSchedulerEmitter(mockCtrl)

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

	_, err := g.Run(ctx, "1")
	assert.Error(t, err)
}

func TestInstrumentRecoverME(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	ctx := context.Background()

	scope := tally.NewTestScope("", nil)
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	emitter := cff.NewMockEmitter(mockCtrl)

	taskEmitter := cff.NewMockTaskEmitter(mockCtrl)
	flowEmitter := cff.NewMockFlowEmitter(mockCtrl)
	schedEmitter := cff.NewMockSchedulerEmitter(mockCtrl)

	flowEmitter.EXPECT().FlowSuccess(ctx)
	flowEmitter.EXPECT().FlowDone(ctx, gomock.Any())

	taskEmitter.EXPECT().TaskSuccess(ctx)
	taskEmitter.EXPECT().TaskErrorRecovered(ctx, gomock.Any())
	taskEmitter.EXPECT().TaskDone(ctx, gomock.Any()).Times(2)

	emitter.EXPECT().FlowInit(&cff.FlowInfo{
		Name:   "AtoiRun",
		File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
		Line:   227,
		Column: 8,
	}).Return(flowEmitter)
	emitter.EXPECT().TaskInit(gomock.Any(), gomock.Any()).Times(2).Return(taskEmitter)
	emitter.EXPECT().SchedulerInit(gomock.Any()).AnyTimes().Return(schedEmitter)

	g := &instrument.CustomEmitter{
		Scope:   scope,
		Logger:  logger,
		Emitter: emitter,
	}

	v, err := g.Run(ctx, "300")

	assert.NoError(t, err)
	assert.Equal(t, uint8(0), v)
}

// TestT3630161ME tests against regression for T3630161
func TestT3630161ME(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	ctx := context.Background()

	emitter := cff.NewMockEmitter(mockCtrl)

	taskEmitter := cff.NewMockTaskEmitter(mockCtrl)
	flowEmitter := cff.NewMockFlowEmitter(mockCtrl)
	schedEmitter := cff.NewMockSchedulerEmitter(mockCtrl)

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

		emitter := cff.NewMockEmitter(mockCtrl)

		taskEmitter := cff.NewMockTaskEmitter(mockCtrl)
		flowEmitter := cff.NewMockFlowEmitter(mockCtrl)
		schedEmitter := cff.NewMockSchedulerEmitter(mockCtrl)

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

		emitter := cff.NewMockEmitter(mockCtrl)

		taskEmitter := cff.NewMockTaskEmitter(mockCtrl)
		flowEmitter := cff.NewMockFlowEmitter(mockCtrl)
		schedEmitter := cff.NewMockSchedulerEmitter(mockCtrl)

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

		emitter := cff.NewMockEmitter(mockCtrl)

		taskEmitter := cff.NewMockTaskEmitter(mockCtrl)
		flowEmitter := cff.NewMockFlowEmitter(mockCtrl)
		schedEmitter := cff.NewMockSchedulerEmitter(mockCtrl)

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

func TestPanic(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	ctx := context.Background()

	emitter := cff.NewMockEmitter(mockCtrl)

	// No flow emitter as flow isn't instrumented.
	taskEmitter := cff.NewMockTaskEmitter(mockCtrl)

	taskEmitter.EXPECT().TaskPanic(ctx, gomock.Any())
	taskEmitter.EXPECT().TaskDone(ctx, gomock.Any())

	schedEmitter := cff.NewMockSchedulerEmitter(mockCtrl)

	emitter.EXPECT().TaskInit(
		&cff.TaskInfo{
			Name:   "Atoi",
			File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
			Line:   349,
			Column: 12,
		},
		&cff.FlowInfo{
			Name:   "",
			File:   "go.uber.org/cff/internal/tests/instrument/instrument.go",
			Line:   346,
			Column: 9,
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
	err := g.FlowAlwaysPanics(ctx)
	require.Error(t, err)
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

			benchmark.MetricsMemoized1000(logger, builder)
		}()
	}
	wg.Wait()

}
