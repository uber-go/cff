// Package cff along with the cff binary, provides a means of easily
// orchestrating a number of related functions with automated concurrent
// execution.
//
// Specify one or more flows in your code with the Flow function, tag your Go
// file with the "cff" build tag, and run the cff tool.
package cff

import (
	"context"

	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

// FlowOption specifies parameters for a Flow.
type FlowOption interface {
	cffOption()
}

// Params specifies inputs for a Flow that do not come from tasks. These
// values are made available to the Flow as-is.
//
//  cff.Params(request)
func Params(args ...interface{}) FlowOption {
	panic("code not generated; run cff")
}

// Results specifies one or more outputs for a Flow as pointers.
//
//  var result *Response
//  err := cff.Flow(ctx,
//    cff.Results(&result),
//    cff.Task(...),
//  )
func Results(results ...interface{}) FlowOption {
	panic("code not generated; run cff")
}

// Metrics provides the Tally scope to which metrics will be logged for Tasks
// and Flows that have been instrumented with Instrument or InstrumentFlow.
func Metrics(scope tally.Scope) FlowOption {
	panic("code not generated; run cff")
}

// Logger provides the logger to which messages will be logged for Tasks and
// Flows that have been instrumented with Instrument or InstrumentFlow.
func Logger(logger *zap.Logger) FlowOption {
	panic("code not generated; run cff")
}

// WithEmitter provides an optional observer for task metric events.
// If provided, we will use the provided Emitter rather than the default
// metrics implementation.
func WithEmitter(Emitter) FlowOption {
	panic("code not generated; run cff")
}

// WithLogFields specifies additional log fields that should be added to emitted
// logs.
func WithLogFields(...zap.Field) FlowOption {
	panic("code not generated; run cff")
}

// Task specifies a task for execution with a flow. A Task is any executable
// function or bound method available in the scope when cff.Flow is called.
//
//  cff.Flow(ctx,
//    ...
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
func Task(fn interface{}, opts ...TaskOption) FlowOption {
	panic("code not generated; run cff")
}

// InstrumentFlow specifies that this Flow should be instrumented for
// observability. The provided name will be used in emitted metrics, logs, and
// spans, if any.
func InstrumentFlow(name string) FlowOption {
	panic("code not generated; run cff")
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
func Flow(ctx context.Context, opts ...FlowOption) error {
	panic("code not generated; run cff")
}

// TaskOption customizes the execution behavior of a single Task.
type TaskOption interface {
	cffTaskOption()
}

// FallbackWith specifies that if the corresponding Task failed with an error,
// we should recover from that error and return the provided values instead.
//
// This function accepts the same number of values as returned by the Task
// with exactly the same types. This DOES NOT include the error type returned
// by the Task.
//
//   cff.Task(client.ListUsers, cff.FallbackWith(cachedUserList))
//
// Note that this function DOES NOT recover from Task panics. If panics should
// be handled, it is the caller or the Task implementation's responsibility.
func FallbackWith(results ...interface{}) TaskOption {
	panic("code not generated; run cff")
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
func Predicate(fn interface{}) TaskOption {
	panic("code not generated; run cff")
}

// Instrument specifies that this Task should be instrumented for
// observability. The provided name will be used in emitted metrics, logs, and
// spans, if any.
func Instrument(name string) TaskOption {
	panic("code not generated; run cff")
}

// Invoke specifies that task must always be executed, even if none of other
// tasks consume its output.
//
// Only tasks marked with Invoke are allowed to have zero non-error or
// single error returns.
func Invoke(enable bool) TaskOption {
	panic("code not generated; run cff")
}
