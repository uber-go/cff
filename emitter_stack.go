package cff

import (
	"context"
	"time"
)

type emitterStack []Emitter

// EmitterStack allows users to combine multiple Emitter objects into a single one
// that sends events to all of them.
func EmitterStack(e []Emitter) Emitter {
	return emitterStack(e)
}

type emitterStackTask struct {
	task  string
	stack []TaskEmitter
}

// TaskInit returns a TaskEmitter which could be memoized based on task name.
func (s emitterStack) TaskInit(taskInfo *TaskInfo, flowInfo *FlowInfo) TaskEmitter {
	emitters := make([]TaskEmitter, 0, len(s))
	for _, e := range s {
		emitters = append(emitters, e.TaskInit(taskInfo, flowInfo))
	}

	return &emitterStackTask{
		task:  taskInfo.Task,
		stack: emitters,
	}
}

// TaskSuccess is called when a task runs successfully.
func (s *emitterStackTask) TaskSuccess(ctx context.Context) {
	for _, e := range s.stack {
		e.TaskSuccess(ctx)
	}
}

// TaskError is called when a task fails due to a task error.
func (s *emitterStackTask) TaskError(ctx context.Context, err error) {
	for _, e := range s.stack {
		e.TaskError(ctx, err)
	}
}

// TaskError is called when a task fails due to a task error.
func (s *emitterStackTask) TaskErrorRecovered(ctx context.Context, err error) {
	for _, e := range s.stack {
		e.TaskErrorRecovered(ctx, err)
	}
}

// TaskSkipped is called when a task is skipped due to predicate or an
// earlier task error.
func (s *emitterStackTask) TaskSkipped(ctx context.Context, err error) {
	for _, e := range s.stack {
		e.TaskSkipped(ctx, err)
	}
}

// TaskPanic is called when a task panics.
func (s *emitterStackTask) TaskPanic(ctx context.Context, pv interface{}) {
	for _, e := range s.stack {
		e.TaskPanic(ctx, pv)
	}
}

// TaskRecovered is called when a task errors but it was recovered by a
// RecoverWith annotation.
func (s *emitterStackTask) TaskPanicRecovered(ctx context.Context, pv interface{}) {
	for _, e := range s.stack {
		e.TaskPanicRecovered(ctx, pv)
	}
}

// TaskDone is called when a task finishes.
func (s *emitterStackTask) TaskDone(ctx context.Context, d time.Duration) {
	for _, e := range s.stack {
		e.TaskDone(ctx, d)
	}
}

type emitterStackFlow struct {
	flow  string
	stack []FlowEmitter
}

// FlowInit returns a FlowEmitter which could be memoized based on flow name.
func (s emitterStack) FlowInit(info *FlowInfo) FlowEmitter {
	emitters := make([]FlowEmitter, 0, len(s))
	for _, e := range s {
		emitters = append(emitters, e.FlowInit(info))
	}

	return &emitterStackFlow{
		flow:  info.Flow,
		stack: emitters,
	}
}

// FlowSuccess is called when a flow runs successfully.
func (s *emitterStackFlow) FlowSuccess(ctx context.Context) {
	for _, e := range s.stack {
		e.FlowSuccess(ctx)
	}
}

// FlowError is called when a flow fails due to a task error.
func (s *emitterStackFlow) FlowError(ctx context.Context, err error) {
	for _, e := range s.stack {
		e.FlowError(ctx, err)
	}
}

// FlowSkipped is called when a flow fails due to a task error. Currently,
// only adding to be backwards compatible. There is discussion in ERD to
// remove this metric.
func (s *emitterStackFlow) FlowSkipped(ctx context.Context, err error) {
	for _, e := range s.stack {
		e.FlowSkipped(ctx, err)
	}
}

// FlowDone is called when a flow finishes.
func (s *emitterStackFlow) FlowDone(ctx context.Context, d time.Duration) {
	for _, e := range s.stack {
		e.FlowDone(ctx, d)
	}
}

// FlowFailedTask is called when a flow fails due to a task error and
// returns a shallow copy of current FlowEmitter with updated tags.
func (s *emitterStackFlow) FlowFailedTask(ctx context.Context, task string, err error) FlowEmitter {
	emitters := make([]FlowEmitter, 0, len(s.stack))

	for _, e := range s.stack {
		emitters = append(emitters, e.FlowFailedTask(ctx, task, err))
	}

	return &emitterStackFlow{
		flow:  s.flow,
		stack: emitters,
	}
}
