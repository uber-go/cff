package cff

import "context"
import "github.com/uber-go/tally"

// Option TODO
type Option interface {
	cffOption()
}

// Provide TODO
func Provide(args ...interface{}) Option {
	panic("code not generated; run cff")
}

// Result TODO
func Result(results ...interface{}) Option {
	panic("code not generated; run cff")
}

// Scope TODO
func Scope(scope tally.Scope) Option {
	panic("code not generated; run cff")
}

// Tasks TODO
func Tasks(tasks ...interface{}) Option {
	panic("code not generated; run cff")
}

// InstrumentFlow TODO
func InstrumentFlow(name string) Option {
	panic("code not generated; run cff")
}

// Flow TODO
func Flow(ctx context.Context, opts ...Option) error {
	panic("code not generated; run cff")
}

// Task TODO
func Task(fn interface{}, opts ...TaskOption) Option {
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
