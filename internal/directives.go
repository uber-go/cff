package internal

// List of functions in the CFF package that are code generation directives.
var _codegenDirectives = []string{
	"FallbackWith",
	"Flow",
	"Instrument",
	"InstrumentFlow",
	"Invoke",
	"Logger",
	"Params",
	"Predicate",
	"Results",
	"Task",
	"WithEmitter",
}

var _directiveIndex map[string]struct{}

func init() {
	_directiveIndex = make(map[string]struct{}, len(_codegenDirectives))
	for _, f := range _codegenDirectives {
		_directiveIndex[f] = struct{}{}
	}
}

// IsCodegenDirective reports whether the function with the given name in the
// CFF package is a code generation directive.
func IsCodegenDirective(name string) bool {
	_, ok := _directiveIndex[name]
	return ok
}
