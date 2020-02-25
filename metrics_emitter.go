package cff

import (
	"sync"
	"time"

	"github.com/uber-go/tally"
)

// FlowEmitter receives events for when flow events occur, for the purpose of
// emitting metrics.
//
// WARNING: This interface is not stable and may change in the future.
type FlowEmitter interface {
	// FlowSuccess is called when a flow runs successfully.
	FlowSuccess()
	// FlowError is called when a flow fails due to a task error.
	FlowError()
	// FlowSkipped is called when a flow fails due to a task error. Currently,
	// only adding to be backwards compatible. There is discussion in ERD to
	// remove this metric.
	FlowSkipped()
	// FlowDone is called when a flow finishes.
	FlowDone(time.Duration)
	// FlowFailedTask is called when a flow fails due to a task error and
	// returns a shallow copy of current FlowEmitter with updated tags.
	FlowFailedTask(task string) FlowEmitter
}

// TaskEmitter receives events for when task events occur, for the purpose of
// emitting metrics.
//
// WARNING: This interface is not stable and may change in the future.
type TaskEmitter interface {
	// TaskSuccess is called when a task runs successfully.
	TaskSuccess()
	// TaskError is called when a task fails due to a task error.
	TaskError()
	// TaskSkipped is called when a task is skipped due to predicate or an
	// earlier task error.
	TaskSkipped()
	// TaskPanic is called when a task panics.
	TaskPanic()
	// TaskRecovered is called when a task errors but it was recovered by a
	// RecoverWith annotation.
	TaskRecovered()
	// TaskDone is called when a task finishes.
	TaskDone(time.Duration)
}

// MetricsEmitter initializes Task and Flow metrics emitters.
//
// WARNING: This interface is not stable and may change in the future.
type MetricsEmitter interface {
	// TaskInit returns a TaskEmitter which could be memoized based on task name.
	TaskInit(task string) TaskEmitter
	// FlowInit returns a FlowEmitter which could be memoized based on flow name.
	FlowInit(flow string) FlowEmitter
	// unexported makes this interface private.
	unexported()
}

type flowEmitter struct {
	scope tally.Scope
}

type taskEmitter struct {
	scope tally.Scope
}

type emitter struct {
	scope tally.Scope

	flows *sync.Map // map[string]FlowEmitter
	tasks *sync.Map // map[string]TaskEmitter
}

// Task Emitter implementation.
//
func (e *taskEmitter) TaskError() {
	e.scope.Counter("task.error").Inc(1)
}

func (e *taskEmitter) TaskPanic() {
	e.scope.Counter("task.panic").Inc(1)
}

func (e *taskEmitter) TaskRecovered() {
	e.scope.Counter("task.recovered").Inc(1)
}

func (e *taskEmitter) TaskSkipped() {
	e.scope.Counter("task.skipped").Inc(1)
}

func (e *taskEmitter) TaskSuccess() {
	e.scope.Counter("task.success").Inc(1)
}

func (e *taskEmitter) TaskDone(d time.Duration) {
	e.scope.Timer("task.timing").Record(d)
}

// FlowEmitter implementation.
//
func (e *flowEmitter) FlowError() {
	e.scope.Counter("taskflow.error").Inc(1)
}

func (e *flowEmitter) FlowSkipped() {
	e.scope.Counter("taskflow.skipped").Inc(1)
}

func (e *flowEmitter) FlowSuccess() {
	e.scope.Counter("taskflow.success").Inc(1)
}

func (e *flowEmitter) FlowFailedTask(task string) FlowEmitter {
	return &flowEmitter{
		scope: e.scope.Tagged(map[string]string{
			"failedtask": task,
		})}
}

// MetricsEmitter implementation.
//
// TODO(T5108563): Improve lookup if scope tags are different in case there is
// a collision between task instrumentation names with those tasks being in
// different flows.
func (e *emitter) TaskInit(task string) TaskEmitter {
	if v, ok := e.tasks.Load(task); ok {
		return v.(TaskEmitter)
	}
	scope := e.scope.Tagged(map[string]string{"task": task})
	te := &taskEmitter{
		scope: scope,
	}
	e.tasks.LoadOrStore(task, te)

	return te
}

func (e *emitter) FlowInit(flow string) FlowEmitter {
	if v, ok := e.tasks.Load(flow); ok {
		return v.(FlowEmitter)
	}
	scope := e.scope.Tagged(map[string]string{"flow": flow})
	fe := &flowEmitter{
		scope: scope,
	}
	e.flows.LoadOrStore(flow, fe)

	return fe
}

func (e *emitter) unexported() {}

func (e *flowEmitter) FlowDone(d time.Duration) {
	e.scope.Timer("taskflow.timing").Record(d)
}

// DefaultMetricsEmitter sets up default implementation of metrics used in the
// template with memoization of the scope.
func DefaultMetricsEmitter(scope tally.Scope) MetricsEmitter {
	return &emitter{
		scope: scope,
		flows: new(sync.Map),
		tasks: new(sync.Map),
	}
}
