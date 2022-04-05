package internal

import (
	"bytes"
	"go/ast"
	"io"
	"text/template"
)

const (
	_parallelRootTmpl = "parallel.go.tmpl"
	_parallelTmplDir  = "templates/parallel/*"
)

func (g *generator) generateParallel(
	file *file,
	p *parallel,
	w io.Writer,
	addImports map[string]string,
	aliases map[string]struct{},
) error {
	// Tracks user-provided expressions used in the generated code.
	// We use this to ensure that the expressions are evaluated
	// in the order they were specified by the user.
	exprs := make(map[ast.Expr]struct{})
	fnMap := g.funcMap(file, addImports, aliases, exprs)
	t := template.New(_parallelRootTmpl).Funcs(fnMap)
	tmpl, err := t.ParseFS(tmplFS, _parallelTmplDir, _sharedTmplDir)
	if err != nil {
		return err
	}

	var b bytes.Buffer
	// Render the template to a staging buffer to understand how user provided
	// expressions are applied in generated code before writing the final
	// result.
	if err := tmpl.ExecuteTemplate(&b, _parallelRootTmpl, parallelTemplateData{
		Parallel: p,
	}); err != nil {
		return err
	}

	if _, err := io.WriteString(w, "func() (err error) {"); err != nil {
		return err
	}

	// Render variable assignments for user provided parameter expressions in
	// the order they were provided to cff.Flow.
	prologueT := template.New(_paramExprTmpl).Funcs(fnMap)
	prologueTmpl, err := prologueT.ParseFS(tmplFS, _prologueTmplDir)
	if err != nil {
		return err
	}
	if err := prologueTmpl.ExecuteTemplate(w, _paramExprTmpl, paramExprs(exprs)); err != nil {
		return err
	}
	if _, err := w.Write(b.Bytes()); err != nil {
		return err
	}
	if _, err := io.WriteString(w, "}()"); err != nil {
		return err
	}

	return nil
}

type parallelTemplateData struct {
	Parallel *parallel
}
