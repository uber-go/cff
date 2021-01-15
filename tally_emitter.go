package cff

import (
	"context"
	"sync"
	"time"

	"github.com/uber-go/tally"
)

// cacheKey uniquely identifies a task or a flow based on the position information.
type cacheKey struct {
	TaskName             string // name of the task
	TaskFile             string // file where task is defined
	TaskLine, TaskColumn int    // line and column in the file where the task is defined
	FlowName             string // name of the flow
	FlowFile             string // file where flow is defined
	FlowLine, FlowColumn int    // line and column in the file where the flow is defined
}

type tallyEmitter struct {
	scope tally.Scope

	flows  *sync.Map // map[cacheKey]FlowEmitter
	tasks  *sync.Map // map[cacheKey]TaskEmitter
	scheds *sync.Map // map[cacheKey]SchedulerEmitter
}

func (tallyEmitter) emitter() {}

// TallyEmitter is a CFF2 emitter that emits metrics to Tally.
//
// A full list of metrics published by TallyEmitter can be found at
// https://eng.uberinternal.com/docs/cff2/observability/#metrics.
func TallyEmitter(scope tally.Scope) Emitter {
	return &tallyEmitter{
		scope:  scope,
		flows:  new(sync.Map),
		tasks:  new(sync.Map),
		scheds: new(sync.Map),
	}
}

func (e *tallyEmitter) TaskInit(taskInfo *TaskInfo, flowInfo *FlowInfo) TaskEmitter {
	cacheKey := cacheKey{
		TaskName:   taskInfo.Name,
		TaskFile:   taskInfo.File,
		TaskLine:   taskInfo.Line,
		TaskColumn: taskInfo.Column,
		FlowName:   flowInfo.Name,
		FlowFile:   flowInfo.File,
		FlowLine:   flowInfo.Line,
		FlowColumn: flowInfo.Column,
	}
	// Note: this lookup is an optimization to avoid the expensive Tagged call.
	if v, ok := e.tasks.Load(cacheKey); ok {
		return v.(TaskEmitter)
	}
	tags := map[string]string{
		"task": taskInfo.Name,
	}
	if flowInfo.Name != "" {
		tags["flow"] = flowInfo.Name
	}

	scope := e.scope.Tagged(tags)
	te := &tallyTaskEmitter{
		scope: scope,
	}
	v, _ := e.tasks.LoadOrStore(cacheKey, te)

	return v.(TaskEmitter)
}

func (e *tallyEmitter) FlowInit(info *FlowInfo) FlowEmitter {
	cacheKey := cacheKey{
		FlowName:   info.Name,
		FlowFile:   info.File,
		FlowLine:   info.Line,
		FlowColumn: info.Column,
	}
	// Note: this lookup is an optimization to avoid the expensive Tagged call.
	if v, ok := e.flows.Load(cacheKey); ok {
		return v.(FlowEmitter)
	}
	scope := e.scope.Tagged(map[string]string{"flow": info.Name})
	fe := &tallyFlowEmitter{
		scope: scope,
	}
	v, _ := e.flows.LoadOrStore(cacheKey, fe)

	return v.(FlowEmitter)
}

type tallyFlowEmitter struct {
	scope tally.Scope
}

func (tallyFlowEmitter) flowEmitter() {}

func (e *tallyFlowEmitter) FlowError(context.Context, error) {
	e.scope.Counter("taskflow.error").Inc(1)
}

func (e *tallyFlowEmitter) FlowSuccess(context.Context) {
	e.scope.Counter("taskflow.success").Inc(1)
}

func (e *tallyFlowEmitter) FlowDone(_ context.Context, d time.Duration) {
	e.scope.Timer("taskflow.timing").Record(d)
}

type tallyTaskEmitter struct {
	scope tally.Scope
}

func (tallyTaskEmitter) taskEmitter() {}

func (e *tallyTaskEmitter) TaskError(context.Context, error) {
	e.scope.Counter("task.error").Inc(1)
}

func (e *tallyTaskEmitter) TaskErrorRecovered(_ context.Context, err error) {
	e.scope.Counter("task.recovered").Inc(1)
}

func (e *tallyTaskEmitter) TaskPanic(_ context.Context, x interface{}) {
	e.scope.Counter("task.panic").Inc(1)
}

func (e *tallyTaskEmitter) TaskPanicRecovered(_ context.Context, x interface{}) {
	e.scope.Counter("task.recovered").Inc(1)
}

func (e *tallyTaskEmitter) TaskSkipped(context.Context, error) {
	e.scope.Counter("task.skipped").Inc(1)
}

func (e *tallyTaskEmitter) TaskSuccess(context.Context) {
	e.scope.Counter("task.success").Inc(1)
}

func (e *tallyTaskEmitter) TaskDone(_ context.Context, d time.Duration) {
	e.scope.Timer("task.timing").Record(d)
}

type tallySchedulerEmitter struct {
	scope tally.Scope
}

func (tallySchedulerEmitter) schedulerEmitter() {}

// SchedulerInit constructs a tally SchedulerEmitter.
func (e *tallyEmitter) SchedulerInit(info *SchedulerInfo) SchedulerEmitter {
	flow := info.FlowInfo.Name
	cacheKey := cacheKey{
		FlowName:   flow,
		FlowFile:   info.FlowInfo.File,
		FlowLine:   info.FlowInfo.Line,
		FlowColumn: info.FlowInfo.Column,
	}
	if v, ok := e.scheds.Load(cacheKey); ok {
		return v.(SchedulerEmitter)
	}
	scope := e.scope
	if flow != "" {
		scope = scope.Tagged(map[string]string{"flow": flow})
	}
	tse := &tallySchedulerEmitter{
		scope: scope,
	}
	v, _ := e.scheds.LoadOrStore(cacheKey, tse)
	return v.(SchedulerEmitter)
}

func (e *tallySchedulerEmitter) EmitScheduler(s SchedulerState) {
	e.scope.Gauge("scheduler.pending").Update(float64(s.Pending))
	e.scope.Gauge("scheduler.ready").Update(float64(s.Ready))
	e.scope.Gauge("scheduler.waiting").Update(float64(s.Waiting))
	e.scope.Gauge("scheduler.idle_workers").Update(float64(s.IdleWorkers))
	e.scope.Gauge("scheduler.concurrency").Update(float64(s.Concurrency))
}
