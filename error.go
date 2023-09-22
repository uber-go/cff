// Copyright (c) 2022 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cff

import (
	"fmt"
)

// PanicError is an error that is thrown when a task panics. It contains the value
// that is recovered from the panic and the stacktrace of where the panic happened.
// For example, the following code checks if an error from [Flow] is due to a panic:
//
//	var r string
//	err := cff.Flow(
//		context.Background(),
//		cff.Results(&r),
//		cff.Task(
//			func() string {
//				panic("panic")
//			},
//		),
//	)
//	var panicError *cff.PanicError
//	if errors.As(err, &panicError) {
//		// err is from a panic
//		fmt.Printf("recovered: %s\n", panicError.Value)
//	} else {
//		// err is not from a panic
//	}
type PanicError struct {
	// Value contains the value recovered from the panic that caused this error.
	Value any

	// Stacktrace contains string of what call stack looks like when the panic happened.
	// This is populated by calling runtime/debug.Stack() when a non-nil value is
	// recovered from a cff-scheduled job.
	Stacktrace []byte
}

var _ error = (*PanicError)(nil)

func (pe *PanicError) Error() string {
	return fmt.Sprintf("panic: %v\n%s", pe.Value, pe.Stacktrace)
}
