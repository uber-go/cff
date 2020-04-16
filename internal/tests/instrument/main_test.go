package instrument

import (
	"context"
	"testing"

	"go.uber.org/cff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
)

func TestInstrument(t *testing.T) {
	scope := tally.NewTestScope("", nil)
	core, observedLogs := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	h := &DefaultEmitter{Scope: scope, Logger: logger}
	ctx := context.Background()
	v, err := h.Run(ctx, "1")

	assert.NoError(t, err)
	assert.Equal(t, uint8(1), v)

	metrics := scope.Snapshot()
	// metrics
	counters := metrics.Counters()
	for k := range counters {
		t.Logf("got counter with key %q", k)
	}
	assert.Equal(t, int64(1), counters["task.success+flow=AtoiRun,task=Atoi"].Value())
	assert.Equal(t, int64(1), counters["task.success+flow=AtoiRun,task=uint8"].Value())
	assert.Equal(t, int64(1), counters["taskflow.success+flow=AtoiRun"].Value())

	timers := metrics.Timers()
	assert.NotNil(t, timers["task.timing+flow=AtoiRun,task=Atoi"])
	assert.NotNil(t, timers["taskflow.timing+flow=AtoiRun"])

	// logs
	expectedLevel := zap.DebugLevel
	expectedMessages := []string{
		"task success",
		"task done",
		"task success",
		"task done",
		"flow success",
		"flow done",
	}
	logEntries := observedLogs.All()
	assert.Equal(t, len(expectedMessages), len(logEntries))
	for _, entry := range logEntries {
		t.Logf("log entry - level: %q, message: %q, fields: %v", entry.Level, entry.Message, entry.ContextMap())
	}
	for i, entry := range logEntries {
		assert.Equal(t, expectedLevel, entry.Level)
		assert.Equal(t, expectedMessages[i], entry.Message)
	}
}

func TestInstrumentError(t *testing.T) {
	scope := tally.NewTestScope("", nil)
	core, observedLogs := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	h := &DefaultEmitter{Scope: scope, Logger: logger}
	ctx := context.Background()
	_, err := h.Run(ctx, "NaN")

	assert.Error(t, err)

	// metrics
	counters := scope.Snapshot().Counters()
	for k, v := range counters {
		t.Logf("got counter with key %q val %v", k, v.Value())
	}
	assert.Equal(t, int64(1), counters["task.error+flow=AtoiRun,task=Atoi"].Value())
	assert.Equal(t, int64(1), counters["taskflow.error+flow=AtoiRun"].Value())

	expected := []struct {
		level   zapcore.Level
		message string
		fields  map[string]interface{}
	}{
		{zap.DebugLevel, "task error", map[string]interface{}{"task": "Atoi"}},
		{zap.DebugLevel, "task done", map[string]interface{}{"task": "Atoi"}},
		{zap.DebugLevel, "flow error", nil},
		{zap.DebugLevel, "task skipped", map[string]interface{}{"task": "uint8"}},
		{zap.DebugLevel, "flow skipped", nil},
		{zap.DebugLevel, "flow done", nil},
	}

	// logs
	logEntries := observedLogs.All()
	for _, entry := range logEntries {
		t.Logf("log entry - level: %q, message: %q, fields: %v", entry.Level, entry.Message, entry.ContextMap())
	}
	assert.Equal(t, len(expected), len(logEntries))
	for i, entry := range logEntries {
		assert.Equal(t, expected[i].level, entry.Level)
		assert.Equal(t, expected[i].message, entry.Message)
		for k, v := range expected[i].fields {
			actualValue, ok := entry.ContextMap()[k]
			assert.True(t, ok)
			assert.Equal(t, v, actualValue)
		}
	}
}

func TestInstrumentCancelledContext(t *testing.T) {
	scope := tally.NewTestScope("", nil)
	core, observedLogs := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()

	h := &DefaultEmitter{Scope: scope, Logger: logger}
	_, err := h.Run(ctx, "1")
	assert.Error(t, err)

	// metrics
	counters := scope.Snapshot().Counters()
	for k := range counters {
		t.Logf("got counter with key %q", k)
	}
	assert.Equal(t, int64(1), counters["task.skipped+flow=AtoiRun,task=Atoi"].Value())
	assert.Equal(t, int64(1), counters["task.skipped+flow=AtoiRun,task=uint8"].Value())
	assert.Equal(t, int64(1), counters["taskflow.skipped+flow=AtoiRun"].Value())

	// logs
	expectedLevel := zap.DebugLevel
	expectedMessages := []string{
		"flow error",
		"task skipped",
		"task skipped",
		"flow skipped",
		"flow done",
	}
	logEntries := observedLogs.All()
	assert.Equal(t, len(expectedMessages), len(logEntries))
	for i, entry := range logEntries {
		t.Logf("log entry - level: %q, message: %q, fields: %v", entry.Level, entry.Message, entry.Context)
		assert.Equal(t, expectedLevel, entry.Level)
		assert.Equal(t, expectedMessages[i], entry.Message)
	}
}

func TestInstrumentRecover(t *testing.T) {
	scope := tally.NewTestScope("", nil)
	core, observedLogs := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	h := &DefaultEmitter{Scope: scope, Logger: logger}
	ctx := context.Background()
	v, err := h.Run(ctx, "300")

	assert.NoError(t, err)
	assert.Equal(t, uint8(0), v)

	// metrics
	counters := scope.Snapshot().Counters()
	for k := range counters {
		t.Logf("got counter with key %q", k)
	}
	assert.Equal(t, int64(1), counters["task.success+flow=AtoiRun,task=Atoi"].Value())
	assert.Equal(t, int64(1), counters["task.recovered+flow=AtoiRun,task=uint8"].Value())
	assert.Equal(t, int64(1), counters["taskflow.success+flow=AtoiRun"].Value())

	// logs
	expected := []struct {
		level   zapcore.Level
		message string
	}{
		{zap.DebugLevel, "task success"},
		{zap.DebugLevel, "task done"},
		{zap.ErrorLevel, "task error recovered"},
		{zap.DebugLevel, "task done"},
		{zap.DebugLevel, "flow success"},
		{zap.DebugLevel, "flow done"},
	}
	logEntries := observedLogs.All()
	assert.Equal(t, len(expected), len(logEntries))
	for i, entry := range logEntries {
		t.Logf("log entry - level: %q, message: %q, fields: %v", entry.Level, entry.Message, entry.Context)
		assert.Equal(t, expected[i].level, entry.Level)
		assert.Equal(t, expected[i].message, entry.Message)
	}
}

func TestInstrumentAnnotationOrder(t *testing.T) {
	scope := tally.NewTestScope("", nil)
	core, observedLogs := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	h := &DefaultEmitter{Scope: scope, Logger: logger}
	ctx := context.Background()
	v, err := h.Do(ctx, "1")

	assert.NoError(t, err)
	assert.Equal(t, 1, v)

	// metrics
	counters := scope.Snapshot().Counters()
	for k := range counters {
		t.Logf("got counter with key %q", k)
	}
	assert.Equal(t, int64(1), counters["task.success+flow=AtoiDo,task=Atoi"].Value())
	assert.Equal(t, int64(1), counters["taskflow.success+flow=AtoiDo"].Value())

	// logs
	expectedLevel := zap.DebugLevel
	expectedMessages := []string{
		"task success",
		"task done",
		"flow success",
		"flow done",
	}
	logEntries := observedLogs.All()
	assert.Equal(t, len(expectedMessages), len(logEntries))
	for i, entry := range logEntries {
		t.Logf("log entry - level: %q, message: %q, fields: %v", entry.Level, entry.Message, entry.Context)
		assert.Equal(t, expectedLevel, entry.Level)
		assert.Equal(t, expectedMessages[i], entry.Message)
	}
}

func TestInstrumentTaskButNotFlow(t *testing.T) {
	scope := tally.NewTestScope("", nil)
	core, observedLogs := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	h := &DefaultEmitter{Scope: scope, Logger: logger}
	ctx := context.Background()
	v, err := h.Work(ctx, "1")

	assert.NoError(t, err)
	assert.Equal(t, 1, v)

	// metrics
	counters := scope.Snapshot().Counters()
	for k := range counters {
		t.Logf("got counter with key %q", k)
	}
	assert.Equal(t, int64(1), counters["task.success+task=Atoi"].Value())

	// logs
	expectedLevel := zap.DebugLevel
	expectedMessages := []string{
		"task success",
		"task done",
	}
	logEntries := observedLogs.All()
	assert.Equal(t, len(expectedMessages), len(logEntries))
	for i, entry := range logEntries {
		t.Logf("log entry - level: %q, message: %q, fields: %v", entry.Level, entry.Message, entry.Context)
		assert.Equal(t, expectedLevel, entry.Level)
		assert.Equal(t, expectedMessages[i], entry.Message)
	}
}

// TestT3630161 tests against regression for T3630161
func TestT3630161(t *testing.T) {
	scope := tally.NewTestScope("", nil)
	logger := zaptest.NewLogger(t)
	h := &DefaultEmitter{Scope: scope, Logger: logger}
	ctx := context.Background()
	h.T3630161(ctx)

	// metrics
	counters := scope.Snapshot().Counters()
	countersByName := make(map[string][]tally.CounterSnapshot)
	for k := range counters {
		name := counters[k].Name()
		countersByName[name] = append(countersByName[name], counters[k])
	}

	assert.Equal(t, 1, len(countersByName["task.success"]))
	assert.Equal(t, map[string]string{"flow": "T3630161", "task": "End"}, countersByName["task.success"][0].Tags())
	assert.Equal(t, 1, len(countersByName["task.recovered"]))
	assert.Equal(t, map[string]string{"flow": "T3630161", "task": "Err"}, countersByName["task.recovered"][0].Tags())
	assert.Equal(t, 1, len(countersByName["task.recovered"]))
	assert.Equal(t, map[string]string{"flow": "T3630161"}, countersByName["taskflow.success"][0].Tags())
}

// TestT3795761 tests against regression for T3795761 where a task that
// returns no error is not reported as skipped when an earlier task that it
// depends on returns an error.
func TestT3795761(t *testing.T) {
	scope := tally.NewTestScope("", nil)
	core, observedLogs := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	h := &DefaultEmitter{
		Scope:  scope,
		Logger: logger,
	}
	ctx := context.Background()

	expectedLevel := zap.DebugLevel

	t.Run("should run error", func(t *testing.T) {
		h.T3795761(ctx, true, true)

		// logs
		expectedMessages := []string{
			"task success",
			"task done",
			"task error",
			"task done",
			"flow error",
			"flow skipped",
			"flow done",
		}
		logEntries := observedLogs.TakeAll()
		for _, entry := range logEntries {
			t.Logf("log entry - level: %q, message: %q, fields: %v", entry.Level, entry.Message, entry.ContextMap())
		}
		require.Equal(t, len(expectedMessages), len(logEntries))
		for i, entry := range logEntries {
			assert.Equal(t, expectedLevel, entry.Level)
			assert.Equal(t, expectedMessages[i], entry.Message)
		}
	})

	t.Run("should run no error", func(t *testing.T) {
		h.T3795761(ctx, true, false)

		expectedMessages := []string{
			"task success",
			"task done",
			"task success",
			"task done",
			"flow success",
			"flow done",
		}
		logEntries := observedLogs.TakeAll()
		for _, entry := range logEntries {
			t.Logf("log entry - level: %q, message: %q, fields: %v", entry.Level, entry.Message, entry.ContextMap())
		}
		require.Equal(t, len(expectedMessages), len(logEntries))
		for i, entry := range logEntries {
			assert.Equal(t, expectedLevel, entry.Level)
			assert.Equal(t, expectedMessages[i], entry.Message)
		}
	})

	t.Run("should not run", func(t *testing.T) {
		// false, false is equivalent
		h.T3795761(ctx, false, true)

		expectedMessages := []string{
			"task success",
			"task done",
			"task done",
			"flow success",
			"task skipped",
			"flow done",
		}
		logEntries := observedLogs.TakeAll()
		for _, entry := range logEntries {
			t.Logf("log entry - level: %q, message: %q, fields: %v", entry.Level, entry.Message, entry.ContextMap())
		}
		require.Equal(t, len(expectedMessages), len(logEntries))
		for i, entry := range logEntries {
			assert.Equal(t, expectedLevel, entry.Level)
			assert.Equal(t, expectedMessages[i], entry.Message)
		}
	})
}

func TestWithMultipleEmitters(t *testing.T) {
	core1, logs1 := observer.New(zapcore.DebugLevel)
	core2, logs2 := observer.New(zapcore.DebugLevel)

	n, err := AtoiWithTwoEmitters(context.Background(),
		cff.LogEmitter(zap.New(core1)),
		cff.LogEmitter(zap.New(core2)),
		"42",
	)
	require.NoError(t, err)
	assert.Equal(t, 42, n)

	assert.Equal(t, logs1.AllUntimed(), logs2.AllUntimed(), "logs did not match")
}
