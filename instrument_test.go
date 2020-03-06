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

// TestMetricsEmitter verifies that new MetricsEmitter interface gets called
// if it's passed in.
func TestInstrumentME(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx := context.Background()

	metricsEmitter := cff.NewMockMetricsEmitter(mockCtrl)

	taskEmitter := cff.NewMockTaskEmitter(mockCtrl)
	flowEmitter := cff.NewMockFlowEmitter(mockCtrl)

	flowsucc := flowEmitter.EXPECT().FlowSuccess(ctx)
	flowEmitter.EXPECT().FlowDone(ctx, gomock.Any()).After(flowsucc)

	// 2 tasks.
	taskEmitter.EXPECT().TaskSuccess(ctx).Times(2)
	taskEmitter.EXPECT().TaskDone(ctx, gomock.Any()).Times(2)

	metricsEmitter.EXPECT().FlowInit("AtoiRun").Return(flowEmitter)
	// 2 in the tasks for loop inside defer() and twice after.
	metricsEmitter.EXPECT().TaskInit(gomock.Any()).Times(2).Return(taskEmitter)

	scope := tally.NewTestScope("", nil)
	// Logging
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	// CustomMetricsEmitter Unit
	g := &instrument.CustomMetricsEmitter{
		Logger:         logger,
		Scope:          scope,
		MetricsEmitter: metricsEmitter,
	}
	v, err := g.Run(ctx, "1")

	assert.NoError(t, err)
	assert.Equal(t, uint8(1), v)
}

func TestInstrumentErrorME(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx := context.Background()

	metricsEmitter := cff.NewMockMetricsEmitter(mockCtrl)

	taskEmitter := cff.NewMockTaskEmitter(mockCtrl)
	flowEmitter := cff.NewMockFlowEmitter(mockCtrl)

	flowFailedEmitter := cff.NewMockFlowEmitter(mockCtrl)

	flowFailedEmitter.EXPECT().FlowError(ctx, gomock.Any())
	flowFailedEmitter.EXPECT().FlowSkipped(ctx, gomock.Any())
	flowFailedEmitter.EXPECT().FlowDone(ctx, gomock.Any())
	flowEmitter.EXPECT().FlowFailedTask(ctx, "Atoi", gomock.Any()).Return(flowFailedEmitter)
	// 2 tasks.
	taskEmitter.EXPECT().TaskError(ctx, gomock.Any())
	taskEmitter.EXPECT().TaskSkipped(ctx, gomock.Any())
	taskEmitter.EXPECT().TaskDone(ctx, gomock.Any())

	metricsEmitter.EXPECT().FlowInit("AtoiRun").Return(flowEmitter)
	metricsEmitter.EXPECT().TaskInit(gomock.Any()).Times(2).Return(taskEmitter)

	scope := tally.NewTestScope("", nil)
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	h := &instrument.CustomMetricsEmitter{
		Scope:          scope,
		Logger:         logger,
		MetricsEmitter: metricsEmitter,
	}
	_, err := h.Run(ctx, "NaN")

	assert.Error(t, err)
}

func TestInstrumentTaskButNotFlowME(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx := context.Background()

	metricsEmitter := cff.NewMockMetricsEmitter(mockCtrl)

	taskEmitter := cff.NewMockTaskEmitter(mockCtrl)

	taskEmitter.EXPECT().TaskSuccess(ctx)
	taskEmitter.EXPECT().TaskDone(ctx, gomock.Any())
	metricsEmitter.EXPECT().TaskInit(gomock.Any()).Return(taskEmitter)

	scope := tally.NewTestScope("", nil)
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	g := &instrument.CustomMetricsEmitter{
		Scope:          scope,
		Logger:         logger,
		MetricsEmitter: metricsEmitter,
	}
	v, err := g.Work(ctx, "1")

	assert.NoError(t, err)
	assert.Equal(t, 1, v)
}

func TestInstrumentCancelledContextME(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	scope := tally.NewTestScope("", nil)
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()

	flowCancelledErr := ctx.Err()

	metricsEmitter := cff.NewMockMetricsEmitter(mockCtrl)

	taskEmitter := cff.NewMockTaskEmitter(mockCtrl)
	flowEmitter := cff.NewMockFlowEmitter(mockCtrl)

	flowEmitter.EXPECT().FlowSkipped(ctx, flowCancelledErr)
	flowEmitter.EXPECT().FlowDone(ctx, gomock.Any())

	taskEmitter.EXPECT().TaskSkipped(ctx, gomock.Any()).Times(2)

	metricsEmitter.EXPECT().FlowInit("AtoiRun").Return(flowEmitter)
	metricsEmitter.EXPECT().TaskInit(gomock.Any()).AnyTimes().Return(taskEmitter)

	g := &instrument.CustomMetricsEmitter{
		Scope:          scope,
		Logger:         logger,
		MetricsEmitter: metricsEmitter,
	}

	_, err := g.Run(ctx, "1")
	assert.Error(t, err)
}

func TestInstrumentRecoverME(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx := context.Background()

	scope := tally.NewTestScope("", nil)
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	metricsEmitter := cff.NewMockMetricsEmitter(mockCtrl)

	taskEmitter := cff.NewMockTaskEmitter(mockCtrl)
	flowEmitter := cff.NewMockFlowEmitter(mockCtrl)

	flowEmitter.EXPECT().FlowSuccess(ctx)
	flowEmitter.EXPECT().FlowDone(ctx, gomock.Any())

	taskEmitter.EXPECT().TaskError(ctx, gomock.Any())
	taskEmitter.EXPECT().TaskSuccess(ctx)
	taskEmitter.EXPECT().TaskRecovered(ctx, gomock.Any())
	taskEmitter.EXPECT().TaskDone(ctx, gomock.Any()).Times(2)

	metricsEmitter.EXPECT().FlowInit("AtoiRun").Return(flowEmitter)
	metricsEmitter.EXPECT().TaskInit(gomock.Any()).Times(2).Return(taskEmitter)

	g := &instrument.CustomMetricsEmitter{
		Scope:          scope,
		Logger:         logger,
		MetricsEmitter: metricsEmitter,
	}

	v, err := g.Run(ctx, "300")

	assert.NoError(t, err)
	assert.Equal(t, uint8(0), v)
}

// TestT3630161ME tests against regression for T3630161
func TestT3630161ME(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx := context.Background()

	metricsEmitter := cff.NewMockMetricsEmitter(mockCtrl)

	taskEmitter := cff.NewMockTaskEmitter(mockCtrl)
	flowEmitter := cff.NewMockFlowEmitter(mockCtrl)

	// flowsucc := flowEmitter.EXPECT().FlowSuccess()
	flowEmitter.EXPECT().FlowSuccess(ctx)
	flowEmitter.EXPECT().FlowDone(ctx, gomock.Any())

	// 2 tasks.
	taskEmitter.EXPECT().TaskError(ctx, gomock.Any())
	taskEmitter.EXPECT().TaskRecovered(ctx, gomock.Any())
	taskEmitter.EXPECT().TaskDone(ctx, gomock.Any()).Times(2)
	taskEmitter.EXPECT().TaskSuccess(ctx)

	metricsEmitter.EXPECT().FlowInit("T3630161").Return(flowEmitter)
	metricsEmitter.EXPECT().TaskInit(gomock.Any()).Times(2).Return(taskEmitter)

	scope := tally.NewTestScope("", nil)
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	g := &instrument.CustomMetricsEmitter{
		Scope:          scope,
		Logger:         logger,
		MetricsEmitter: metricsEmitter,
	}

	g.T3630161(ctx)
}

// TestT3795761 tests against regression for T3795761 where a task that returns no error is not reported as
// skipped when an earlier task that it depends on returns an error.
func TestT3795761ME(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx := context.Background()

	metricsEmitter := cff.NewMockMetricsEmitter(mockCtrl)

	taskEmitter := cff.NewMockTaskEmitter(mockCtrl)
	flowEmitter := cff.NewMockFlowEmitter(mockCtrl)
	flowFailedEmitter := cff.NewMockFlowEmitter(mockCtrl)

	flowEmitter.EXPECT().FlowFailedTask(ctx, "NeedsInt", gomock.Any()).Return(flowFailedEmitter)
	flowEmitter.EXPECT().FlowDone(ctx, gomock.Any()).AnyTimes()
	flowEmitter.EXPECT().FlowSuccess(ctx).AnyTimes()

	flowFailedEmitter.EXPECT().FlowError(ctx, gomock.Any())
	flowFailedEmitter.EXPECT().FlowSkipped(ctx, gomock.Any())
	flowFailedEmitter.EXPECT().FlowDone(ctx, gomock.Any())

	taskEmitter.EXPECT().TaskSuccess(ctx).AnyTimes()
	taskEmitter.EXPECT().TaskError(ctx, gomock.Any())
	taskEmitter.EXPECT().TaskSkipped(ctx, gomock.Any())
	taskEmitter.EXPECT().TaskDone(ctx, gomock.Any()).AnyTimes()

	metricsEmitter.EXPECT().FlowInit(gomock.Any()).AnyTimes().Return(flowEmitter)
	metricsEmitter.EXPECT().TaskInit(gomock.Any()).AnyTimes().Return(taskEmitter)

	scope := tally.NewTestScope("", nil)
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	g := &instrument.CustomMetricsEmitter{
		Scope:          scope,
		Logger:         logger,
		MetricsEmitter: metricsEmitter,
	}

	t.Run("should run error", func(t *testing.T) {
		g.T3795761(ctx, true, true)
	})

	t.Run("should run no error", func(t *testing.T) {
		g.T3795761(ctx, true, false)
	})

	t.Run("should not run", func(t *testing.T) {
		// false, false is equivalent
		g.T3795761(ctx, false, true)
	})
}

func TestPanic(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx := context.Background()

	metricsEmitter := cff.NewMockMetricsEmitter(mockCtrl)

	// No flow emitter as flow isn't instrumented.
	taskEmitter := cff.NewMockTaskEmitter(mockCtrl)

	taskEmitter.EXPECT().TaskSkipped(ctx, gomock.Any())
	tpanic := taskEmitter.EXPECT().TaskPanic(ctx, gomock.Any())

	taskEmitter.EXPECT().TaskDone(ctx, gomock.Any()).After(tpanic)

	metricsEmitter.EXPECT().TaskInit("Atoi").Return(taskEmitter)

	scope := tally.NewTestScope("", nil)
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	g := &instrument.CustomMetricsEmitter{
		Scope:          scope,
		Logger:         logger,
		MetricsEmitter: metricsEmitter,
	}
	err := g.FlowAlwaysPanics(ctx)
	require.Error(t, err)
}

// TestConcurrentFlow detects data races when multiple flows share the same
// metrics emitter.
func TestConcurrentFlow(t *testing.T) {
	scope := tally.NewTestScope("", nil)
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	builder := cff.DefaultMetricsEmitter(scope)

	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			benchmark.MetricsMemoized1000(logger, scope, builder)
		}()
	}
	wg.Wait()

}
