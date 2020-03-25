package cff

import (
	"context"
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
	// TaskSkipped is called when a task is skipped due to predicate or an
	// earlier task error.
	TaskSkipped(context.Context, error)
	// TaskPanic is called when a task panics.
	TaskPanic(context.Context, interface{})
	// TaskRecovered is called when a task errors but it was recovered by a
	// RecoverWith annotation.
	TaskRecovered(context.Context, interface{})
	// TaskDone is called when a task finishes.
	TaskDone(context.Context, time.Duration)
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

// MetricsEmitter initializes Task and Flow metrics emitters.
//
// WARNING: This interface is not stable and may change in the future.
type MetricsEmitter interface {
	// TaskInit returns a TaskEmitter which could be memoized based on task name.
	TaskInit(*TaskInfo, *FlowInfo) TaskEmitter
	// FlowInit returns a FlowEmitter which could be memoized based on flow name.
	FlowInit(*FlowInfo) FlowEmitter
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
func (e *taskEmitter) TaskError(context.Context, error) {
	e.scope.Counter("task.error").Inc(1)
}

func (e *taskEmitter) TaskPanic(context.Context, interface{}) {
	e.scope.Counter("task.panic").Inc(1)
}

func (e *taskEmitter) TaskRecovered(context.Context, interface{}) {
	e.scope.Counter("task.recovered").Inc(1)
}

func (e *taskEmitter) TaskSkipped(context.Context, error) {
	e.scope.Counter("task.skipped").Inc(1)
}

func (e *taskEmitter) TaskSuccess(context.Context) {
	e.scope.Counter("task.success").Inc(1)
}

func (e *taskEmitter) TaskDone(_ context.Context, d time.Duration) {
	e.scope.Timer("task.timing").Record(d)
}

// FlowEmitter implementation.
//
func (e *flowEmitter) FlowError(context.Context, error) {
	e.scope.Counter("taskflow.error").Inc(1)
}

func (e *flowEmitter) FlowSkipped(context.Context, error) {
	e.scope.Counter("taskflow.skipped").Inc(1)
}

func (e *flowEmitter) FlowSuccess(context.Context) {
	e.scope.Counter("taskflow.success").Inc(1)
}

func (e *flowEmitter) FlowFailedTask(_ context.Context, task string, _ error) FlowEmitter {
	return &flowEmitter{
		scope: e.scope.Tagged(map[string]string{
			"failedtask": task,
		})}
}

func (e *flowEmitter) FlowDone(_ context.Context, d time.Duration) {
	e.scope.Timer("taskflow.timing").Record(d)
}

// cacheKey uniquely identifies a task or a flow based on the position information.
type cacheKey struct {
	TaskName             string // name of the task
	TaskFile             string // file where task is defined
	TaskLine, TaskColumn int    // line and column in the file where the task is defined
	FlowName             string // name of the flow
	FlowFile             string // file where flow is defined
	FlowLine, FlowColumn int    // line and column in the file where the flow is defined
}

// MetricsEmitter implementation.
//
func (e *emitter) TaskInit(taskInfo *TaskInfo, flowInfo *FlowInfo) TaskEmitter {
	cacheKey := cacheKey{
		TaskName:   taskInfo.Task,
		TaskFile:   taskInfo.File,
		TaskLine:   taskInfo.Line,
		TaskColumn: taskInfo.Column,
		FlowName:   flowInfo.Flow,
		FlowFile:   flowInfo.File,
		FlowLine:   flowInfo.Line,
		FlowColumn: flowInfo.Column,
	}
	// Note: this lookup is an optimization to avoid the expensive Tagged call.
	if v, ok := e.tasks.Load(cacheKey); ok {
		return v.(TaskEmitter)
	}
	tags := map[string]string{
		"task": taskInfo.Task,
	}
	if flowInfo.Flow != "" {
		tags["flow"] = flowInfo.Flow
	}

	scope := e.scope.Tagged(tags)
	te := &taskEmitter{
		scope: scope,
	}
	v, _ := e.tasks.LoadOrStore(cacheKey, te)

	return v.(TaskEmitter)
}

func (e *emitter) FlowInit(info *FlowInfo) FlowEmitter {
	cacheKey := cacheKey{
		FlowName:   info.Flow,
		FlowFile:   info.File,
		FlowLine:   info.Line,
		FlowColumn: info.Column,
	}
	// Note: this lookup is an optimization to avoid the expensive Tagged call.
	if v, ok := e.flows.Load(cacheKey); ok {
		return v.(FlowEmitter)
	}
	scope := e.scope.Tagged(map[string]string{"flow": info.Flow})
	fe := &flowEmitter{
		scope: scope,
	}
	v, _ := e.flows.LoadOrStore(cacheKey, fe)

	return v.(FlowEmitter)
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
