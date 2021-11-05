package cff

import (
	"context"
	"time"
)

type emitterStack []Emitter

func (emitterStack) emitter() {}

// EmitterStack allows users to combine multiple Emitters together.
//
// Events are sent to the emitters in an unspecified order. Emitters should
// not assume the ordering of events.
func EmitterStack(emitters ...Emitter) Emitter {
	switch len(emitters) {
	case 0:
		return NopEmitter()
	case 1:
		return emitters[0]
	default:
		var stack emitterStack
		for _, e := range emitters {
			if s, ok := e.(emitterStack); ok {
				// Flatten nested stacks.
				stack = append(stack, s...)
			} else {
				stack = append(stack, e)
			}
		}
		return stack
	}
}

type taskEmitterStack []TaskEmitter

func (taskEmitterStack) taskEmitter() {}

// TaskInit returns a TaskEmitter which could be memoized based on task name.
func (es emitterStack) TaskInit(taskInfo *TaskInfo, dInfo *DirectiveInfo) TaskEmitter {
	emitters := make(taskEmitterStack, 0, len(es))
	for _, e := range es {
		emitters = append(emitters, e.TaskInit(taskInfo, dInfo))
	}
	return emitters
}

// TaskSuccess is called when a task runs successfully.
func (ts taskEmitterStack) TaskSuccess(ctx context.Context) {
	for _, e := range ts {
		e.TaskSuccess(ctx)
	}
}

// TaskError is called when a task fails due to a task error.
func (ts taskEmitterStack) TaskError(ctx context.Context, err error) {
	for _, e := range ts {
		e.TaskError(ctx, err)
	}
}

// TaskError is called when a task fails due to a task error.
func (ts taskEmitterStack) TaskErrorRecovered(ctx context.Context, err error) {
	for _, e := range ts {
		e.TaskErrorRecovered(ctx, err)
	}
}

// TaskSkipped is called when a task is skipped due to predicate or an
// earlier task error.
func (ts taskEmitterStack) TaskSkipped(ctx context.Context, err error) {
	for _, e := range ts {
		e.TaskSkipped(ctx, err)
	}
}

// TaskPanic is called when a task panics.
func (ts taskEmitterStack) TaskPanic(ctx context.Context, pv interface{}) {
	for _, e := range ts {
		e.TaskPanic(ctx, pv)
	}
}

// TaskRecovered is called when a task errors but it was recovered by a
// RecoverWith annotation.
func (ts taskEmitterStack) TaskPanicRecovered(ctx context.Context, pv interface{}) {
	for _, e := range ts {
		e.TaskPanicRecovered(ctx, pv)
	}
}

// TaskDone is called when a task finishes.
func (ts taskEmitterStack) TaskDone(ctx context.Context, d time.Duration) {
	for _, e := range ts {
		e.TaskDone(ctx, d)
	}
}

type flowEmitterStack []FlowEmitter

func (flowEmitterStack) flowEmitter() {}

// FlowInit returns a FlowEmitter which could be memoized based on flow name.
func (es emitterStack) FlowInit(info *FlowInfo) FlowEmitter {
	emitters := make(flowEmitterStack, 0, len(es))
	for _, e := range es {
		emitters = append(emitters, e.FlowInit(info))
	}

	return emitters
}

// FlowSuccess is called when a flow runs successfully.
func (fs flowEmitterStack) FlowSuccess(ctx context.Context) {
	for _, e := range fs {
		e.FlowSuccess(ctx)
	}
}

// FlowError is called when a flow fails due to a task error.
func (fs flowEmitterStack) FlowError(ctx context.Context, err error) {
	for _, e := range fs {
		e.FlowError(ctx, err)
	}
}

// FlowDone is called when a flow finishes.
func (fs flowEmitterStack) FlowDone(ctx context.Context, d time.Duration) {
	for _, e := range fs {
		e.FlowDone(ctx, d)
	}
}

type parallelEmitterStack []ParallelEmitter

func (parallelEmitterStack) parallelEmitter() {}

// ParallelInit returns a ParallelEmitter which could be memoized based on parallel name.
func (es emitterStack) ParallelInit(info *ParallelInfo) ParallelEmitter {
	emitters := make(parallelEmitterStack, 0, len(es))
	for _, e := range es {
		emitters = append(emitters, e.ParallelInit(info))
	}
	return emitters
}

// ParallelSuccess is called when a parallel runs successfully.
func (ps parallelEmitterStack) ParallelSuccess(ctx context.Context) {
	for _, e := range ps {
		e.ParallelSuccess(ctx)
	}
}

// ParallelError is called when a parallel fails due to a task error.
func (ps parallelEmitterStack) ParallelError(ctx context.Context, err error) {
	for _, e := range ps {
		e.ParallelError(ctx, err)
	}
}

// ParallelDone is called when a parallel finishes.
func (ps parallelEmitterStack) ParallelDone(ctx context.Context, d time.Duration) {
	for _, e := range ps {
		e.ParallelDone(ctx, d)
	}
}

// SchedulerInit builds a SchedulerEmitter backed by the SchedulerEmitters of
// the underlying Emitters.
func (es emitterStack) SchedulerInit(info *SchedulerInfo) SchedulerEmitter {
	emitters := make(schedulerEmitterStack, len(es))
	for i, e := range es {
		emitters[i] = e.SchedulerInit(info)
	}
	return emitters
}

type schedulerEmitterStack []SchedulerEmitter

func (schedulerEmitterStack) schedulerEmitter() {}

// EmitScheduler emits the state of the CFF scheduler.
func (ses schedulerEmitterStack) EmitScheduler(s SchedulerState) {
	for _, e := range ses {
		e.EmitScheduler(s)
	}
}
