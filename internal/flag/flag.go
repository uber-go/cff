package flag

import "flag"

// This file re-exports functions and types from the standard library flag
// package.

// ErrHelp reports that the user requested help with --help/-h.
var ErrHelp = flag.ErrHelp

// Set is a collection of flags.
type Set = flag.FlagSet

// Getter is the interface implemented by custom flag values.
type Getter = flag.Getter

// NewSet builds a new flag set.
func NewSet(name string) *Set {
	return flag.NewFlagSet(name, flag.ContinueOnError)
}
