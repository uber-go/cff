// Package cff along with the cff CLI, provides a means of easily writing
// common concurrent code patterns in a type-safe manner.
//
// # Code generation directives
//
// Some APIs in this package are marked as "code generation directives."
// If you use a code generation directive in a file,
// that file must have the 'cff' build constraint on top:
//
//	//go:build cff
//
// Following that, you must run the following command before you use
// 'go build' or 'go test' with that file.
//
//	cff ./...
package cff

import "context"

const _noGenMsg = `If you're seeing this error, you probably built code that uses cff without processing it with cff.
Ensure that .go files that use cff have '//go:build cff' on top and run 'cff ./...'`

// NOTE: All code generation directives must be added to this file. The list
// of directives is updated automatically based on the contents of this file.

// Option is an argument for a [Flow] or [Parallel].
//
// See individual option documentation for details.
type Option interface {
	cffOption()
}

// Params specifies inputs for a Flow that do not have any dependencies.
// These values are made available to the Flow as-is.
//
// For example:
//
//	var req *GetUserRequest = // ...
//	cff.Flow(
//		cff.Params(req),
//		// ...
//	)
//
// This is a code generation directive.
func Params(args ...interface{}) Option {
	panic(_noGenMsg)
}

// Results specifies one or more outputs for a Flow.
// Arguments to Results must be pointers to variables
// that will hold the result values.
//
// For example:
//
//	var result *GetUserResponse
//	err := cff.Flow(ctx,
//		cff.Results(&result),
//		// ...
//	)
//
// This is a code generation directive.
func Results(results ...interface{}) Option {
	panic(_noGenMsg)
}

// WithEmitter provides an optional observer for [Flow] or [Parallel] events.
// Emitters can track metrics, logs, or other observability data.
//
//	cff.Flow(ctx,
//		// ...
//		cff.WithEmitter(em),
//	)
//
// Provide this option multiple times to use multiple emitters.
//
// WARNING: Do not use this API.
// We intend to replace it in an upcoming release.
//
// This is a code generation directive.
func WithEmitter(Emitter) Option {
	panic(_noGenMsg)
}

// Task specifies a task for execution with a [Flow] or [Parallel].
// A task can be a reference to:
//
//   - a top-level function; or
//   - a bound method; or
//   - an anonymous function
//
// For example:
//
//	// Given,
//	//   var client *Client
//	//   func (*Client) GetUser(...) (...)
//	// The following is a bound method reference.
//	cff.Task(client.GetUser)
//
//	// Given,
//	//   func bindUser(...) (...)
//	// The following is a top-level function reference.
//	cff.Task(bindUser),
//
//	// The following is an anonymous function reference.
//	cff.Task(func(...) (...,, error) {
//		// ...
//	})
//
// A Task's usage and constraints change based on whether you're using it
// inside a Flow or a Parallel.
// See the documentation for [Flow] or [Parallel] for more details.
//
// This is a code generation directive.
func Task(fn interface{}, opts ...TaskOption) Option {
	panic(_noGenMsg)
}

// InstrumentFlow specifies that this Flow should be instrumented for
// observability.
// The provided name will be passed to the [Emitter] you passed into
// WithEmitter.
//
// This is a code generation directive.
func InstrumentFlow(name string) Option {
	panic(_noGenMsg)
}

// Concurrency specifies the maximum number of goroutines cff should use to
// execute tasks of this Flow or Parallel.
//
// The default value for this is,
//
//	max(GOMAXPROCS, 4)
//
// That is, by default cff will use [runtime.GOMAXPROCS] goroutines,
// with a minimum of 4.
//
// This is a code generation directive.
func Concurrency(n int) Option {
	panic(_noGenMsg)
}

// ContinueOnError configures a [Parallel] to keep running all other tasks
// despite errors returned by tasks over the course of its execution.
// By default, Parallel will stop execution at the first error it encounters.
//
//	err = cff.Parallel(ctx,
//		cff.Task(task1),
//		cff.Task(task2),
//		// ...
//		cff.ContinueOnError(true),
//	)
//
// If one or more tasks return errors with ContinueOnError(true),
// Parallel will still run all the other tasks,
// and accumulate and combine the errors together into a single error object.
// You can access the full list of errors with [go.uber.org/multierr.Errors].
//
// ContinueOnError(true) is incompatible with [Flow], [SliceEnd] and [MapEnd].
//
// This is a code generation directive.
func ContinueOnError(bool) Option {
	panic(_noGenMsg)
}

// Flow specifies a single Flow for execution with cff.
// A child of the provided context is made available to all tasks in the Flow
// if they request it.
//
// A Flow MUST have at least one task (specified with [Task] or [Tasks]),
// and at least one result (specified with [Results]).
//
//	var result *Result
//	cff.Flow(ctx,
//		cff.Results(&result),
//		cff.Task(
//			// ...
//		),
//	)
//
// Tasks may be specified to a Flow in any order.
// They will be connected based on their inputs and outputs.
// If a task fails with an error,
// the entire Flow terminates and the error is returned.
//
// # Flow tasks
//
// Within a cff.Flow, each task has:
//
//   - zero or more inputs, specified by its parameters
//   - *one* or more outputs, specified by its return values
//   - optionally, a context.Context as the first parameter
//   - optionally, an error as the last return value
//
// This is roughly expressed as:
//
//	func([context.Context], I1, I2, ...) (R1, R2, ..., [error])
//
// The types of the inputs specify the dependencies of this task.
// cff will run other tasks that provide these dependencies
// and feed their results back into this task to run it.
// Similarly, it will feed the results of this task into other tasks that
// depend on them.
//
// Tasks may use the optional context argument to cancel operations early in
// case of failures:
// the context is valid only as long as the flow is running.
// If the flow terminates early because of a failure, the context is
// invalidated.
//
//	func(context.Context, I1, I2, ...) (R1, R2, ...)
//
// Fallible tasks may declare an error as their last return value.
// If a task fails, the flow is terminated and all ongoing tasks are canceled.
//
//	func(I1, I2, ...) (R1, R2, ..., error)
//
// Task behaviors may further be customized with [TaskOption].
//
// This is a code generation directive.
func Flow(ctx context.Context, opts ...Option) error {
	panic(_noGenMsg)
}

// TaskOption customizes the behavior of a single Task.
type TaskOption interface {
	cffTaskOption()
}

// FallbackWith specifies that
// if the corresponding task fails with an error or panics,
// we should recover from that failure and return the provided values instead.
//
// This function accepts the same number of values as returned by the task
// with exactly the same types -- not including the error return value (if
// any).
//
// For example:
//
//	// Given,
//	//   func (*Client) ListUsers(context.Context) ([]*User, error)
//	// And,
//	//   var cachedUserList []*User = ...
//	cff.Task(client.ListUsers, cff.FallbackWith(cachedUserList))
//
// If client.ListUsers returns an error or panics,
// cff will return cachedUserList instead.
//
// This is a code generation directive.
func FallbackWith(results ...interface{}) TaskOption {
	panic(_noGenMsg)
}

// Predicate specifies a function that determines if the corresponding task
// should run.
//
// The predicate function has the following signature:
//
//	func(I1, I2, ...) bool
//
// Where the arguments I1, I2, ... are inputs similar to a task.
// Arguments added to the predicate become a dependency of the task,
// so the predicate or the task will not run until that value is available.
//
// When specified, the corresponding task will be executed only if this
// function returns true.
// If the function evaluates to false, the cff will skip execution of this
// task.
// If any other tasks depend on this task,
// cff will give them zero values of the outputs of this task.
//
// For example:
//
//	cff.Task(
//		authorizeUser,
//		cff.Predicate(func(cfg *Config) bool {
//			return cfg.Prorudction == true
//		}),
//	)
//
// This is a code generation directive.
func Predicate(fn interface{}) TaskOption {
	panic(_noGenMsg)
}

// Instrument specifies that this Task should be instrumented for
// observability.
// The provided name will be passed to the [Emitter] you passed into
// WithEmitter.
//
// This is a code generation directive.
func Instrument(name string) TaskOption {
	panic(_noGenMsg)
}

// Invoke specifies that a flow task must be executed
// even if none of other tasks consume its output.
//
// By default, flow tasks have the following restrictions:
//
//   - must have a non-error return value (outputs)
//   - the output must be consumed by another task or flow result (via
//     [Results])
//
// A task tagged with Invoke(true) loses these restriction.
// It may have zero outputs, or if it has outputs,
// other tasks or flow results don't have to consume them.
//
//	cff.Task(func(ctx context.Context, req *Request) {
//		res, err := shadowClient.Send(req)
//		log.Info("shadowed request", "response", res, "error", err)
//	}, cff.Invoke(true))
//
// This is a code generation directive.
func Invoke(enable bool) TaskOption {
	panic(_noGenMsg)
}

// Parallel specifies a parallel operation for execution with cff.
//
// A Parallel must have at least one [Task], [Tasks], [Map], or [Slice].
//
//	cff.Parallel(ctx,
//		cff.Task(/* ... */)
//		cff.Slice(/* ... */)
//		cff.Map(/* ... */)
//	)
//
// Tasks inside a Parallel are all independent.
// They run concurrently with bounded parallelism.
//
// If any of the tasks fail with an error or panic,
// Parallel terminates the entire operation.
// You can change this with [ContinueOnError].
// With ContinueOnError(true), Parallel will run through all provided tasks
// despite errors,
// and return an aggregated error object representing all encountered failures.
//
// A child of the provided context is made available to all tasks in the
// parallel if they request it.
// If the context is cancelled or otherwise errors,
// Parallel does not run further tasks.
// This behaviour cannot be changed.
//
// # Parallel tasks
//
// Within a cff.Parallel, each task has:
//
//   - optionally, a context.Context as the first parameter
//   - optionally, an error as the last return value
//
// Note that tasks inside Parallel cannot have dependencies.
// Use [Flow] for that.
//
// This is roughly expressed as:
//
//	func([context.Context]) ([error])
//
// Tasks may use the context argument to cancel operations early in case of
// failures.
// Fallible tasks may return a non-nil error to signal failure.
//
// Task behaviors may further be customized with [TaskOption].
//
// This is a code generation directive.
func Parallel(ctx context.Context, opts ...Option) error {
	panic(_noGenMsg)
}

// InstrumentParallel specifies that this Parallel should be instrumented for
// observability.
// The provided name will be passed to the [Emitter] you passed into
// WithEmitter.
//
// This is a code generation directive.
func InstrumentParallel(name string) Option {
	panic(_noGenMsg)
}

// Tasks specifies multiple functions for execution with [Parallel].
// As with [Task], each argument to Tasks is a reference to:
//
//   - a top-level function; or
//   - a bound method; or
//   - an anonymous function
//
// They may all match the signature specified for parallel tasks (see
// [Parallel]).
//
// Tasks cannot be used with Flow. Use [Task] for that.
//
// This is a code generation directive.
func Tasks(fn ...interface{}) Option {
	panic(_noGenMsg)
}

// Slice runs fn in parallel on elements of the provided slice
// with a bounded number of goroutines.
//
//	cff.Parallel(ctx,
//		cff.Slice(
//			func(el someType) { ... },
//			[]someType{...},
//		),
//	)
//
// For a slice []T, fn has the following signature:
//
//	func([ctx context.Context,] [idx int,] value T) ([error])
//
// That is, it has the following parameters in-order:
//
//   - an optional context.Context
//   - an optional integer holding the index of the element in the slice
//   - a value in the slice
//
// And if the operation is fallible, it may have an error return value.
// A non-nil error returned by the function halts the entire Parallel
// operation.
// Use [ContinueOnError] to change this.
//
// Slice may only be used with [Parallel].
//
// This is a code generation directive.
func Slice(fn interface{}, slice interface{}, opts ...SliceOption) Option {
	panic(_noGenMsg)
}

// SliceOption customizes the execution behavior of [Slice].
type SliceOption interface {
	cffSliceOption()
}

// SliceEnd specifies a function for execution when a [Slice] operation
// finishes.
// This function will run after all items in the slice have been processed.
//
// As with parallel tasks, the function passed to SliceEnd may have:
//
//   - an optional context.Context parameter
//   - an optional error return value
//
// Therefore, these are all valid:
//
//	cff.SliceEnd(func() {...})
//	cff.SliceEnd(func() error {...})
//	cff.SliceEnd(func(ctx context.Context) {...})
//	cff.SliceEnd(func(ctx context.Context) error {...})
//
// SliceEnd cannot be used with [ContinueOnError].
//
// This is a code generation directive.
func SliceEnd(fn interface{}) SliceOption {
	panic(_noGenMsg)
}

// Map runs fn in parallel on elements of the provided map
// with a bounded number of goroutines.
//
//	cff.Parallel(ctx,
//		cff.Map(
//			func(k string, v *User) { /* ... */ },
//			map[string]*User{ /* ... */ },
//		),
//	)
//
// For a slice map[K]V, fn has the following signature:
//
//	func([ctx context.Context,] k K, v V) ([error])
//
// That is, it has the following parameters in-order:
//
//   - an optional context.Context
//   - a key in the map
//   - the value of that key in the map
//
// And if the operation is fallible, it may have an error return value.
// A non-nil error returned by the function halts the entire Parallel
// operation.
// Use [ContinueOnError] to change this.
//
// Map may only be used with [Parallel].
//
// This is a code generation directive.
func Map(fn interface{}, m interface{}, opts ...MapOption) Option {
	panic(_noGenMsg)
}

// MapOption customizes the execution behavior of [Map].
type MapOption interface {
	cffMapOption()
}

// MapEnd specifies a function for execution when a [Map] operation finishes.
// This function will run after all items in the map have been processed.
//
// As with parallel tasks, the function passed to MapEnd may have:
//
//   - an optional context.Context parameter
//   - an optional error return value
//
// Therefore, these are all valid:
//
//	cff.MapEnd(func() {...})
//	cff.MapEnd(func() error {...})
//	cff.MapEnd(func(ctx context.Context) {...})
//	cff.MapEnd(func(ctx context.Context) error {...})
//
// MapEnd cannot be used with [ContinueOnError].
//
// This is a code generation directive.
func MapEnd(fn interface{}) MapOption {
	panic(_noGenMsg)
}
