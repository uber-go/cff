package internal

import "go/ast"

type visitorFunc func(ast.Node) (recurse bool)

func (f visitorFunc) Visit(n ast.Node) ast.Visitor {
	if f(n) {
		return f
	}
	return nil
}

// Typed version of ast.Walk for convenience.
func astWalk(n ast.Node, f visitorFunc) {
	ast.Walk(f, n)
}
