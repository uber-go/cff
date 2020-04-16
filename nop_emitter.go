package cff

import (
	"context"
	"time"
)

// NopEmitter is a CFF2 emitter that does not do anything.
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

// NopTaskEmitter is a Task emitter that does not do anything.
func NopTaskEmitter() TaskEmitter {
	return &nopEmitter{}
}

type nopEmitter struct{}

func (*nopEmitter) emitter()     {}
func (*nopEmitter) flowEmitter() {}
func (*nopEmitter) taskEmitter() {}

func (e *nopEmitter) FlowInit(*FlowInfo) FlowEmitter { return e }

func (*nopEmitter) FlowSuccess(context.Context) {}

func (*nopEmitter) FlowError(context.Context, error) {}

func (*nopEmitter) FlowDone(context.Context, time.Duration) {}

func (e *nopEmitter) TaskInit(*TaskInfo, *FlowInfo) TaskEmitter { return e }

func (*nopEmitter) TaskSuccess(context.Context) {}

func (*nopEmitter) TaskError(context.Context, error) {}

func (*nopEmitter) TaskErrorRecovered(context.Context, error) {}

func (*nopEmitter) TaskSkipped(context.Context, error) {}

func (*nopEmitter) TaskPanic(context.Context, interface{}) {}

func (*nopEmitter) TaskPanicRecovered(context.Context, interface{}) {}

func (*nopEmitter) TaskDone(context.Context, time.Duration) {}
