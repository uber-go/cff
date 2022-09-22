package modifier

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"text/template"
)

type resultsModifier struct {
	Position token.Position
	Results  []ast.Expr

	fset *token.FileSet
	info *types.Info
	expr ast.Expr
}

var _ Modifier = (*resultsModifier)(nil)

// NewResultsModifier returns a Modifier that corresponds to
// a cff.Results call.
func NewResultsModifier(fset *token.FileSet, n ast.Expr, results []ast.Expr, i *types.Info) Modifier {
	return &resultsModifier{
		Position: fset.Position(n.Pos()),
		Results:  results,
		fset:     fset,
		info:     i,
		expr:     n,
	}
}

func (rm *resultsModifier) FuncExpr() string {
	return fmt.Sprintf("_cffResults%d_%d", rm.Position.Line, rm.Position.Column)
}

func (rm *resultsModifier) FuncArgs() []Arg {
	args := make([]Arg, len(rm.Results))
	for i, res := range rm.Results {
		args[i] = Arg{
			Name:    ExprHash(rm.fset, res),
			Type:    rm.info.TypeOf(res),
			LastIdx: i == len(rm.Results)-1,
		}
	}
	return args
}

func (rm *resultsModifier) GenImpl(p GenParams) error {
	t := template.New(_funcTmpl).Funcs(p.FuncMap)
	mt, err := t.ParseFS(ModifierTmplFS, TmplDir)
	if err != nil {
		return err
	}
	if err := mt.ExecuteTemplate(p.Writer, _funcTmpl, rm); err != nil {
		return err
	}
	return nil
}

func (rm *resultsModifier) Expr() ast.Expr {
	return rm.expr
}

func (rm *resultsModifier) Provides() []ast.Expr {
	return rm.Results
}
