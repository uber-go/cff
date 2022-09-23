package modifier

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"text/template"
)

type funcModifier struct {
	name     string
	modified ast.Expr
	provided []ast.Expr

	position token.Position
	fset     *token.FileSet
	info     *types.Info
}

const (
	// TaskName is the prefix for the name that replaces a cff.Task.
	TaskName = "_cffTask"
	// ResultsName is the prefix for the name that replaces a cff.Results.
	ResultsName = "_cffResults"
)

var _ Modifier = (*funcModifier)(nil)

// Params are the inputs to generating a modifier function.
type Params struct {
	// Name prefix that will replace the modified CFF directive in generated
	// code.
	Name string
	// Modified is the ast.Expr that is replaced by the modifier.
	Modified ast.Expr
	// Provided are arguments of modified CFF directive. These values are the
	// same as the values returned by the replacement modifier function.
	Provided []ast.Expr
	Fset     *token.FileSet
	Info     *types.Info
}

// NewModifier creates a basic modifier from the arguments of a CFF directive.
func NewModifier(p Params) Modifier {
	return &funcModifier{
		name:     p.Name,
		modified: p.Modified,
		provided: p.Provided,
		fset:     p.Fset,
		position: p.Fset.Position(p.Modified.Pos()),
		info:     p.Info,
	}
}

// FuncExpr returns the name of the modifier replacement function.
func (fm *funcModifier) FuncExpr() string {
	return fmt.Sprintf("%s%d_%d", fm.name, fm.position.Line, fm.position.Column)
}

func (fm *funcModifier) FuncArgs() []Arg {
	provided := make([]Arg, len(fm.provided))
	for i, arg := range fm.provided {
		provided[i] = Arg{
			Name:    ExprHash(fm.fset, arg),
			Type:    fm.info.TypeOf(arg),
			LastIdx: i == len(fm.provided)-1,
		}
	}
	return provided
}

func (fm *funcModifier) GenImpl(p GenParams) error {
	t := template.New(filepath.Join(_funcTmpl)).Funcs(p.FuncMap)
	mt, err := t.ParseFS(ModifierTmplFS, TmplDir)
	if err != nil {
		return err
	}
	if err := mt.ExecuteTemplate(p.Writer, _funcTmpl, fm); err != nil {
		return err
	}
	return nil
}

// Expr returns the ast.Expr replaced by the modifier function.
func (fm *funcModifier) Expr() ast.Expr {
	return fm.modified
}

// Provides returns the expression(s) supplied by the modifier.
func (fm *funcModifier) Provides() []ast.Expr {
	return fm.provided
}
