package cff

import (
	"go.uber.org/cff/scheduler"
	"go.uber.org/atomic"
)

// We re-export things here so that users of CFF don't have to add other
// packages as dependencies to their BUILD.bazel.

// Job is a job prepared to be enqueued to the CFF scheduler.
type Job = scheduler.Job

// AtomicBool is a type-safe means of reading and writing boolean values.
type AtomicBool = atomic.Bool

// ScheduledJob is a job that has been scheduled for execution with the CFF
// scheduler.
type ScheduledJob = scheduler.ScheduledJob

// BeginFlow begins execution of a flow with a maximum of n workers. Enqueue
// jobs into the returned scheduler in topological order using the Enqueue
// method, and wait for results with Wait.
//
//  sched := cff.BeginFlow(..)
//  j1 := sched.Enqueue(cff.Job{...}
//  j2 := sched.Enqueue(cff.Job{..., Dependencies: []*cff.ScheduledJob{j1}}
//  // ...
//  err := sched.Wait()
func BeginFlow(n int, e SchedulerEmitter) *scheduler.Scheduler {
	cfg := scheduler.Config{
		Concurrency: n,
		Emitter:     adaptSchedulerEmitter(e),
	}
	return cfg.New()
}

// schedulerAdapter adapts a SchedulerEmitter into a scheduler.Emitter.
type schedulerAdapter struct {
	emitter SchedulerEmitter
}

func adaptSchedulerEmitter(e SchedulerEmitter) scheduler.Emitter {
	if _, isNop := e.(*nopEmitter); e == nil || isNop {
		// Avoid the cost of a live ticker in the scheduler if we're using
		// no emitter or a no-op emitter.
		return nil
	}
	return schedulerAdapter{
		emitter: e,
	}
}

func (s schedulerAdapter) Emit(state scheduler.State) {
	s.emitter.EmitScheduler(SchedulerState(state))
}
