package instrument

import (
	"context"
	"testing"

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
	h := &h{scope: scope, logger: logger}
	ctx := context.Background()
	v, err := h.run(ctx, "1")

	assert.NoError(t, err)
	assert.Equal(t, uint8(1), v)

	// metrics
	counters := scope.Snapshot().Counters()
	for k := range counters {
		t.Logf("got counter with key %q", k)
	}
	assert.Equal(t, int64(1), counters["task.success+name=Atoi"].Value())
	assert.Equal(t, int64(1), counters["task.success+name=uint8"].Value())
	assert.Equal(t, int64(1), counters["taskflow.success+name=AtoiRun"].Value())

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
	h := &h{scope: scope, logger: logger}
	ctx := context.Background()
	_, err := h.run(ctx, "NaN")

	assert.Error(t, err)

	// metrics
	counters := scope.Snapshot().Counters()
	for k := range counters {
		t.Logf("got counter with key %q", k)
	}
	assert.Equal(t, int64(1), counters["task.error+name=Atoi"].Value())
	assert.Equal(t, int64(1), counters["taskflow.error+failedTask=Atoi,name=AtoiRun"].Value())

	// logs
	logEntries := observedLogs.All()
	assert.Equal(t, 0, len(logEntries))
}

func TestInstrumentCancelledContext(t *testing.T) {
	scope := tally.NewTestScope("", nil)
	core, observedLogs := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()

	h := &h{scope: scope, logger: logger}
	_, err := h.run(ctx, "1")
	assert.Error(t, err)

	// metrics
	counters := scope.Snapshot().Counters()
	for k := range counters {
		t.Logf("got counter with key %q", k)
	}
	assert.Equal(t, int64(1), counters["task.skipped+name=Atoi"].Value())
	assert.Equal(t, int64(1), counters["task.skipped+name=uint8"].Value())
	assert.Equal(t, int64(1), counters["taskflow.skipped+name=AtoiRun"].Value())

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
	h := &h{scope: scope, logger: logger}
	ctx := context.Background()
	v, err := h.run(ctx, "300")

	assert.NoError(t, err)
	assert.Equal(t, uint8(0), v)

	// metrics
	counters := scope.Snapshot().Counters()
	for k := range counters {
		t.Logf("got counter with key %q", k)
	}
	assert.Equal(t, int64(1), counters["task.success+name=Atoi"].Value())
	assert.Equal(t, int64(1), counters["task.error+name=uint8"].Value())
	assert.Equal(t, int64(1), counters["task.recovered+name=uint8"].Value())
	assert.Equal(t, int64(1), counters["taskflow.success+failedTask=uint8,name=AtoiRun"].Value())

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
	h := &h{scope: scope, logger: logger}
	ctx := context.Background()
	v, err := h.do(ctx, "1")

	assert.NoError(t, err)
	assert.Equal(t, 1, v)

	// metrics
	counters := scope.Snapshot().Counters()
	for k := range counters {
		t.Logf("got counter with key %q", k)
	}
	assert.Equal(t, int64(1), counters["task.success+name=Atoi"].Value())
	assert.Equal(t, int64(1), counters["taskflow.success+name=AtoiDo"].Value())

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
	h := &h{scope: scope, logger: logger}
	ctx := context.Background()
	v, err := h.work(ctx, "1")

	assert.NoError(t, err)
	assert.Equal(t, 1, v)

	// metrics
	counters := scope.Snapshot().Counters()
	for k := range counters {
		t.Logf("got counter with key %q", k)
	}
	assert.Equal(t, int64(1), counters["task.success+name=Atoi"].Value())

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
