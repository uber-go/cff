package cff

import (
	"time"
)

type metricsEmitterStack []MetricsEmitter

// MetricsEmitterStack allows users to combine multiple MetricsEmitter objects into a single one
// that sends events to all of them.
func MetricsEmitterStack(e []MetricsEmitter) MetricsEmitter {
	return metricsEmitterStack(e)
}

type metricsEmitterStackTask struct {
	task  string
	stack []TaskEmitter
}

// TaskInit returns a TaskEmitter which could be memoized based on task name.
func (s metricsEmitterStack) TaskInit(task string) TaskEmitter {
	emitters := make([]TaskEmitter, 0, len(s))
	for _, e := range s {
		emitters = append(emitters, e.TaskInit(task))
	}

	return &metricsEmitterStackTask{
		task:  task,
		stack: emitters,
	}
}

// TaskSuccess is called when a task runs successfully.
func (s *metricsEmitterStackTask) TaskSuccess() {
	for _, e := range s.stack {
		e.TaskSuccess()
	}
}

// TaskError is called when a task fails due to a task error.
func (s *metricsEmitterStackTask) TaskError() {
	for _, e := range s.stack {
		e.TaskError()
	}
}

// TaskSkipped is called when a task is skipped due to predicate or an
// earlier task error.
func (s *metricsEmitterStackTask) TaskSkipped() {
	for _, e := range s.stack {
		e.TaskSkipped()
	}
}

// TaskPanic is called when a task panics.
func (s *metricsEmitterStackTask) TaskPanic() {
	for _, e := range s.stack {
		e.TaskPanic()
	}
}

// TaskRecovered is called when a task errors but it was recovered by a
// RecoverWith annotation.
func (s *metricsEmitterStackTask) TaskRecovered() {
	for _, e := range s.stack {
		e.TaskRecovered()
	}
}

// TaskDone is called when a task finishes.
func (s *metricsEmitterStackTask) TaskDone(d time.Duration) {
	for _, e := range s.stack {
		e.TaskDone(d)
	}
}

type metricsEmitterStackFlow struct {
	flow  string
	stack []FlowEmitter
}

// FlowInit returns a FlowEmitter which could be memoized based on flow name.
func (s metricsEmitterStack) FlowInit(flow string) FlowEmitter {
	emitters := make([]FlowEmitter, 0, len(s))
	for _, e := range s {
		emitters = append(emitters, e.FlowInit(flow))
	}

	return &metricsEmitterStackFlow{
		flow:  flow,
		stack: emitters,
	}
}

// FlowSuccess is called when a flow runs successfully.
func (s *metricsEmitterStackFlow) FlowSuccess() {
	for _, e := range s.stack {
		e.FlowSuccess()
	}
}

// FlowError is called when a flow fails due to a task error.
func (s *metricsEmitterStackFlow) FlowError() {
	for _, e := range s.stack {
		e.FlowError()
	}
}

// FlowSkipped is called when a flow fails due to a task error. Currently,
// only adding to be backwards compatible. There is discussion in ERD to
// remove this metric.
func (s *metricsEmitterStackFlow) FlowSkipped() {
	for _, e := range s.stack {
		e.FlowSkipped()
	}
}

// FlowDone is called when a flow finishes.
func (s *metricsEmitterStackFlow) FlowDone(d time.Duration) {
	for _, e := range s.stack {
		e.FlowDone(d)
	}
}

// FlowFailedTask is called when a flow fails due to a task error and
// returns a shallow copy of current FlowEmitter with updated tags.
func (s *metricsEmitterStackFlow) FlowFailedTask(task string) FlowEmitter {
	emitters := make([]FlowEmitter, 0, len(s.stack))

	for _, e := range s.stack {
		emitters = append(emitters, e.FlowFailedTask(task))
	}

	return &metricsEmitterStackFlow{
		flow:  s.flow,
		stack: emitters,
	}
}
