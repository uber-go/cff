package internal

// List of functions in the cff package that are code generation directives.
var _codegenDirectives = map[string]struct{}{
	"Params":             {},
	"Results":            {},
	"WithEmitter":        {},
	"Task":               {},
	"InstrumentFlow":     {},
	"Concurrency":        {},
	"ContinueOnError":    {},
	"Flow":               {},
	"FallbackWith":       {},
	"Predicate":          {},
	"Instrument":         {},
	"Invoke":             {},
	"Parallel":           {},
	"InstrumentParallel": {},
	"Tasks":              {},
	"Slice":              {},
	"SliceEnd":           {},
	"Map":                {},
	"MapEnd":             {},
}

// IsCodegenDirective reports whether the function with the given name in the
// cff package is a code generation directive.
func IsCodegenDirective(name string) bool {
	_, ok := _codegenDirectives[name]
	return ok
}
