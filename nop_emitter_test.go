package cff

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNopEmitter(t *testing.T) {
	e := NopEmitter()

	ctx := context.Background()

	t.Run("flow", func(t *testing.T) {
		e := e.FlowInit(&FlowInfo{Flow: "foo"})

		e.FlowSuccess(ctx)
		e.FlowError(ctx, errors.New("great sadness"))

		e = e.FlowFailedTask(ctx, "foo", errors.New("bar"))

		e.FlowSkipped(ctx, errors.New("something went wrong"))
		e.FlowDone(ctx, 3*time.Second)
	})

	t.Run("task", func(t *testing.T) {
		e := e.TaskInit(&TaskInfo{Task: "foo"}, &FlowInfo{Flow: "bar"})

		e.TaskSuccess(ctx)
		e.TaskError(ctx, errors.New("great sadness"))
		e.TaskErrorRecovered(ctx, errors.New("not that bad"))
		e.TaskSkipped(ctx, errors.New("something went wrong"))
		e.TaskPanic(ctx, "you found a bug")
		e.TaskPanicRecovered(ctx, "you found a bug that wasn't that bad")
		e.TaskDone(ctx, time.Second)
	})
}
