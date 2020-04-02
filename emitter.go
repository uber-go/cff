package cff

import (
	"context"
	"time"

	"github.com/uber-go/tally"
)

// FlowEmitter receives events for when flow events occur, for the purpose of
// emitting metrics.
//
// WARNING: This interface is not stable and may change in the future.
type FlowEmitter interface {
	// FlowSuccess is called when a flow runs successfully.
	FlowSuccess(context.Context)
	// FlowError is called when a flow fails due to a task error.
	FlowError(context.Context, error)
	// FlowSkipped is called when a flow fails due to a task error. Currently,
	// only adding to be backwards compatible. There is discussion in ERD to
	// remove this metric.
	FlowSkipped(context.Context, error)
	// FlowDone is called when a flow finishes.
	FlowDone(context.Context, time.Duration)
	// FlowFailedTask is called when a flow fails due to a task error and
	// returns a shallow copy of current FlowEmitter with updated tags.
	FlowFailedTask(ctx context.Context, task string, err error) FlowEmitter

	flowEmitter() // private interface
}

// TaskEmitter receives events for when task events occur, for the purpose of
// emitting metrics.
//
// WARNING: This interface is not stable and may change in the future.
type TaskEmitter interface {
	// TaskSuccess is called when a task runs successfully.
	TaskSuccess(context.Context)
	// TaskError is called when a task fails due to a task error.
	TaskError(context.Context, error)
	// TaskErrorRecovered is called when a task fails due to a task error
	// and recovers in a FallbackWith.
	TaskErrorRecovered(context.Context, error)
	// TaskSkipped is called when a task is skipped due to predicate or an
	// earlier task error.
	TaskSkipped(context.Context, error)
	// TaskPanic is called when a task panics.
	TaskPanic(context.Context, interface{})
	// TaskPanicRecovered is called when a task panics but is recovered by
	// a FallbackWith.
	TaskPanicRecovered(context.Context, interface{})
	// TaskDone is called when a task finishes.
	TaskDone(context.Context, time.Duration)

	taskEmitter() // private interface
}

// FlowInfo provides information to uniquely identify a flow.
type FlowInfo struct {
	Flow         string
	File         string
	Line, Column int
}

// TaskInfo provides information to uniquely identify a task.
type TaskInfo struct {
	Task         string
	File         string
	Line, Column int
}

// Emitter initializes Task and Flow emitters.
//
// WARNING: This interface is not stable and may change in the future.
type Emitter interface {
	// TaskInit returns a TaskEmitter which could be memoized based on task name.
	TaskInit(*TaskInfo, *FlowInfo) TaskEmitter
	// FlowInit returns a FlowEmitter which could be memoized based on flow name.
	FlowInit(*FlowInfo) FlowEmitter

	emitter() // private interface
}

// DefaultEmitter sets up default implementation of metrics used in the
// template with memoization of the scope.
func DefaultEmitter(scope tally.Scope) Emitter {
	return TallyEmitter(scope)
}
