package modifier

import (
	"fmt"
	"go/ast"
	"go/token"
	"text/template"
)

const (
	_concurrencyTmplPath = "templates/concurrency.go.tmpl"
	_concurrencyTmplName = "concurrency.go.tmpl"
)

type concurrencyModifier struct {
	Position token.Position

	expr        ast.Expr
	concurrency ast.Expr
}

var _ Modifier = (*concurrencyModifier)(nil)

// NewConcurrencyModifier returns a Modifier that corresponds to
// a cff.Concurrency call.
func NewConcurrencyModifier(fset *token.FileSet, n ast.Expr, concurrency ast.Expr) Modifier {
	return &concurrencyModifier{
		Position:    fset.Position(n.Pos()),
		expr:        n,
		concurrency: concurrency,
	}
}

func (cm *concurrencyModifier) FuncExpr() string {
	return fmt.Sprintf("_cffConcurrency%d_%d", cm.Position.Line, cm.Position.Column)
}

func (cm *concurrencyModifier) GenImpl(p GenParams) error {
	modifierT := template.New(_concurrencyTmplName)
	modifierTmpl, err := modifierT.ParseFS(ModifierTmplFS, _concurrencyTmplPath)
	if err != nil {
		return err
	}
	if err := modifierTmpl.ExecuteTemplate(p.Writer, _concurrencyTmplName, cm); err != nil {
		return err
	}
	return nil
}

func (cm *concurrencyModifier) Expr() ast.Expr {
	return cm.expr
}

func (cm *concurrencyModifier) Provides() []ast.Expr {
	return []ast.Expr{cm.concurrency}
}
