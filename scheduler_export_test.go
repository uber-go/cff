package cff

import "go.uber.org/cff/scheduler"

func AdaptSchedulerEmitter(e SchedulerEmitter) scheduler.Emitter {
	return adaptSchedulerEmitter(e)
}
