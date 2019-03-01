package internal

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/types/typeutil"
)

type generator struct {
	fset *token.FileSet

	typeIDs    *typeutil.Map // map[types.Type]int
	nextTypeID int

	// Directory to which generated code is written instead of the package
	// directory.
	outputDir string
}

func newGenerator(fset *token.FileSet, outputDir string) *generator {
	return &generator{
		fset:       fset,
		typeIDs:    new(typeutil.Map),
		nextTypeID: 1,
		outputDir:  outputDir,
	}
}

func (g *generator) GenerateFile(f *file) error {
	if len(f.Flows) == 0 {
		// Don't regenerate files that don't have flows.
		return nil
	}

	bs, err := ioutil.ReadFile(f.Filepath)
	if err != nil {
		return err
	}

	// Output buffer
	var buff bytes.Buffer

	// This tracks positioning information for the file.
	posFile := g.fset.File(f.AST.Pos())

	addImports := make(map[string]string) // import path -> name or empty for implicit name
	var lastOff int
	for _, flow := range f.Flows {
		// Everything from previous position up to this flow call.
		if _, err := buff.Write(bs[lastOff:posFile.Offset(flow.Pos())]); err != nil {
			return err
		}

		// Generate code for the flow.
		if err := g.generateFlow(f, flow, &buff, addImports); err != nil {
			return err
		}

		lastOff = posFile.Offset(flow.End())
	}

	// Write remaining code as-is.
	if _, err := buff.Write(bs[lastOff:]); err != nil {
		return err
	}

	// Parse the generated file and clean up.

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, posFile.Name(), buff.Bytes(), parser.ParseComments)
	if err != nil {
		// When there is a parsing error, we should output the file to a temporary file to help debugging
		// the template.
		tmpFile, tmpErr := ioutil.TempFile("", "*.go")
		if tmpErr == nil {
			if _, writeErr := buff.WriteTo(tmpFile); writeErr == nil {
				err = fmt.Errorf("%v\noutputted temporary file to %s", err, tmpFile.Name())
			}
		}

		return err
	}

	newImports := make([]string, 0, len(addImports))
	for imp := range addImports {
		newImports = append(newImports, imp)
	}
	sort.Strings(newImports)

	for _, importPath := range newImports {
		astutil.AddNamedImport(fset, file, addImports[importPath], importPath)
	}

	// The user code will have imports to cffImportPath but we should remove
	// them because it will be unused.
	if _, ok := f.UnnamedImports[cffImportPath]; ok {
		astutil.DeleteImport(fset, file, cffImportPath)
	}
	for _, name := range f.Imports[cffImportPath] {
		astutil.DeleteNamedImport(fset, file, name, cffImportPath)
	}

	// Remove build tag.
	for _, cg := range file.Comments {
		// Only consider comments before the "package" clause.
		if cg.Pos() >= file.Package {
			break
		}

		// Replace +build cff with +build !cff.
		for _, c := range cg.List {
			if strings.TrimSpace(strings.TrimPrefix(c.Text, "//")) == "+build cff" {
				c.Text = "// +build !cff"
				break
			}
		}
	}

	buff.Reset()
	err = format.Node(&buff, fset, file)

	// TODO(abg): Configurable output file name/template
	outputDir := g.outputDir
	if outputDir == "" {
		outputDir = filepath.Dir(f.Filepath)
	}
	newName := strings.TrimSuffix(filepath.Base(f.Filepath), ".go") + "_gen.go"
	return ioutil.WriteFile(filepath.Join(outputDir, newName), buff.Bytes(), 0644)
}

func (g *generator) generateFlow(file *file, f *flow, w io.Writer, addImports map[string]string) error {
	tmpl, err := template.New("cff").Funcs(template.FuncMap{
		"type":     g.typePrinter(file.AST),
		"typeHash": g.printTypeHash,
		"expr":     g.printExpr,
		"import": func(importPath string) string {
			if names := file.Imports[importPath]; len(names) > 0 {
				// already imported
				return names[0]
			}

			name, ok := addImports[importPath]
			if !ok {
				addImports[importPath] = ""
			}

			// TODO(abg): If the name is already taken, we will want to use
			// a named import. This can be done by having the compiler record
			// a list of unavailable names in the scope where cff.Flow was
			// called.
			if name == "" {
				name = filepath.Base(importPath)
			}

			return name
		},
	}).Parse(_tmpl)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, flowTemplateData{Flow: f})
}

func (g *generator) typePrinter(f *ast.File) func(types.Type) string {
	return func(t types.Type) string {
		return types.TypeString(t, func(pkg *types.Package) string {
			for _, imp := range f.Imports {
				ip, _ := strconv.Unquote(imp.Path.Value)
				if ip != pkg.Path() {
					continue
				}

				// Using a named import.
				if imp.Name != nil {
					return imp.Name.Name
				}

				// Unnamed imports use the package's name.
				return pkg.Name()
			}

			return ""
		})
	}
}

func (g *generator) printTypeHash(t types.Type) string {
	return strconv.Itoa(g.typeID(t))
}

func (g *generator) printExpr(e ast.Expr) string {
	var buff bytes.Buffer
	format.Node(&buff, g.fset, e)
	return buff.String()
}

func (g *generator) typeID(t types.Type) int {
	if i := g.typeIDs.At(t); i != nil {
		return i.(int)
	}

	id := g.nextTypeID
	g.nextTypeID++
	g.typeIDs.Set(t, id)
	return id
}

type flowTemplateData struct {
	Flow *flow
}

const _tmpl = `
{{- $context := import "context" -}}
{{- $schedule := .Flow.Schedule -}}
{{- $flow := .Flow -}}
{{- with .Flow -}}
func(ctx {{ $context }}.Context,
{{- if $flow.Scope -}}
{{- $tally := import "github.com/uber-go/tally" -}}
scope {{$tally}}.Scope,
{{- end -}}
{{- range .Inputs -}}
	v{{ typeHash .Type }} {{ type .Type }},
{{- end }}) (err error) {
	{{ range $schedIdx, $sched := $schedule }}
		if ctx.Err() != nil {
			{{- range $thisSchedIdx, $sched := $schedule -}}
				{{- range $task := $sched -}}
					{{- if (and ($task.Instrument) (ge $thisSchedIdx $schedIdx)) -}}
						scope.Counter("task.skipped").Inc(1)
					{{ end -}}
				{{- end -}}
			{{- end -}}
			{{ if $flow.Instrument -}}
				scope.Counter("taskflow.skipped").Inc(1)
			{{- end }}
			return ctx.Err()
		}
		{{ if eq (len .) 1 -}}
			{{/* If there is only one task, don't use errgroup. We're using
			     range but it'll be called only once. */}}
			{{- range . }}
				{{ template "taskResultVarDecl" . -}}
				{{ if .Predicate }}
					if {{ template "callTask" .Predicate }} {
				{{- end }}
				{{ template "taskResultList" . }} = {{ template "callTask" . }}
				{{ if .HasError -}}
					if err != nil {
						{{ if .RecoverWith -}}
							{{ if .Instrument -}}
								scope.Counter("task.error").Inc(1)
								scope.Counter("task.recovered").Inc(1)
							{{- end }} 
							{{ template "taskResultList" . }} = {{ range $i, $v := .RecoverWith -}}
								{{ if gt $i 0 }},{{ end }}{{ expr $v }}
							{{- end }}, nil
						{{- else -}}
							{{ if .Instrument -}}
								scope.Counter("task.error").Inc(1)
							{{- end }} 
							{{ if $flow.Instrument -}}
								scope.Counter("taskflow.error").Inc(1)
							{{- end }} 
							return err
						{{- end }}
					} {{ if .Instrument }} else {
						scope.Counter("task.success").Inc(1)
					} {{ end }}
				{{ end }}
				{{ if .Predicate }}
					} {{ if .Instrument }} else {
					scope.Counter("task.skipped").Inc(1)
				} {{ end }}
				{{ end }}
			{{ end }}
		{{ else -}}
			{{/* For >1 tasks, we need the variables to be in scope since we
			     can't return them. */}}
			{{- $once := printf "once%v" $schedIdx -}}
			{{- $wg := printf "wg%v" $schedIdx -}}
			{{- $serr := printf "err%v" $schedIdx -}}
			{{- $sync := import "sync" -}}
			var (
				{{ $wg }} {{ $sync }}.WaitGroup
				{{ $once }} {{ $sync }}.Once
			)

			{{ $wg }}.Add({{ len . }})
			{{ range . }}
				{{ template "taskResultVarDecl" . }}
				go func() {
					defer {{ $wg }}.Done()
					{{ if .Instrument -}}
						timer := scope.Timer("task.timing").Start()
						defer timer.Stop()
					{{- end }} 

					{{ if .Predicate }}
						if {{ template "callTask" .Predicate }} {
					{{ end }}
					{{ if .HasError }}var {{ $serr }} error
					{{ end -}}
					{{ template "taskResultList" . }} = {{ template "callTask" . }}
					{{ if .HasError -}}
						if {{ $serr }} != nil {
							{{ if .RecoverWith -}}
								{{ if .Instrument -}}
									scope.Counter("task.error").Inc(1)
									scope.Counter("task.recovered").Inc(1)
								{{- end }} 

								{{ template "taskResultList" . }} = {{ range $i, $v := .RecoverWith -}}
									{{ if gt $i 0 }},{{ end }}{{ expr $v }}
								{{- end }}, nil
							{{- else -}}
								{{ if .Instrument -}}
									scope.Counter("task.error").Inc(1)
								{{- end }} 
								{{ $once }}.Do(func() {
									err = {{ $serr }}
								})
							{{- end }}
						} {{ if .Instrument }} else {
							scope.Counter("task.success").Inc(1)
						} {{ end }}
					{{ end }}
					{{ if .Predicate }}
						}
					{{ end }}
				}()
			{{ end }}
			{{ $wg }}.Wait()
			if err != nil {
				{{ if $flow.Instrument -}}
					scope.Counter("taskflow.error").Inc(1)
				{{- end }} 
				return err
			}

			// Prevent variable unused errors.
			var (
				_ = &{{ $once }}
				{{ range . -}}
					{{ range .Outputs -}}
						_ = &v{{ typeHash .}}
					{{ end -}}
				{{ end }}
			)

		{{ end }}
	{{ end }}

	{{ range .Outputs }}
	*({{ expr .Node }}) = v{{ typeHash .Type }}
	{{- end }}

	{{ if $flow.Instrument -}}
	if err != nil {
		scope.Counter("taskflow.error").Inc(1)
	} else {
		scope.Counter("taskflow.success").Inc(1)
	}

	{{- end }} 

	return err
}({{ expr .Ctx }}{{ with $flow.Scope }}, {{ expr . }} {{end}}{{ range .Inputs }}, {{ expr .Node }}{{ end }})
{{- end -}}

{{- define "taskResultVarDecl" -}}
{{ if eq (len .Outputs) 1 -}}
	{{ range .Outputs }}var v{{ typeHash . }} {{ type . }}{{ end }}
{{- else -}}
	var (
		{{ range .Outputs -}}
			v{{ typeHash . }} {{ type . }}
		{{ end -}}
	)
{{- end }}
{{- end -}}

{{- define "taskResultList" -}}
{{- range $i, $t := .Outputs -}}
	{{ if gt $i 0 }},{{ end }}v{{ typeHash $t }}
{{- end }}{{ if .HasError }}, err{{ end }}
{{- end -}}

{{- define "callTask" -}}
	{{- expr .Node }}({{- if .WantCtx }}ctx,{{ end }} {{- range .Inputs }}v{{ typeHash . }}, {{- end }})
{{- end -}}
`
