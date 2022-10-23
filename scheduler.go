package cff

import (
	"go.uber.org/atomic"
	"go.uber.org/cff/scheduler"
)

// We re-export things here so that users of cff don't have to add other
// packages as dependencies to their BUILD.bazel.

// Job is a job prepared to be enqueued to the cff scheduler.
type Job = scheduler.Job

// AtomicBool is a type-safe means of reading and writing boolean values.
type AtomicBool = atomic.Bool

// ScheduledJob is a job that has been scheduled for execution with the cff
// scheduler.
type ScheduledJob = scheduler.ScheduledJob

// SchedulerParams configures the cff scheduler.
type SchedulerParams struct {
	// Concurrency specifies the number of concurrent workers
	// used by the scheduler to run jobs.
	//
	// See cff.Concurrency for more details.
	Concurrency int
	// Emitter provides an emitter for the scheduler.
	//
	// See cff.SchedulerEmitter for more details.
	Emitter SchedulerEmitter
	// ContinueOnError when true directs the scheduler to continue running
	// through job errors.
	//
	// See cff.ContinueOnError for more details.
	ContinueOnError bool
}

// NewScheduler returns a new Scheduler with a maximum of n workers. Enqueue
// jobs into the returned scheduler in topological order using the Enqueue
// method, and wait for results with Wait.
//
//	sched := cff.NewScheduler(..)
//	j1 := sched.Enqueue(cff.Job{...}
//	j2 := sched.Enqueue(cff.Job{..., Dependencies: []*cff.ScheduledJob{j1}}
//	// ...
//	err := sched.Wait()
func NewScheduler(p SchedulerParams) *scheduler.Scheduler {
	cfg := scheduler.Config{
		Concurrency:     p.Concurrency,
		Emitter:         adaptSchedulerEmitter(p.Emitter),
		ContinueOnError: p.ContinueOnError,
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
