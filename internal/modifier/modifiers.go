// Package modifier implements modifier-based code generation
// for cff directives.
package modifier

import (
	"embed"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"io"
	"path/filepath"
	"strings"
	"text/template"
)

// ModifierTmplFS embeds templates for generating modifier implementations.
//
//go:embed templates/*
var ModifierTmplFS embed.FS

const (
	// TmplDir is the path pattern for modifier templates.
	TmplDir   = "templates/*"
	_funcTmpl = "func.go.tmpl"
)

// Modifier changes the existing code by doing two things.
//  1. It generates a function that corresponds to what the cff "API"s do by inspecting the
//     arguments.
//  2. It inline replaces the cff "API" call with calls to corrresponding generated function.
//
// Each call to cff "API" is translates to a modifier.
type Modifier interface {
	Expr() ast.Expr            // The ast Expr that produced this modifier.
	FuncExpr() string          // The name of the modifier-generated function.
	Provides() []ast.Expr      // The expressions that are provided to and returned by this modifier function.
	GenImpl(p GenParams) error // Generates the function body of the modifier-generated function.
}

// GenParams is the parameter for Modifiers to generate the body
// of their corresponding methods.
type GenParams struct {
	Writer  io.Writer
	FuncMap template.FuncMap
}

// Arg is a parameter of a Modifier function.
type Arg struct {
	Name    string
	Type    types.Type
	LastIdx bool
}

// ExprHash returns a unique identifier for an Expr.
func ExprHash(fset *token.FileSet, n ast.Expr) string {
	pos := fset.Position(n.Pos())
	return fmt.Sprintf("m%v%d_%d", TrimFilename(pos.Filename), pos.Line, pos.Column)
}

// Placeholder constructs a placeholder modifier that will be incrementally
// replaced by real modifier implementations. This should be removed when
// all modifiers are implemented.
func Placeholder(n ast.Expr) Modifier {
	return &placeholder{expr: n}
}

type placeholder struct {
	expr ast.Expr
}

func (p *placeholder) Expr() ast.Expr {
	return p.expr
}

func (p *placeholder) FuncExpr() string {
	return ""
}

func (p *placeholder) GenImpl(_ GenParams) error {
	return nil
}

func (p *placeholder) Provides() []ast.Expr {
	return nil
}

// TrimFilename is a utility function used for trimming just the file
// name without ".go" suffix to guarantee uniqueness of generated cff
// functions.
func TrimFilename(path string) string {
	return strings.ReplaceAll(strings.TrimSuffix(filepath.Base(path), ".go"), "_", "")
}
