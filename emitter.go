package cff

import (
	"context"
	"time"

	"go.uber.org/cff/scheduler"
)

//go:generate mockgen -destination mock_emitter_test.go -package cff go.uber.org/cff Emitter,TaskEmitter,FlowEmitter,ParallelEmitter,SchedulerEmitter

// Emitter initializes Task, Flow, and Parallel emitters.
type Emitter interface {
	// TaskInit returns a TaskEmitter which could be memoized based on task name.
	TaskInit(*TaskInfo, *DirectiveInfo) TaskEmitter
	// FlowInit returns a FlowEmitter which could be memoized based on flow name.
	FlowInit(*FlowInfo) FlowEmitter
	// ParallelInit returns a ParallelEmitter which could be memoized based on
	// parallel name.
	ParallelInit(*ParallelInfo) ParallelEmitter
	// SchedulerInit returns an emitter for the CFF scheduler.
	SchedulerInit(s *SchedulerInfo) SchedulerEmitter
}

// SchedulerState describes the status of jobs managed by the CFF scheduler.
type SchedulerState = scheduler.State

// SchedulerEmitter provides observability into the state of the CFF
// scheduler.
type SchedulerEmitter interface {
	// EmitScheduler emits the state of the CFF scheduler.
	EmitScheduler(s SchedulerState)
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

// DirectiveInfo provides information to uniquely identify a CFF Directive.
type DirectiveInfo struct {
	Name string
	// Directive is the type of directive (e.g flow or parallel)
	Directive    DirectiveType
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
type FlowEmitter interface {
	// FlowSuccess is called when a flow runs successfully.
	FlowSuccess(context.Context)
	// FlowError is called when a flow fails due to a task error.
	FlowError(context.Context, error)
	// FlowDone is called when a flow finishes.
	FlowDone(context.Context, time.Duration)
}

// ParallelEmitter receives events for when parallel events occur, for the
// purpose of emitting metrics.
type ParallelEmitter interface {
	// ParallelSuccess is called when a parallel runs successfully.
	ParallelSuccess(context.Context)
	// ParallelError is called when a parallel fails due to a task error.
	ParallelError(context.Context, error)
	// ParallelDone is called when a parallel finishes.
	ParallelDone(context.Context, time.Duration)
}

// TaskEmitter receives events for when task events occur, for the purpose of
// emitting metrics.
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
}
