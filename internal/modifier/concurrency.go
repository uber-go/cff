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
	return fmt.Sprintf("_cffConcurrency%v_%d_%d", TrimFilename(cm.Position.Filename), cm.Position.Line, cm.Position.Column)
}

func (cm *concurrencyModifier) GenImpl(p GenParams) error {
	modifierT := template.New(_concurrencyTmplName)
	modifierTmpl, err := modifierT.ParseFS(ModifierTmplFS, _concurrencyTmplPath)
	if err != nil {
		return err
	}
	return modifierTmpl.ExecuteTemplate(p.Writer, _concurrencyTmplName, cm)
}

func (cm *concurrencyModifier) Expr() ast.Expr {
	return cm.expr
}

func (cm *concurrencyModifier) Provides() []ast.Expr {
	return []ast.Expr{cm.concurrency}
}
