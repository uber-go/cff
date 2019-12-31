package noinfra

import (
	"go/ast"
	"strconv"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	_cff2Pkg = "go.uber.org/cff"
	_cff2ERD = "https://t.uber.com/cff2"
)

// Analyzer is a lint rule that checks for disallowed usage of CFF2.
var Analyzer = &analysis.Analyzer{
	Name:     "cff2_noinfra",
	Run:      run,
	Doc:      "Checks for use of CFF2 in Infra projects.",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspct := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.ImportSpec)(nil),
	}

	inspct.Preorder(nodeFilter, func(n ast.Node) {
		// If there is an error, we return "" which won't match.
		path := importPath(n.(*ast.ImportSpec))
		if !isCFF2Import(path) {
			return
		}

		pass.Reportf(n.Pos(),
			"found usage of %s in a disallowed directory. Please see %s for why this isn't allowed."+
				" If you would like an exemption, please email cff-group@uber.com.",
			path,
			_cff2ERD,
		)
	})

	return nil, nil
}

// importPath returns unquoted import path.
func importPath(n *ast.ImportSpec) string {
	path, err := strconv.Unquote(n.Path.Value)
	// Don't want to fail here as this linter won't affect code flow, and is more of
	// a guidance which can also be spotted through a CR.
	if err != nil {
		return ""
	}
	return path
}

func isCFF2Import(path string) bool {
	return _cff2Pkg == path
}
