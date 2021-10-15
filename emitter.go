package cff

import (
	"context"
	"time"

	"go.uber.org/cff/scheduler"
)

// Emitter initializes Task and Flow emitters.
//
// WARNING: This interface is not stable and may change in the future.
type Emitter interface {
	// TaskInit returns a TaskEmitter which could be memoized based on task name.
	TaskInit(*TaskInfo, *FlowInfo) TaskEmitter
	// FlowInit returns a FlowEmitter which could be memoized based on flow name.
	FlowInit(*FlowInfo) FlowEmitter
	// SchedulerInit returns an emitter for the CFF scheduler.
	SchedulerInit(s *SchedulerInfo) SchedulerEmitter

	emitter() // private interface (GO-258).
}

// SchedulerState describes the status of jobs managed by the CFF scheduler.
type SchedulerState = scheduler.State

// SchedulerEmitter provides observability into the state of the CFF
// scheduler.
type SchedulerEmitter interface {
	// EmitScheduler emits the state of the CFF scheduler.
	EmitScheduler(s SchedulerState)

	schedulerEmitter() // private interface (GO-258).
}

// SchedulerInfo provides information about the context the scheduler
// is running in.
type SchedulerInfo struct {
	// Name of the directive the scheduler runs tasks for.
	Name string
	// DirectiveType is the type of Directive scheduler is running for
	// (e.g. flow, parallel).
	Directive    DirectiveType
	File         string
	Line, Column int
}

// FlowInfo provides information to uniquely identify a flow.
type FlowInfo struct {
	Name         string
	File         string
	Line, Column int
}

// TaskInfo provides information to uniquely identify a task.
type TaskInfo struct {
	Name         string
	File         string
	Line, Column int
}

// ParallelInfo provides information to uniquely identify a Parallel operation.
type ParallelInfo struct {
	Name         string
	File         string
	Line, Column int
}

// FlowEmitter receives events for when flow events occur, for the purpose of
// emitting metrics.
//
// WARNING: This interface is not stable and may change in the future.
type FlowEmitter interface {
	// FlowSuccess is called when a flow runs successfully.
	FlowSuccess(context.Context)
	// FlowError is called when a flow fails due to a task error.
	FlowError(context.Context, error)
	// FlowDone is called when a flow finishes.
	FlowDone(context.Context, time.Duration)

	flowEmitter() // private interface (GO-258).
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

	taskEmitter() // private interface (GO-258).
}
