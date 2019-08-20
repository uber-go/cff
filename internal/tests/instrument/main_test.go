package instrument

import (
	"context"
	"testing"

	"go.uber.org/zap/zaptest"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestInstrument(t *testing.T) {
	scope := tally.NewTestScope("", nil)
	core, observedLogs := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	h := &H{Scope: scope, Logger: logger}
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
	assert.Equal(t, int64(1), counters["task.success+task=Atoi"].Value())
	assert.Equal(t, int64(1), counters["task.success+task=uint8"].Value())
	assert.Equal(t, int64(1), counters["taskflow.success+flow=AtoiRun"].Value())

	timers := metrics.Timers()
	assert.NotNil(t, timers["task.timing+task=Atoi"])
	assert.NotNil(t, timers["taskflow.timing+flow=AtoiRun"])

	// logs
	expectedLevel := zap.DebugLevel
	expectedMessages := []string{
		"task succeeded",
		"task succeeded",
		"taskflow succeeded",
	}
	logEntries := observedLogs.All()
	assert.Equal(t, len(expectedMessages), len(logEntries))
	for i, entry := range logEntries {
		assert.Equal(t, expectedLevel, entry.Level)
		assert.Equal(t, expectedMessages[i], entry.Message)
		t.Logf("log entry - level: %q, message: %q, fields: %v", entry.Level, entry.Message, entry.Context)
	}
}

func TestInstrumentError(t *testing.T) {
	scope := tally.NewTestScope("", nil)
	core, observedLogs := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	h := &H{Scope: scope, Logger: logger}
	ctx := context.Background()
	_, err := h.Run(ctx, "NaN")

	assert.Error(t, err)

	// metrics
	counters := scope.Snapshot().Counters()
	for k, v := range counters {
		t.Logf("got counter with key %q val %v", k, v.Value())
	}
	assert.Equal(t, int64(1), counters["task.error+task=Atoi"].Value())
	assert.Equal(t, int64(1), counters["taskflow.error+failedtask=Atoi,flow=AtoiRun"].Value())

	expected := []struct {
		level   zapcore.Level
		message string
	}{
		{
			zap.DebugLevel,
			"task skipped",
		},
		{
			zap.DebugLevel,
			"task skipped",
		},
		{
			zap.DebugLevel,
			"taskflow skipped",
		},
	}

	// logs
	logEntries := observedLogs.All()
	assert.Equal(t, len(expected), len(logEntries))
	for i, entry := range logEntries {
		assert.Equal(t, expected[i].level, entry.Level)
		assert.Equal(t, expected[i].message, entry.Message)
		t.Logf("log entry - level: %q, message: %q, fields: %v", entry.Level, entry.Message, entry.Context)
	}
}

func TestInstrumentCancelledContext(t *testing.T) {
	scope := tally.NewTestScope("", nil)
	core, observedLogs := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()

	h := &H{Scope: scope, Logger: logger}
	_, err := h.Run(ctx, "1")
	assert.Error(t, err)

	// metrics
	counters := scope.Snapshot().Counters()
	for k := range counters {
		t.Logf("got counter with key %q", k)
	}
	assert.Equal(t, int64(1), counters["task.skipped+task=Atoi"].Value())
	assert.Equal(t, int64(1), counters["task.skipped+task=uint8"].Value())
	assert.Equal(t, int64(1), counters["taskflow.skipped+flow=AtoiRun"].Value())

	// logs
	expectedLevel := zap.DebugLevel
	expectedMessages := []string{
		"task skipped",
		"task skipped",
		"taskflow skipped",
	}
	logEntries := observedLogs.All()
	assert.Equal(t, len(expectedMessages), len(logEntries))
	for i, entry := range logEntries {
		assert.Equal(t, expectedLevel, entry.Level)
		assert.Equal(t, expectedMessages[i], entry.Message)
		t.Logf("log entry - level: %q, message: %q, fields: %v", entry.Level, entry.Message, entry.Context)
	}
}

func TestInstrumentRecover(t *testing.T) {
	scope := tally.NewTestScope("", nil)
	core, observedLogs := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	h := &H{Scope: scope, Logger: logger}
	ctx := context.Background()
	v, err := h.Run(ctx, "300")

	assert.NoError(t, err)
	assert.Equal(t, uint8(0), v)

	// metrics
	counters := scope.Snapshot().Counters()
	for k := range counters {
		t.Logf("got counter with key %q", k)
	}
	assert.Equal(t, int64(1), counters["task.success+task=Atoi"].Value())
	assert.Equal(t, int64(1), counters["task.error+task=uint8"].Value())
	assert.Equal(t, int64(1), counters["task.recovered+task=uint8"].Value())
	assert.Equal(t, int64(1), counters["taskflow.success+flow=AtoiRun"].Value())

	// logs
	expected := []struct {
		level   zapcore.Level
		message string
	}{
		{
			zap.DebugLevel,
			"task succeeded",
		},
		{
			zap.ErrorLevel,
			"task error recovered",
		},
		{
			zap.DebugLevel,
			"taskflow succeeded",
		},
	}
	logEntries := observedLogs.All()
	assert.Equal(t, len(expected), len(logEntries))
	for i, entry := range logEntries {
		assert.Equal(t, expected[i].level, entry.Level)
		assert.Equal(t, expected[i].message, entry.Message)
		t.Logf("log entry - level: %q, message: %q, fields: %v", entry.Level, entry.Message, entry.Context)
	}
}

func TestInstrumentAnnotationOrder(t *testing.T) {
	scope := tally.NewTestScope("", nil)
	core, observedLogs := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	h := &H{Scope: scope, Logger: logger}
	ctx := context.Background()
	v, err := h.Do(ctx, "1")

	assert.NoError(t, err)
	assert.Equal(t, 1, v)

	// metrics
	counters := scope.Snapshot().Counters()
	for k := range counters {
		t.Logf("got counter with key %q", k)
	}
	assert.Equal(t, int64(1), counters["task.success+task=Atoi"].Value())
	assert.Equal(t, int64(1), counters["taskflow.success+flow=AtoiDo"].Value())

	// logs
	expectedLevel := zap.DebugLevel
	expectedMessages := []string{
		"task succeeded",
		"taskflow succeeded",
	}
	logEntries := observedLogs.All()
	assert.Equal(t, len(expectedMessages), len(logEntries))
	for i, entry := range logEntries {
		assert.Equal(t, expectedLevel, entry.Level)
		assert.Equal(t, expectedMessages[i], entry.Message)
		t.Logf("log entry - level: %q, message: %q, fields: %v", entry.Level, entry.Message, entry.Context)
	}
}

func TestInstrumentTaskButNotFlow(t *testing.T) {
	scope := tally.NewTestScope("", nil)
	core, observedLogs := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	h := &H{Scope: scope, Logger: logger}
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
		"task succeeded",
	}
	logEntries := observedLogs.All()
	assert.Equal(t, len(expectedMessages), len(logEntries))
	for i, entry := range logEntries {
		assert.Equal(t, expectedLevel, entry.Level)
		assert.Equal(t, expectedMessages[i], entry.Message)
		t.Logf("log entry - level: %q, message: %q, fields: %v", entry.Level, entry.Message, entry.Context)
	}
}

// TestT3630161 tests against regression for T3630161
func TestT3630161(t *testing.T) {
	scope := tally.NewTestScope("", nil)
	logger := zaptest.NewLogger(t)
	h := &H{Scope: scope, Logger: logger}
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
	assert.Equal(t, map[string]string{"task": "End"}, countersByName["task.success"][0].Tags())
	assert.Equal(t, 1, len(countersByName["task.error"]))
	assert.Equal(t, map[string]string{"task": "Err"}, countersByName["task.error"][0].Tags())
	assert.Equal(t, 1, len(countersByName["task.recovered"]))
	assert.Equal(t, map[string]string{"task": "Err"}, countersByName["task.recovered"][0].Tags())
	assert.Equal(t, 1, len(countersByName["task.recovered"]))
	assert.Equal(t, map[string]string{"flow": "T3630161"}, countersByName["taskflow.success"][0].Tags())
}
