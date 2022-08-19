package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"text/template"
)

const (
	_modifierTmplDir = "templates/modifiers/*"
	_concurrencyTmpl = "concurrency.go.tmpl"
)

type concurrencyModifier struct {
	Position token.Position

	node ast.Node
}

var _ modifier = (*concurrencyModifier)(nil)

func newConcurrencyModifier(n ast.Node, pos token.Position) *concurrencyModifier {
	return &concurrencyModifier{
		Position: pos,
		node:     n,
	}
}

func (cm *concurrencyModifier) FuncExpr() string {
	return fmt.Sprintf("_cffConcurrency%d_%d", cm.Position.Line, cm.Position.Column)
}

func (cm *concurrencyModifier) RetExpr() string {
	return "int"
}

func (cm *concurrencyModifier) GenImpl(p genParams) error {
	modifierT := template.New(_concurrencyTmpl)
	modifierTmpl, err := modifierT.ParseFS(tmplFS, _modifierTmplDir)
	if err != nil {
		return err
	}

	if err := modifierTmpl.ExecuteTemplate(p.writer, _concurrencyTmpl, cm); err != nil {
		return err
	}

	return nil
}

func (cm *concurrencyModifier) Node() ast.Node {
	return cm.node
}

type flowModifier struct {
	Position token.Position

	node ast.Node
}

var _ modifier = (*flowModifier)(nil)

func newFlowModifier(n ast.Node, pos token.Position) *flowModifier {
	return &flowModifier{
		Position: pos,
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

func (fm *flowModifier) GenImpl(p genParams) error {
	// TODO
	return nil
}

func (fm *flowModifier) Node() ast.Node {
	return fm.node
}
