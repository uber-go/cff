package cff

// DirectiveType identifies the type of code generation directive
// for [Emitter] operations.
type DirectiveType int

const (
	// UnknownDirective is an invalid value for a DirectiveType.
	UnknownDirective DirectiveType = iota

	// FlowDirective marks a Flow.
	FlowDirective

	// ParallelDirective marks a Parallel.
	ParallelDirective
)

// String returns the directive string.
func (d DirectiveType) String() string {
	switch d {
	case FlowDirective:
		return "flow"
	case ParallelDirective:
		return "parallel"
	}
	return "unknown"
}
