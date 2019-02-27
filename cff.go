package cff

import (
	"context"

	"github.com/uber-go/tally"
)

// FlowOption TODO
type FlowOption interface {
	cffOption()
}

// Provide TODO
func Provide(args ...interface{}) FlowOption {
	panic("code not generated; run cff")
}

// Result TODO
func Result(results ...interface{}) FlowOption {
	panic("code not generated; run cff")
}

// Scope TODO
func Scope(scope tally.Scope) FlowOption {
	panic("code not generated; run cff")
}

// Tasks TODO
func Tasks(tasks ...interface{}) FlowOption {
	panic("code not generated; run cff")
}

// InstrumentFlow TODO
func InstrumentFlow(name string) FlowOption {
	panic("code not generated; run cff")
}

// Flow TODO
func Flow(ctx context.Context, opts ...FlowOption) error {
	panic("code not generated; run cff")
}

// Task TODO
func Task(fn interface{}, opts ...TaskOption) FlowOption {
	panic("code not generated; run cff")
}

// TaskOption TODO
type TaskOption interface {
	cffTaskOption()
}

// RecoverWith TODO
func RecoverWith(results ...interface{}) TaskOption {
	panic("code not generated; run cff")
}

// Predicate TODO
func Predicate(fn interface{}) TaskOption {
	panic("code not generated; run cff")
}

// Instrument TODO
func Instrument(name string) TaskOption {
	panic("code not generated; run cff")
}
