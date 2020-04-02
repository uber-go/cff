package cff

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLogFlowEmitter_IncludesFlowName(t *testing.T) {
	core, observed := observer.New(zapcore.DebugLevel)

	em := LogEmitter(zap.New(core)).FlowInit(&FlowInfo{Flow: "myflow"})
	em.FlowSuccess(context.Background())
	em.FlowSkipped(context.Background(), errors.New("foo"))

	for _, logEntry := range observed.TakeAll() {
		fields := logEntry.ContextMap()
		assert.Equalf(t, "myflow", fields["flow"],
			"flow name expected in %#v", fields)
	}
}

func TestLogTaskEmitter(t *testing.T) {
	ctx := context.Background()
	core, observed := observer.New(zapcore.DebugLevel)
	emitter := LogEmitter(zap.New(core))
	tem := emitter.TaskInit(&TaskInfo{Task: "mytask"}, &FlowInfo{Flow: "myflow"})

	t.Run("includes task and flow name", func(t *testing.T) {
		tem.TaskSuccess(ctx)
		tem.TaskErrorRecovered(ctx, errors.New("great sadness"))

		for _, logEntry := range observed.TakeAll() {
			fields := logEntry.ContextMap()
			assert.Equalf(t, "mytask", fields["task"],
				"task name expected in %#v", fields)
			assert.Equalf(t, "myflow", fields["flow"],
				"flow name expected in %#v", fields)
		}
	})

	t.Run("panic with value", func(t *testing.T) {
		tem.TaskPanic(ctx, "foo")
		logs := observed.TakeAll()
		require.Len(t, logs, 1)
		assert.Equal(t, "task panic: foo", fmt.Sprint(logs[0].ContextMap()["error"]))
	})

	t.Run("panic with error", func(t *testing.T) {
		tem.TaskPanic(ctx, errors.New("great sadness"))
		logs := observed.TakeAll()
		require.Len(t, logs, 1)
		assert.Equal(t, "great sadness", fmt.Sprint(logs[0].ContextMap()["error"]))
	})

	t.Run("panic recovered with value", func(t *testing.T) {
		tem.TaskPanicRecovered(ctx, "foo")
		logs := observed.TakeAll()
		require.Len(t, logs, 1)
		assert.Equal(t, "task panic: foo", fmt.Sprint(logs[0].ContextMap()["error"]))
	})

	t.Run("panic recovered with error", func(t *testing.T) {
		tem.TaskPanicRecovered(ctx, errors.New("great sadness"))
		logs := observed.TakeAll()
		require.Len(t, logs, 1)
		assert.Equal(t, "great sadness", fmt.Sprint(logs[0].ContextMap()["error"]))
	})
}
