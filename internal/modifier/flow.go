package modifier

import (
	"fmt"
	"go/ast"
	"go/token"
)

type flowModifier struct {
	Position token.Position

	node ast.Node
}

var _ Modifier = (*flowModifier)(nil)

// NewFlowModifier returns a Modifier that corresponds to a
// cff.Flow call.
func NewFlowModifier(fset *token.FileSet, n ast.Node) Modifier {
	return &flowModifier{
		Position: fset.Position(n.Pos()),
		node:     n,
	}
}

func (fm *flowModifier) FuncExpr() string {
	return fmt.Sprintf("_cffFlow%d_%d", fm.Position.Line, fm.Position.Column)
}

func (fm *flowModifier) RetExpr() string {
	// TODO: this probably needs a type generation as well.
	return "error"
}

func (fm *flowModifier) GenImpl(p GenParams) error {
	// TODO
	return nil
}

func (fm *flowModifier) Node() ast.Node {
	return fm.node
}
