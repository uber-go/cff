package panic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestCatchesPanicParallel(t *testing.T) {
	scope := tally.NewTestScope("", nil)
	logger, err := zap.NewDevelopment()
	p := &panicker{
		scope:  scope,
		logger: logger,
	}
	require.NoError(t, err)
	err = p.FlowPanicsParallel()
	require.Error(t, err)
}

func TestCatchesPanicSerial(t *testing.T) {
	scope := tally.NewTestScope("", nil)
	core, observedLogs := observer.New(zap.NewAtomicLevel())
	logger := zap.New(core)
	p := &panicker{
		scope:  scope,
		logger: logger,
	}
	err := p.FlowPanicsSerial()
	require.Error(t, err)
	snapshot := scope.Snapshot()
	counters := snapshot.Counters()
	assert.Equal(t, int64(1), counters["task.panic+"+tally.KeyForStringMap(map[string]string{"name": "T1"})].Value())
	assert.Equal(t, int64(1), counters["taskflow.error+"+tally.KeyForStringMap(map[string]string{"name": "FlowPanicsSerial"})].Value())
	logs := observedLogs.All()
	assert.Equal(t, "task panic", logs[0].Message)
	assert.Equal(t, "T1", logs[0].ContextMap()["name"])
	_, ok := logs[0].ContextMap()["stack"]
	assert.Equal(t, true, ok)
}
