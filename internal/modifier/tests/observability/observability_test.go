package observability

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestInstrumentFlow(t *testing.T) {
	scope := tally.NewTestScope("", nil)

	core, observed := observer.New(zapcore.DebugLevel)
	logger := zap.New(core)
	fmt.Println(observed)

	res, err := InstrumentFlow(scope, logger)
	require.NoError(t, err)
	assert.Equal(t, res, int64(1))

	t.Run("flow scope", func(t *testing.T) {
		counters := scope.Snapshot().Counters()
		for k := range counters {
			t.Logf("got counter with key %q", k)
		}
		assert.Equal(t, int64(1), counters["taskflow.success+flow=Instrumented"].Value())

		timers := scope.Snapshot().Timers()
		assert.NotNil(t, timers["taskflow.timing+flow=Instrumented"])
	})

	t.Run("flow logs", func(t *testing.T) {
		expectedLevel := zap.DebugLevel
		expectedMessages := []string{
			"flow success",
			"flow done",
		}
		logEntries := observed.All()
		assert.Equal(t, len(expectedMessages), len(logEntries))
		for _, entry := range logEntries {
			t.Logf("log entry - level: %q, message: %q, fields: %v", entry.Level, entry.Message, entry.ContextMap())
		}
		for i, entry := range logEntries {
			assert.Equal(t, expectedLevel, entry.Level)
			assert.Equal(t, expectedMessages[i], entry.Message)
		}
	})
}
