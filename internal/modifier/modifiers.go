package modifier

import (
	"go/ast"
	"io"
)

// Modifier changes the existing code by doing two things.
// 1. It generates a function that corresponds to what the cff "API"s do by inspecting the
//    arguments.
// 2. It inline replaces the cff "API" call with calls to corrresponding generated function.
// Each call to cff "API" is translates to a modifier.
type Modifier interface {
	Node() ast.Node            // The ast Node that produced this modifier.
	FuncExpr() string          // The name of the modifier-generated function.
	RetExpr() string           // The return signature expression of a modifier
	GenImpl(p GenParams) error // Generates the function body of the modifier-generated function.
}

// GenParams is the parameter for Modifiers to generate the body
// of their corresponding methods.
type GenParams struct {
	Writer io.Writer
}
