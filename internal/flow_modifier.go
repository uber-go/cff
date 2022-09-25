package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"text/template"

	"go.uber.org/cff/internal/modifier"
)

const (
	_flowModifierRootTmpl = "flow.go.tmpl"
)

type flowModifier struct {
	Position token.Position
	Flow     *flow

	fset *token.FileSet
	expr ast.Expr
	info *types.Info
}

var _ modifier.Modifier = (*flowModifier)(nil)

// NewFlowModifier returns a Modifier that corresponds to a cff.Flow call.
func NewFlowModifier(fset *token.FileSet, f *flow, n ast.Expr, i *types.Info) modifier.Modifier {
	return &flowModifier{
		Position: fset.Position(n.Pos()),
		Flow:     f,
		fset:     fset,
		expr:     n,
		info:     i,
	}
}

func (fm *flowModifier) FuncExpr() string {
	return fmt.Sprintf("_cffFlow%v_%d_%d",
		modifier.TrimFilename(fm.Position.Filename),
		fm.Position.Line,
		fm.Position.Column,
	)
}

type rootModifierParams struct {
	Ctx     ast.Expr
	CtxType types.Type

	Args   []rootArg
	Values []providedValues
}

type rootArg struct {
	Name  string
	Types []types.Type
}

type providedValues struct {
	ModifierID string
	Exprs      []ast.Expr
	LastIdx    int
}

func (fm *flowModifier) FuncArgs() rootModifierParams {
	var provided []providedValues
	args := make([]rootArg, len(fm.Flow.modifiers))

	for i, m := range fm.Flow.modifiers {
		if m.FuncExpr() == "" {
			args[i] = rootArg{}
			continue
		}

		var (
			typs  []types.Type
			exprs []ast.Expr
		)

		for _, p := range m.Provides() {
			typs = append(typs, fm.info.TypeOf(p))
			exprs = append(exprs, p)
		}

		args[i] = rootArg{
			Name:  modifier.ExprHash(fm.fset, m.Expr()),
			Types: typs,
		}

		provided = append(
			provided,
			providedValues{
				ModifierID: modifier.ExprHash(fm.fset, m.Expr()),
				Exprs:      exprs,
				LastIdx:    len(exprs) - 1,
			},
		)
	}

	return rootModifierParams{
		Ctx:     fm.Flow.Ctx,
		CtxType: fm.info.TypeOf(fm.Flow.Ctx),
		Args:    args,
		Values:  provided,
	}
}

func (fm *flowModifier) GenImpl(p modifier.GenParams) error {
	t := template.New(_flowModifierRootTmpl).Funcs(p.FuncMap)
	mt, err := t.ParseFS(modifier.ModifierTmplFS, modifier.TmplDir)
	if err != nil {
		return err
	}
	if err := mt.ExecuteTemplate(p.Writer, _flowModifierRootTmpl, fm); err != nil {
		return err
	}
	return nil
}

func (fm *flowModifier) Expr() ast.Expr {
	return fm.expr
}

func (fm *flowModifier) Provides() []ast.Expr {
	// cff.Flow Root modifier constructs its arguments from all modifiers
	// associated with it during CFF compilation.
	return nil
}
