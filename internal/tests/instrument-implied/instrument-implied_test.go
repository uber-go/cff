package instrumentimplied

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestInstrumentImpliedName(t *testing.T) {
	scope := tally.NewTestScope("", nil)
	core, observedLogs := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	h := &H{Scope: scope, Logger: logger}
	ctx := context.Background()
	_, err := h.ImpliedName(ctx, "1")

	assert.NoError(t, err)

	// metrics
	counters := scope.Snapshot().Counters()
	for k := range counters {
		t.Logf("got counter with key %q", k)
	}
	assert.Equal(t, int64(1), counters["task.success+flow=ImpliedName,task=instrument-implied.go.30"].Value())
	assert.Equal(t, int64(1), counters["task.success+flow=ImpliedName,task=instrument-implied.go.34"].Value())

	// logs
	logEntries := observedLogs.All()
	expectedLevel := zap.DebugLevel
	expectedMessages := []struct {
		message string
		fields  map[string]interface{}
	}{
		{"task success", map[string]interface{}{"flow": "ImpliedName", "task": "instrument-implied.go.30"}},
		{"task done", map[string]interface{}{"flow": "ImpliedName", "task": "instrument-implied.go.30"}},
		{"task success", map[string]interface{}{"flow": "ImpliedName", "task": "instrument-implied.go.34"}},
		{"task done", map[string]interface{}{"flow": "ImpliedName", "task": "instrument-implied.go.34"}},
		{"flow success", map[string]interface{}{"flow": "ImpliedName"}},
		{"flow done", map[string]interface{}{"flow": "ImpliedName"}},
	}
	for i, entry := range logEntries {
		t.Logf("log entry - level: %q, message: %q, fields: %v", entry.Level, entry.Message, entry.ContextMap())
		assert.Equal(t, expectedLevel, entry.Level)
		assert.Equal(t, expectedMessages[i].message, entry.Message)
		assert.Equal(t, expectedMessages[i].fields, entry.ContextMap())
	}
}
