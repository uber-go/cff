package modifier

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"text/template"
)

type taskModifier struct {
	Position token.Position

	args []ast.Expr
	fset *token.FileSet
	info *types.Info
	expr ast.Expr
}

var _ Modifier = (*taskModifier)(nil)

// NewTaskModifier creates a modifier for a cff.Task.
func NewTaskModifier(fset *token.FileSet, n ast.Expr, args []ast.Expr, info *types.Info) Modifier {
	return &taskModifier{
		Position: fset.Position(n.Pos()),
		args:     args,
		fset:     fset,
		info:     info,
		expr:     n,
	}
}

// FuncExpr returns the name of the cff.Task modifier function.
func (tm *taskModifier) FuncExpr() string {
	return fmt.Sprintf("_cffTask%d_%d", tm.Position.Line, tm.Position.Column)
}

func (tm *taskModifier) FuncArgs() []Arg {
	args := make([]Arg, len(tm.args))
	for i, arg := range tm.args {
		args[i] = Arg{
			Name:    ExprHash(tm.fset, arg),
			Type:    tm.info.TypeOf(arg),
			LastIdx: i == len(tm.args)-1,
		}
	}
	return args
}

func (tm *taskModifier) GenImpl(p GenParams) error {
	t := template.New(filepath.Join(_funcTmpl)).Funcs(p.FuncMap)
	mt, err := t.ParseFS(ModifierTmplFS, TmplDir)
	if err != nil {
		return err
	}
	if err := mt.ExecuteTemplate(p.Writer, _funcTmpl, tm); err != nil {
		return err
	}
	return nil
}

func (tm *taskModifier) Expr() ast.Expr {
	return tm.expr
}

func (tm *taskModifier) Provides() []ast.Expr {
	return tm.args
}
