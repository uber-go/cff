package internal

import (
	"io"
	"path/filepath"
	"strconv"
	"text/template"
)

const (
	_parallelTmpl     = "parallel.go.tmpl"
	_parallelTaskTmpl = "parallel_task.go.tmpl"
)

func (g *generator) generateParallel(
	file *file,
	p *parallel,
	w io.Writer,
	addImports map[string]string,
	aliases map[string]struct{},
) error {
	tmpl := parseTemplates(g.parallelFuncMap(file, addImports, aliases), _parallelTmpl, _parallelTaskTmpl)
	return tmpl.ExecuteTemplate(w, _parallelTmpl, parallelTemplateData{
		Parallel: p,
	})
}

func (g *generator) parallelFuncMap(file *file, addImports map[string]string, aliases map[string]struct{}) template.FuncMap {
	return template.FuncMap{
		"expr":  g.printExpr,
		"quote": strconv.Quote,
		"import": func(importPath string) string {
			if names := file.Imports[importPath]; len(names) > 0 {
				// importPath exists in the file already.
				return names[0]
			}
			return printImportAlias(importPath, filepath.Base(importPath), addImports, aliases)
		},
	}
}

type parallelTemplateData struct {
	Parallel *parallel
}
