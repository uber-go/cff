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
{{- $fmt := import "fmt" -}}
{{- $schedule := .Flow.Schedule -}}
{{- $flow := .Flow -}}
{{- with .Flow -}}
func(ctx {{ $context }}.Context,
{{- if $flow.ObservabilityEnabled -}}
{{- $tally := import "github.com/uber-go/tally" -}}
{{- $zap := import "go.uber.org/zap" -}}
scope {{ $tally }}.Scope,
logger *{{ $zap }}.Logger,
{{- end -}}
{{- range .Inputs -}}
	v{{ typeHash .Type }} {{ type .Type }},
{{- end }}) (err error) {
	{{- if $flow.Instrument -}}
	flowTags := map[string]string{"name": {{ expr $flow.Instrument.Name }}}
	flowTagsMutex := new(sync.Mutex)
	{{- end -}}
	{{ range $schedIdx, $sched := $schedule }}
		if ctx.Err() != nil {
			{{- range $thisSchedIdx, $sched := $schedule -}}
				{{- range $idx, $task := $sched -}}
					{{- if (and ($task.Instrument) (ge $thisSchedIdx $schedIdx)) -}}
						{{- $taskName := expr $task.Instrument.Name }}
						{{- $tagsVar := print "s" $thisSchedIdx "t" $idx "Tags" }}
						{{ $tagsVar }} := map[string]string{"name": {{ $taskName }}}
						scope.Tagged({{ $tagsVar }}).Counter("task.skipped").Inc(1)
						logger.Debug("task skipped",
									 zap.String("name", {{ $taskName }}),
								     zap.Error(ctx.Err()),
									)
					{{ end -}}
				{{- end -}}
			{{- end -}}
			{{ if $flow.Instrument -}}
				scope.Tagged(flowTags).Counter("taskflow.skipped").Inc(1)
				logger.Debug("taskflow skipped", zap.String("name", {{ expr $flow.Instrument.Name }}))
			{{- end }}
			return ctx.Err()
		}

		{{- $once := printf "once%v" $schedIdx -}}
		{{- $wg := printf "wg%v" $schedIdx -}}
		{{- $sync := import "sync" }}
		var (
			{{ $wg }} {{ $sync }}.WaitGroup
			{{ $once }} {{ $sync }}.Once
		)

		{{ $wg }}.Add({{ len . }})
		{{ range . }}
			{{- $serr := printf "err%v" .Serial -}}
			{{ template "taskResultVarDecl" . }}
			go func() {
				defer {{ $wg }}.Done()
				{{ if .Instrument -}}
					tags := map[string]string{"name": {{ expr .Instrument.Name }}}
					timer := scope.Tagged(tags).Timer("task.timing").Start()
					defer timer.Stop()
				{{- end }}
				defer func() {
					recovered := recover()
					if recovered != nil {
						{{ $once }}.Do(func() {
							recoveredErr := {{ $fmt }}.Errorf("task panic: %v", recovered)
							{{ if .Instrument -}}
							scope.Tagged(map[string]string{"name": {{ expr .Instrument.Name }}}).Counter("task.panic").Inc(1)
							logger.Error("task panic", 
								zap.String("name", {{ expr .Instrument.Name }}),
								zap.Stack("stack"),
								zap.Error(recoveredErr))
							{{- end }}
							err = recoveredErr
						})
					}
				}()

				{{ if .Predicate }}
					if {{ template "callTask" .Predicate }} {
				{{ end }}
				{{ template "taskResultList" . }}{{ if or .HasError (len .Outputs) }} = {{ end }}{{ template "callTask" . }}
				{{ if .HasError -}}
					if {{ $serr }} != nil {
						{{ if .FallbackWith -}}
							{{ if .Instrument -}}
								{{ if $flow.Instrument -}}
								flowTagsMutex.Lock()
								flowTags["failedTask"] = {{ expr .Instrument.Name }}
								flowTagsMutex.Unlock()
								{{- end }}
								scope.Tagged(tags).Counter("task.error").Inc(1)
								scope.Tagged(tags).Counter("task.recovered").Inc(1)
								logger.Error("task error recovered",
											 zap.String("name", {{ expr .Instrument.Name }}),
											 zap.Error({{ $serr }}),
											)
							{{- end }}

							{{ template "taskResultList" . }} = {{ range $i, $v := .FallbackWith -}}
								{{ if gt $i 0 }},{{ end }}{{ expr $v }}
							{{- end }}, nil
						{{- else -}}
							{{ if .Instrument -}}
								{{ if $flow.Instrument -}}
								flowTagsMutex.Lock()
								flowTags["failedTask"] = {{ expr .Instrument.Name }}
								flowTagsMutex.Unlock()
								{{- end }}
								scope.Tagged(tags).Counter("task.error").Inc(1)
							{{- end }}
							{{ $once }}.Do(func() {
								err = {{ $serr }}
							})
						{{- end }}
					} {{ if .Instrument }} else {
						scope.Tagged(tags).Counter("task.success").Inc(1)
						logger.Debug("task succeeded", zap.String("name", {{ expr .Instrument.Name }}))
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
				scope.Tagged(flowTags).Counter("taskflow.error").Inc(1)
			{{- end }}
			return err
		}

		// Prevent variable unused errors.
		var (
			{{- if $flow.Instrument -}}
			_ = flowTagsMutex
			{{ end }}
			_ = &{{ $once }}
			{{ range . -}}
				{{ range .Outputs -}}
					_ = &v{{ typeHash .}}
				{{ end -}}
			{{ end }}
		)
	{{ end }}

	{{ range .Outputs }}
	*({{ expr .Node }}) = v{{ typeHash .Type }}
	{{- end }}

	{{ if $flow.Instrument -}}
	if err != nil {
		scope.Tagged(flowTags).Counter("taskflow.error").Inc(1)
	} else {
		scope.Tagged(flowTags).Counter("taskflow.success").Inc(1)
		logger.Debug("taskflow succeeded", zap.String("name", {{ expr $flow.Instrument.Name }}))
	}

	{{- end }}

	return err
}({{ expr .Ctx }}{{ if $flow.ObservabilityEnabled }}, {{ expr $flow.Metrics }}, {{ expr $flow.Logger }} {{ end }}{{ range .Inputs }}, {{ expr .Node }}{{ end }})
{{- end -}}

{{- define "taskResultVarDecl" -}}
{{ range .Outputs }}
var v{{ typeHash . }} {{ type . }}
{{- end }}
{{ if .HasError }}var {{ printf "err%d" .Serial }} error{{ end }}
{{- end -}}

{{- define "taskResultList" -}}
{{- range $i, $t := .Outputs -}}
	{{ if gt $i 0 }},{{ end }}v{{ typeHash $t }}
{{- end }}{{ if .HasError }}{{ if len .Outputs }}, {{ end }}{{ printf "err%d" .Serial }}{{ end }}
{{- end -}}

{{- define "callTask" -}}
	{{- expr .Node }}({{- if .WantCtx }}ctx,{{ end }} {{- range .Inputs }}v{{ typeHash . }}, {{- end }})
{{- end -}}
`
