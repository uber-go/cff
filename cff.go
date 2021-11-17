// Package cff along with the cff binary, provides a means of easily
// orchestrating a number of related functions with automated concurrent
// execution.
//
// Specify one or more flows in your code with the Flow function.
package cff

import "context"

const _noGenMsg = "code not generated; run cff"

// NOTE: All code generation directives must be added to this file. The list
// of directives is updated automatically based on the contents of this file.

// Option specifies parameters for a Flow or Parallel.
type Option interface {
	cffOption()
}

// Params specifies inputs for a Flow that do not come from tasks. These
// values are made available to the Flow as-is.
//
//  cff.Params(request)
//
// This is a code generation directive.
func Params(args ...interface{}) Option {
	panic(_noGenMsg)
}

// Results specifies one or more outputs for a Flow as pointers.
//
//  var result *Response
//  err := cff.Flow(ctx,
//    cff.Results(&result),
//    cff.Task(...),
//  )
//
// This is a code generation directive.
func Results(results ...interface{}) Option {
	panic(_noGenMsg)
}

// WithEmitter provides an optional observer for flow events. Emitters can
// track metrics, logs, or other observability data.
//
//  cff.Flow(ctx,
//    ...
//    cff.WithEmitter(cff.TallyEmitter(scope)),
//  )
//
// Provide this option multiple times to connect multiple emitters.
//
//  cff.Flow(ctx,
//    ...
//    cff.WithEmitter(cff.TallyEmitter(scope)),
//    cff.WithEmitter(cff.LogEmitter(logger)),
//  )
//
// This is a code generation directive.
func WithEmitter(Emitter) Option {
	panic(_noGenMsg)
}

// Task specifies a task for execution with a cff.Flow or cff.Parallel. A Task
// is any executable function or bound method available in the scope when
// cff.Flow or cff.Parallel is called.
//
// A Task's usage and constraints change between its usage in a cff.Flow or
// cff.Parallel.
//
// Within a cff.Flow,
//
//  cff.Flow(
//    ctx,
//    cff.Task(h.client.GetUser),
//    cff.Task(bindUser),
//    cff.Task(h.processRequest),
//  )
//
// Each Task has zero or more inputs, specified by the arguments of the
// function, and one or more results, specified by the return values of the
// function.
//
//  func(I1, I2, ...) (R1, R2, ...)
//
// Before this function is executed, all the tasks providing the inputs it
// depends on will have finished executing. Similarly, no task that depends
// on a result of this function will be executed until this function finishes
// executing.
//
// Tasks can request the context for the current execution scope by adding a
// context.Context as their first argument.
//
//  func(context.Context, I1, I2, ...) (R1, R2, ...)
//
// Additionally, tasks that may fail can do so by adding an error as their
// last return value.
//
//  func(I1, I2, ...) (R1, R2, ..., error)
//  func(context.Context, I1, I2, ...) (R1, R2, ..., error)
//
// Task behavior may further be customized with TaskOptions.
//
// Within a cff.Parallel, a Task is a function executed in parallel.
//
// Task can request the context for the current execution scope by optionally
// adding a context.Context as the only argument.
//
// Additionally, Tasks that may fail can do so by optionally adding an error
// as the only return value.
//
//  cff.Parallel(
//    ctx,
//    cff.Task(func() {...}),
//    cff.Task(func() error {...}),
//    cff.Task(func(ctx context.Context) {...}),
//    cff.Task(func(ctx context.Context) error {...}),
//  )
//
// Task functions under cff.Parallel cannot accept other arguments
// or return other values.
//
// This is a code generation directive.
func Task(fn interface{}, opts ...TaskOption) Option {
	panic(_noGenMsg)
}

// InstrumentFlow specifies that this Flow should be instrumented for
// observability. The provided name will be used in emitted metrics, logs, and
// spans, if any.
//
// This is a code generation directive.
func InstrumentFlow(name string) Option {
	panic(_noGenMsg)
}

// Concurrency specifies the maximum number of goroutines CFF2 should use to
// execute the tasks of this Flow.
//
// Defaults to max(GOMAXPROCS, 4).
//
// This option has effect only if the online_scheduling option is enabled.
func Concurrency(n int) Option {
	panic(_noGenMsg)
}

// ContinueOnError configures a Parallel to run through any errors returned
// by tasks over the course of its execution and continue executing remaining
// tasks when it would have otherwise stopped at the first error.
//
// err := cff.Parallel(ctx,
//     cff.Tasks(
//         func() error { ... },
//     ),
//     cff.ContinueOnError(true),
// )
//
// If one or more errors are encountered during when ContinueOnError is
// configured to true, the Parallel will return an error that accumulates the
// messages of all encountered errors after executing all remaining tasks.
//
// This is a code generation directive.
func ContinueOnError(bool) Option {
	panic(_noGenMsg)
}

// Flow specifies a single Flow for execution with CFF. The provided context
// is made available to all tasks in the Flow.
//
// A Flow MUST have at least one Task (specified with Task or Tasks), and at
// least one Results.
//
//  cff.Flow(ctx,
//    cff.Results(&result),
//    cff.Task(
//      ...
//    ),
//  )
//
// Tasks may be specified in any order. They will be connected based on their
// inputs and outputs. If any of the tasks fail, the entire Flow fails and the
// corresponding error is returned.
//
// This is a code generation directive.
func Flow(ctx context.Context, opts ...Option) error {
	panic(_noGenMsg)
}

// TaskOption customizes the execution behavior of a single Task.
type TaskOption interface {
	cffTaskOption()
}

// FallbackWith specifies that if the corresponding Task failed with an error
// or panic, we should recover from that failure and return the provided
// values instead.
//
// This function accepts the same number of values as returned by the Task
// with exactly the same types. This DOES NOT include the error type returned
// by the Task.
//
//   cff.Task(client.ListUsers, cff.FallbackWith(cachedUserList))
//
// This is a code generation directive.
func FallbackWith(results ...interface{}) TaskOption {
	panic(_noGenMsg)
}

// Predicate specifies that the corresponding Task should be executed only if
// the provided function returns true.
//
// This accepts a function which has the signature,
//
//   func(I1, I2, ...) bool
//
// Where the arguments of the functions are inputs similar to a Task. Any type
// added here becomes a dependency of the Task if it is not already. This
// means that the Task will not be executed until these types are available.
//
// When specified, the corresponding Task will be executed only if this
// function returns true.
//
// If the function evaluates to false, the annotated function is skipped and
// tasks that depend on the type provided by that function get the zero value
// for that type.
//
// This is a code generation directive.
func Predicate(fn interface{}) TaskOption {
	panic(_noGenMsg)
}

// Instrument specifies that this Task should be instrumented for
// observability. The provided name will be used in emitted metrics, logs, and
// spans, if any.
//
// This is a code generation directive.
func Instrument(name string) TaskOption {
	panic(_noGenMsg)
}

// Invoke specifies that task must always be executed, even if none of other
// tasks consume its output.
//
// Only tasks marked with Invoke are allowed to have zero non-error or
// single error returns.
//
// This is a code generation directive.
func Invoke(enable bool) TaskOption {
	panic(_noGenMsg)
}

// Parallel specifies a Parallel operation for execution with CFF. The provided
// context is made available to all tasks in the Parallel.
//
// A Parallel MUST have at least one Tasks function.
//
//  cff.Parallel(ctx,
//    cff.Concurrency(4),
//    cff.Tasks(
//      func(ctx context.Context) error {
//         ...
//      },
//      ...
//    ),
//  )
//
// Tasks will run independently with bounded parallelism with all other
// Parallel declared tasks. If any of the tasks fail, Parallel stops
// processsing outstanding tasks and an error is returned.
//
// This is a code generation directive.
func Parallel(ctx context.Context, opts ...Option) error {
	panic(_noGenMsg)
}

// InstrumentParallel specifies that this Parallel should be instrumented for
// observability. The provided name will be used in emitted metrics, logs, and
// spans, if any.
//
// This is a code generation directive.
func InstrumentParallel(name string) Option {
	panic(_noGenMsg)
}

// Tasks specifies functions for execution with Parallel. Tasks are any
// executable function or bound method available in the scope when cff.Parallel
// is called.
// Tasks can request the context for the current execution scope by optionally
// adding a context.Context as the only argument.
// Additionally, Tasks that may fail can do so by optionally adding an error
// as the only return value.
//
//  func(context.Context) error
//
// Tasks functions do not accept other arguments or return values.
//
// This is a code generation directive.
func Tasks(fn ...interface{}) Option {
	panic(_noGenMsg)
}

// Slice executes a parallel operation on the elements of the provided slice.
//
// The fn parameter is function that is invoked on each element of the slice
// parameter.
//
// The fn parameter's first non-context argument is a slice index of type int
// followed by a value of same type as the slice parameter's elements.
//
//  cff.Parallel(
//  	ctx,
//  	cff.Concurrency(...),
//  	cff.Slice(func(idx int, elem someType) { ... }, []someType{...})
//  )
//
// Optionally, a context.Context can be provided as a first argument to the
// execution function.
//
//  func(ctx context.Context, idx int, item someType) { ... }
//
// Optionally, an error can be returned as the execution function's sole
// return value.
//
//  func(idx int, item someType) error { ... }
//
// The second argument to Slice is a slice on which the execution function
// is invoked.
//
// cff.Slice is only an option for cff.Parallel.
//
// This is a code generation directive.
func Slice(fn interface{}, slice interface{}) Option {
	panic(_noGenMsg)
}
