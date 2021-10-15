package cff

// DirectiveType identifies the type of code generation directive.
type DirectiveType int

const (
	// UnknownDirective is an unknown directive.
	UnknownDirective DirectiveType = iota
	// FlowDirective is a cff.Flow directive.
	FlowDirective
	// ParallelDirective is a cff.Parallel directive.
	ParallelDirective
)

// String returns the directive string.
func (d DirectiveType) String() string {
	if d == FlowDirective {
		return "flow"
	} else if d == ParallelDirective {
		return "parallel"
	}
	return "unknown"
}
