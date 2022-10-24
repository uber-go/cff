package cff

import "go.uber.org/cff/scheduler"

// AdaptSchedulerEmitter adapts a cff.SchedulerEmitter to a scheduler.Emitter.
func AdaptSchedulerEmitter(e SchedulerEmitter) scheduler.Emitter {
	return adaptSchedulerEmitter(e)
}
