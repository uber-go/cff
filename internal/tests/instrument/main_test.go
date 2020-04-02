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
		"task succeeded",
		"task succeeded",
		"taskflow succeeded",
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

func TestInstrumentWithLogFields(t *testing.T) {
	t.Run("Ellipsis", func(t *testing.T) {
		core, observedLogs := observer.New(zap.DebugLevel)
		logger := zap.New(core)
		h := &DefaultEmitter{Scope: tally.NoopScope, Logger: logger}

		_, err := h.Run(context.Background(), "1",
			zap.String("foo", "bar"), zap.Int("baz", 42))
		require.NoError(t, err)

		logEntries := observedLogs.All()
		assert.NotEmpty(t, logEntries)
		for _, entry := range logEntries {
			fields := entry.ContextMap()
			foo, ok := fields["foo"]
			assert.True(t, ok, "foo must be present in all log messages")
			assert.Equal(t, "bar", foo)

			baz, ok := fields["baz"]
			assert.True(t, ok, "baz must be present in all log messages")
			assert.EqualValues(t, 42, baz)
		}
	})

	t.Run("Explicit", func(t *testing.T) {
		core, observedLogs := observer.New(zap.DebugLevel)
		logger := zap.New(core)
		h := &DefaultEmitter{Scope: tally.NoopScope, Logger: logger}

		_, err := h.ExplicitListOfFields(context.Background(), "1")
		require.NoError(t, err)

		logEntries := observedLogs.All()
		assert.NotEmpty(t, logEntries)
		for _, entry := range logEntries {
			fields := entry.ContextMap()

			foo, ok := fields["foo"]
			assert.True(t, ok, "foo must be present in all log messages")
			assert.Equal(t, "bar", foo)

			baz, ok := fields["baz"]
			assert.True(t, ok, "baz must be present in all log messages")
			assert.EqualValues(t, 42, baz)
		}
	})
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
	assert.Equal(t, int64(1), counters["taskflow.error+failedtask=Atoi,flow=AtoiRun"].Value())

	expected := []struct {
		level   zapcore.Level
		message string
		fields  map[string]interface{}
	}{
		{
			zap.DebugLevel,
			"task skipped",
			map[string]interface{}{"task": "uint8"},
		},
		{
			zap.DebugLevel,
			"taskflow skipped",
			nil,
		},
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
	assert.Equal(t, int64(1), counters["task.error+flow=AtoiRun,task=uint8"].Value())
	assert.Equal(t, int64(1), counters["task.recovered+flow=AtoiRun,task=uint8"].Value())
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
	assert.Equal(t, 1, len(countersByName["task.error"]))
	assert.Equal(t, map[string]string{"flow": "T3630161", "task": "Err"}, countersByName["task.error"][0].Tags())
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
			"task succeeded",
			"taskflow skipped",
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
			"task succeeded",
			"task succeeded",
			"taskflow succeeded",
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
			"task succeeded",
			"taskflow succeeded",
			"task skipped",
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

// logEvolutionTester helps verify that a new code path behaves similarly to
// an older code path in terms of logged information.
//
// Logs emitted by the new code path must have at least as much information as
// the old code path.
//
//   lev := newLogEvolutionTester(t, zapcore.DebugLevel)
//   defer lev.Verify()
//
//   lev.Old().Info("foo")
//   lev.New().Info("foo", zap.String("bar", "baz"))
//
//
type logEvolutionTester struct {
	t                *testing.T
	oldCore, newCore zapcore.Core
	oldLogs, newLogs *observer.ObservedLogs

	skipFields []string
}

func newLogEvolutionTester(lvl zapcore.Level) *logEvolutionTester {
	oldCore, oldLogs := observer.New(lvl)
	newCore, newLogs := observer.New(lvl)
	return &logEvolutionTester{
		oldCore: oldCore,
		oldLogs: oldLogs,
		newCore: newCore,
		newLogs: newLogs,
	}
}

// SkipFields informs the logEvolutionTester that it should not verify
// equality of the provided fields in log messages.
func (lev *logEvolutionTester) SkipFields(names ...string) {
	lev.skipFields = append(lev.skipFields, names...)
}

// Old returns the logger for the old code path.
func (lev *logEvolutionTester) Old() *zap.Logger { return zap.New(lev.oldCore) }

// New returns the logger for the new code path.
func (lev *logEvolutionTester) New() *zap.Logger { return zap.New(lev.newCore) }

// Verify confirms that log entries emitted on the new code path match those
// on the old code path, optionally with new fields.
func (lev *logEvolutionTester) Verify(t *testing.T) {

	oldLogs := lev.oldLogs.AllUntimed()
	newLogs := lev.newLogs.AllUntimed()

	if !assert.Equal(t, len(oldLogs), len(newLogs), "number of log entries") {
		t.Log("Old:")
		for _, e := range oldLogs {
			t.Logf("\t%v\t%#v", e.Entry, e.ContextMap())
		}

		t.Log("New:")
		for _, e := range newLogs {
			t.Logf("\t%v\t%#v", e.Entry, e.ContextMap())
		}
		return
	}

	skipFields := make(map[string]struct{})
	for _, f := range lev.skipFields {
		skipFields[f] = struct{}{}
	}

	for i, oldEntry := range oldLogs {
		newEntry := newLogs[i]
		if !assert.Equalf(t, oldEntry.Entry, newEntry.Entry, "entry %v", i) {
			continue
		}

		oldContext := oldEntry.ContextMap()
		newContext := newEntry.ContextMap()
		for k, oldv := range oldContext {
			if _, ok := skipFields[k]; ok {
				continue
			}
			assert.Equalf(t, oldv, newContext[k], "field %q of entry %v", k, i)
		}
	}
}

func TestLogEmitterDefaultLoggerParity(t *testing.T) {
	// This test verifies that the LogEmitter generates the same logs as
	// the CFF2-generated code.
	//
	// This test may be deleted once CFF2-generated code is transitioned
	// to using the LogEmitter.

	tests := []struct {
		desc string
		call func(context.Context, *CustomEmitter)
	}{
		{
			desc: "Run",
			call: func(ctx context.Context, e *CustomEmitter) {
				e.Run(ctx, "42")
			},
		},
		{
			desc: "Run/error",
			call: func(ctx context.Context, e *CustomEmitter) {
				e.Run(ctx, "5000") // out of range
			},
		},
		{
			desc: "Do",
			call: func(ctx context.Context, e *CustomEmitter) {
				e.Do(ctx, "42")
			},
		},
		{
			desc: "Do/error",
			call: func(ctx context.Context, e *CustomEmitter) {
				e.Do(ctx, "not a number")
			},
		},
		{
			desc: "Work",
			call: func(ctx context.Context, e *CustomEmitter) {
				e.Work(ctx, "42")
			},
		},
		{
			desc: "Work/error",
			call: func(ctx context.Context, e *CustomEmitter) {
				e.Work(ctx, "not a number")
			},
		},
		{
			desc: "T3630161",
			call: func(ctx context.Context, e *CustomEmitter) { e.T3630161(ctx) },
		},
		{
			desc: "T3795761",
			call: func(ctx context.Context, e *CustomEmitter) {
				e.T3795761(ctx, true /* run */, false /* error */)
			},
		},
		{
			desc: "T3795761/no run",
			call: func(ctx context.Context, e *CustomEmitter) {
				e.T3795761(ctx, false /* run */, false /* error */)
			},
		},
		{
			desc: "T3795761/error",
			call: func(ctx context.Context, e *CustomEmitter) {
				e.T3795761(ctx, true /* run */, true /* error */)
			},
		},
		{
			desc: "FlowAlwaysPanics",
			call: func(ctx context.Context, e *CustomEmitter) { e.FlowAlwaysPanics(ctx) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			lev := newLogEvolutionTester(zapcore.DebugLevel)
			lev.SkipFields("stack") // don't verify stack trace equality
			defer lev.Verify(t)

			ctx := context.Background()

			tt.call(ctx, &CustomEmitter{
				Scope:   tally.NoopScope,
				Logger:  lev.Old(),
				Emitter: cff.LogEmitter(zap.NewNop()),
			})

			tt.call(ctx, &CustomEmitter{
				Scope:   tally.NoopScope,
				Logger:  zap.NewNop(),
				Emitter: cff.LogEmitter(lev.New()),
			})
		})
	}

}
