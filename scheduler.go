package cff

import (
	"go.uber.org/atomic"
	"go.uber.org/cff/scheduler"
)

// We re-export things here so that users of cff don't have to add other
// packages as dependencies to their BUILD.bazel.

// Job is a job prepared to be enqueued to the cff scheduler.
//
// This is intended to be used by cff's generated code.
// Do not use directly.
// This can change without warning.
type Job = scheduler.Job

// AtomicBool is a type-safe means of reading and writing boolean values.
//
// This is intended to be used by cff's generated code.
// Do not use directly.
// This can change without warning.
type AtomicBool = atomic.Bool

// TODO(abg): For Go 1.19 or newer, we can use sync/atomic.Bool
// which drops one more dependency for users.

// ScheduledJob is a job that has been scheduled for execution with the cff
// scheduler.
//
// This is intended to be used by cff's generated code.
// Do not use directly.
// This can change without warning.
type ScheduledJob = scheduler.ScheduledJob

// SchedulerParams configures the cff scheduler.
//
// This is intended to be used by cff's generated code.
// Do not use directly.
// This can change without warning.
type SchedulerParams struct {
	// Concurrency specifies the number of concurrent workers
	// used by the scheduler to run jobs.
	Concurrency int
	// Emitter provides an emitter for the scheduler.
	Emitter SchedulerEmitter
	// ContinueOnError when true directs the scheduler to continue running
	// through job errors.
	ContinueOnError bool
}

// NewScheduler starts up a cff scheduler for use by Flow or Parallel.
//
//	sched := cff.NewScheduler(..)
//	j1 := sched.Enqueue(cff.Job{...}
//	j2 := sched.Enqueue(cff.Job{..., Dependencies: []*cff.ScheduledJob{j1}}
//	// ...
//	err := sched.Wait()
//
// This is intended to be used by cff's generated code.
// Do not use directly.
// This can change without warning.
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
