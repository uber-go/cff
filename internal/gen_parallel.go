package internal

import (
	"io"
	"path/filepath"
	"strconv"
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
	t := template.New(_parallelRootTmpl).Funcs(g.parallelFuncMap(file, addImports, aliases))
	tmpl, err := t.ParseFS(tmplFS, _parallelTmplDir)
	if err != nil {
		return err
	}
	return tmpl.ExecuteTemplate(w, _parallelRootTmpl, parallelTemplateData{
		Parallel: p,
	})
}

func (g *generator) parallelFuncMap(file *file, addImports map[string]string, aliases map[string]struct{}) template.FuncMap {
	return template.FuncMap{
		"quote":   strconv.Quote,
		"rawExpr": g.printRawExpr,
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
