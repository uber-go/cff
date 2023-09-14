// Package emittertest provides testing utilities for cff emitters.
package emittertest

//go:generate mockgen -destination mock_emitter.go -package emittertest go.uber.org/cff Emitter,TaskEmitter,FlowEmitter,ParallelEmitter,SchedulerEmitter
