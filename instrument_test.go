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

	metricsEmitter := cff.NewMockMetricsEmitter(mockCtrl)

	taskEmitter := cff.NewMockTaskEmitter(mockCtrl)
	flowEmitter := cff.NewMockFlowEmitter(mockCtrl)

	flowsucc := flowEmitter.EXPECT().FlowSuccess()
	flowEmitter.EXPECT().FlowDone(gomock.Any()).After(flowsucc)

	// 2 tasks.
	taskEmitter.EXPECT().TaskSuccess().Times(2)
	taskEmitter.EXPECT().TaskDone(gomock.Any()).Times(2)

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
	v, err := g.Run(context.Background(), "1")

	assert.NoError(t, err)
	assert.Equal(t, uint8(1), v)
}

func TestInstrumentErrorME(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	metricsEmitter := cff.NewMockMetricsEmitter(mockCtrl)

	taskEmitter := cff.NewMockTaskEmitter(mockCtrl)
	flowEmitter := cff.NewMockFlowEmitter(mockCtrl)

	flowFailedEmitter := cff.NewMockFlowEmitter(mockCtrl)

	flowFailedEmitter.EXPECT().FlowError()
	flowFailedEmitter.EXPECT().FlowSkipped()
	flowFailedEmitter.EXPECT().FlowDone(gomock.Any())
	flowEmitter.EXPECT().FlowFailedTask("Atoi").Return(flowFailedEmitter)
	// 2 tasks.
	taskEmitter.EXPECT().TaskError()
	taskEmitter.EXPECT().TaskSkipped()
	taskEmitter.EXPECT().TaskDone(gomock.Any())

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
	ctx := context.Background()
	_, err := h.Run(ctx, "NaN")

	assert.Error(t, err)
}

func TestInstrumentTaskButNotFlowME(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	metricsEmitter := cff.NewMockMetricsEmitter(mockCtrl)

	taskEmitter := cff.NewMockTaskEmitter(mockCtrl)

	taskEmitter.EXPECT().TaskSuccess()
	taskEmitter.EXPECT().TaskDone(gomock.Any())
	metricsEmitter.EXPECT().TaskInit(gomock.Any()).Return(taskEmitter)

	scope := tally.NewTestScope("", nil)
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	g := &instrument.CustomMetricsEmitter{
		Scope:          scope,
		Logger:         logger,
		MetricsEmitter: metricsEmitter,
	}
	v, err := g.Work(context.Background(), "1")

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

	metricsEmitter := cff.NewMockMetricsEmitter(mockCtrl)

	taskEmitter := cff.NewMockTaskEmitter(mockCtrl)
	flowEmitter := cff.NewMockFlowEmitter(mockCtrl)

	flowEmitter.EXPECT().FlowSkipped()
	flowEmitter.EXPECT().FlowDone(gomock.Any())

	taskEmitter.EXPECT().TaskSkipped().Times(2)

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

	scope := tally.NewTestScope("", nil)
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	metricsEmitter := cff.NewMockMetricsEmitter(mockCtrl)

	taskEmitter := cff.NewMockTaskEmitter(mockCtrl)
	flowEmitter := cff.NewMockFlowEmitter(mockCtrl)

	flowEmitter.EXPECT().FlowSuccess()
	flowEmitter.EXPECT().FlowDone(gomock.Any())

	taskEmitter.EXPECT().TaskError()
	taskEmitter.EXPECT().TaskSuccess()
	taskEmitter.EXPECT().TaskRecovered()
	taskEmitter.EXPECT().TaskDone(gomock.Any()).Times(2)

	metricsEmitter.EXPECT().FlowInit("AtoiRun").Return(flowEmitter)
	metricsEmitter.EXPECT().TaskInit(gomock.Any()).Times(2).Return(taskEmitter)

	g := &instrument.CustomMetricsEmitter{
		Scope:          scope,
		Logger:         logger,
		MetricsEmitter: metricsEmitter,
	}

	v, err := g.Run(context.Background(), "300")

	assert.NoError(t, err)
	assert.Equal(t, uint8(0), v)
}

// TestT3630161ME tests against regression for T3630161
func TestT3630161ME(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	metricsEmitter := cff.NewMockMetricsEmitter(mockCtrl)

	taskEmitter := cff.NewMockTaskEmitter(mockCtrl)
	flowEmitter := cff.NewMockFlowEmitter(mockCtrl)

	// flowsucc := flowEmitter.EXPECT().FlowSuccess()
	flowEmitter.EXPECT().FlowSuccess()
	flowEmitter.EXPECT().FlowDone(gomock.Any())

	// 2 tasks.
	taskEmitter.EXPECT().TaskError()
	taskEmitter.EXPECT().TaskRecovered()
	taskEmitter.EXPECT().TaskDone(gomock.Any()).Times(2)
	taskEmitter.EXPECT().TaskSuccess()

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

	g.T3630161(context.Background())
}

// TestT3795761 tests against regression for T3795761 where a task that returns no error is not reported as
// skipped when an earlier task that it depends on returns an error.
func TestT3795761ME(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	metricsEmitter := cff.NewMockMetricsEmitter(mockCtrl)

	taskEmitter := cff.NewMockTaskEmitter(mockCtrl)
	flowEmitter := cff.NewMockFlowEmitter(mockCtrl)
	flowFailedEmitter := cff.NewMockFlowEmitter(mockCtrl)

	flowEmitter.EXPECT().FlowFailedTask("NeedsInt").Return(flowFailedEmitter)
	flowEmitter.EXPECT().FlowDone(gomock.Any()).AnyTimes()
	flowEmitter.EXPECT().FlowSuccess().AnyTimes()

	flowFailedEmitter.EXPECT().FlowError()
	flowFailedEmitter.EXPECT().FlowSkipped()
	flowFailedEmitter.EXPECT().FlowDone(gomock.Any())

	taskEmitter.EXPECT().TaskSuccess().AnyTimes()
	taskEmitter.EXPECT().TaskError()
	taskEmitter.EXPECT().TaskSkipped()
	taskEmitter.EXPECT().TaskDone(gomock.Any()).AnyTimes()

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
	ctx := context.Background()

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

	metricsEmitter := cff.NewMockMetricsEmitter(mockCtrl)

	// No flow emitter as flow isn't instrumented.
	taskEmitter := cff.NewMockTaskEmitter(mockCtrl)

	taskEmitter.EXPECT().TaskSkipped()
	tpanic := taskEmitter.EXPECT().TaskPanic()

	taskEmitter.EXPECT().TaskDone(gomock.Any()).After(tpanic)

	metricsEmitter.EXPECT().TaskInit("Atoi").Return(taskEmitter)

	scope := tally.NewTestScope("", nil)
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	g := &instrument.CustomMetricsEmitter{
		Scope:          scope,
		Logger:         logger,
		MetricsEmitter: metricsEmitter,
	}
	err := g.FlowAlwaysPanics()
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
