package cff

import (
	"context"
	"time"
)

// NopEmitter is a cff emitter that does not do anything.
func NopEmitter() Emitter {
	// We implement the interface on the pointer receiver to avoid
	// allocations when the emitter is used. Conversion from pointer to
	// interface requires no allocs, but conversion from value to
	// interface does.
	return &nopEmitter{}
}

// NopFlowEmitter is a Flow emitter that does not do anything.
func NopFlowEmitter() FlowEmitter {
	return &nopEmitter{}
}

// NopParallelEmitter is a Parallel emitter that does not do anything.
func NopParallelEmitter() ParallelEmitter {
	return &nopEmitter{}
}

// NopTaskEmitter is a Task emitter that does not do anything.
func NopTaskEmitter() TaskEmitter {
	return &nopEmitter{}
}

type nopEmitter struct{}

func (e *nopEmitter) FlowInit(*FlowInfo) FlowEmitter { return e }

func (*nopEmitter) FlowSuccess(context.Context) {}

func (*nopEmitter) FlowError(context.Context, error) {}

func (*nopEmitter) FlowDone(context.Context, time.Duration) {}

func (e *nopEmitter) ParallelInit(*ParallelInfo) ParallelEmitter { return e }

func (*nopEmitter) ParallelSuccess(context.Context) {}

func (*nopEmitter) ParallelError(context.Context, error) {}

func (*nopEmitter) ParallelDone(context.Context, time.Duration) {}

func (e *nopEmitter) TaskInit(*TaskInfo, *DirectiveInfo) TaskEmitter { return e }

func (*nopEmitter) TaskSuccess(context.Context) {}

func (*nopEmitter) TaskError(context.Context, error) {}

func (*nopEmitter) TaskErrorRecovered(context.Context, error) {}

func (*nopEmitter) TaskSkipped(context.Context, error) {}

func (*nopEmitter) TaskPanic(context.Context, interface{}) {}

func (*nopEmitter) TaskPanicRecovered(context.Context, interface{}) {}

func (*nopEmitter) TaskDone(context.Context, time.Duration) {}

func (e *nopEmitter) SchedulerInit(*SchedulerInfo) SchedulerEmitter { return e }

func (e *nopEmitter) EmitScheduler(SchedulerState) {}
