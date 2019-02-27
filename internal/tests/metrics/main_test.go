package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"testing"
)

func TestMetrics(t *testing.T) {
	scope := tally.NewTestScope("", nil)
	h := &h{scope: scope}
	ctx := context.Background()
	v, err := h.run(ctx, "1")
	assert.NoError(t, err)
	assert.Equal(t, 1, v)
	counters := scope.Snapshot().Counters()
	for k := range counters {
		t.Logf("got counter with key %q", k)
	}
	assert.Equal(t, int64(1), counters["task.success+"].Value())
	assert.Equal(t, int64(1), counters["taskflow.success+"].Value())
}

func TestMetricsError(t *testing.T) {
	scope := tally.NewTestScope("", nil)
	ctx := context.Background()
	h := &h{scope: scope}
	_, err := h.run(ctx, "NaN")
	assert.Error(t, err)
	counters := scope.Snapshot().Counters()
	for k := range counters {
		t.Logf("got counter with key %q", k)
	}
	assert.Equal(t, int64(1), counters["task.error+"].Value())
	assert.Equal(t, int64(1), counters["taskflow.error+"].Value())
}

func TestMetricsCancelledContext(t *testing.T) {
	scope := tally.NewTestScope("", nil)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()

	h := &h{scope: scope}
	_, err := h.run(ctx, "1")
	assert.Error(t, err)
	counters := scope.Snapshot().Counters()
	for k := range counters {
		t.Logf("got counter with key %q", k)
	}
	assert.Equal(t, int64(1), counters["task.skipped+"].Value())
	assert.Equal(t, int64(1), counters["taskflow.skipped+"].Value())
}
