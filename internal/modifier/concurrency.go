package modifier

import (
	"embed"
	"fmt"
	"go/ast"
	"go/token"
	"text/template"
)

var tmplFS embed.FS

const (
	_modifierTmplDir = "templates/modifiers/*"
	_concurrencyTmpl = "concurrency.go.tmpl"
)

type concurrencyModifier struct {
	Position token.Position

	node ast.Node
}

var _ Modifier = (*concurrencyModifier)(nil)

// NewConcurrencyModifier returns a Modifier that corresponds to
// a cff.Concurrency call.
func NewConcurrencyModifier(fset *token.FileSet, n ast.Node) Modifier {
	return &concurrencyModifier{
		Position: fset.Position(n.Pos()),
		node:     n,
	}
}

func (cm *concurrencyModifier) FuncExpr() string {
	return fmt.Sprintf("_cffConcurrency%d_%d", cm.Position.Line, cm.Position.Column)
}

func (cm *concurrencyModifier) RetExpr() string {
	return "int"
}

func (cm *concurrencyModifier) GenImpl(p GenParams) error {
	modifierT := template.New(_concurrencyTmpl)
	modifierTmpl, err := modifierT.ParseFS(tmplFS, _modifierTmplDir)
	if err != nil {
		return err
	}
	if err := modifierTmpl.ExecuteTemplate(p.Writer, _concurrencyTmpl, cm); err != nil {
		return err
	}
	return nil
}

func (cm *concurrencyModifier) Node() ast.Node {
	return cm.node
}
