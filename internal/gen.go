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

	"go.uber.org/cff/internal/templates"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/types/typeutil"
)

const (
	_flowTmpl = "flow.go.tmpl"
	_taskTmpl = "task.go.tmpl"
	_funcTmpl = "func.go.tmpl"
	_predTmpl = "predicate.go.tmpl"
)

type generator struct {
	fset *token.FileSet

	typeIDs    *typeutil.Map // map[types.Type]int
	nextTypeID int

	predIDs    *typeutil.Map // map[types.Type]int
	nextPredID int

	// File path to which generated code is written.
	outputPath string
}

type generatorOpts struct {
	Fset       *token.FileSet
	OutputPath string
}

func newGenerator(opts generatorOpts) *generator {
	return &generator{
		fset:       opts.Fset,
		typeIDs:    new(typeutil.Map),
		predIDs:    new(typeutil.Map),
		nextTypeID: 1,
		outputPath: opts.OutputPath,
	}
}

func (g *generator) GenerateFile(f *file) error {
	if len(f.Generators) == 0 {
		// Don't regenerate files that don't have directiveGenerators.
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
	aliases := make(map[string]struct{})  // map to record aliases that have been used during in code gen

	// Track all aliases that already exist in the file.
	for _, names := range f.Imports {
		aliases[names[0]] = struct{}{}
	}

	var lastOff int
	for _, gen := range f.Generators {
		// Everything from previous position up to this cff generator call.
		if _, err := buff.Write(bs[lastOff:posFile.Offset(gen.Pos())]); err != nil {
			return err
		}

		// Generate code for top-level CFF2 constructs and update the
		// addImports map.
		if err = gen.generate(
			genParams{
				generator:  g,
				file:       f,
				writer:     &buff,
				addImports: addImports,
				aliases:    aliases,
			},
		); err != nil {
			return err
		}

		lastOff = posFile.Offset(gen.End())
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

	// Remove build tag.
	for _, cg := range file.Comments {
		// Only consider comments before the "package" clause.
		if cg.Pos() >= file.Package {
			break
		}

		// Replace +build cff with +build !cff and add a generated comment that tells
		// Phabricator to skip showing this file in diffs.
		for _, c := range cg.List {
			if strings.TrimSpace(strings.TrimPrefix(c.Text, "//")) == "+build cff" {
				// Build tags must be followed by a blank line (https://golang.org/pkg/go/build/#hdr-Build_Constraints)
				// Trick Phabricator not to consider *this* file to be generated.
				c.Text = "// +build !cff\n\n// @g" + "enerated by CFF"
				break
			}
		}
	}

	buff.Reset()
	err = format.Node(&buff, fset, file)

	return ioutil.WriteFile(g.outputPath, buff.Bytes(), 0644)
}

// generateFlow runs the CFF template for the given flow and writes it to w, modifying addImports if the template
// requires additional imports to be added.
func (g *generator) generateFlow(file *file, f *flow, w io.Writer, addImports map[string]string, aliases map[string]struct{}) error {
	tmpl := parseTemplates(g.funcMap(file, addImports, aliases), _flowTmpl, _funcTmpl, _predTmpl, _taskTmpl)
	return tmpl.ExecuteTemplate(w, _flowTmpl, flowTemplateData{
		Flow: f,
	})
}

func (g *generator) funcMap(file *file, addImports map[string]string, aliases map[string]struct{}) template.FuncMap {
	return template.FuncMap{
		"type": g.typePrinter(file, addImports, aliases),
		"typeName": func(t types.Type) string {
			// Report the name of the type without importing it.
			// Useful for comments.
			return types.TypeString(t, nil)
		},
		"typeHash":    g.printTypeHash,
		"predHash":    g.printPredicateHash,
		"isPredicate": g.isPredicate,
		"expr":        g.printExpr,
		"quote":       strconv.Quote,
		"import": func(importPath string) string {
			if names := file.Imports[importPath]; len(names) > 0 {
				// importPath exists in the file already.
				return names[0]
			}
			return printImportAlias(importPath, filepath.Base(importPath), addImports, aliases)
		},
	}
}

// typePrinter returns the qualifier for the type to form an identifier using that type, modifying addImports if the
// type refers to a package that is not already imported
func (g *generator) typePrinter(f *file, addImports map[string]string, aliases map[string]struct{}) func(types.Type) string {
	return func(t types.Type) string {
		return types.TypeString(t, func(pkg *types.Package) string {
			for _, imp := range f.AST.Imports {
				ip, _ := strconv.Unquote(imp.Path.Value)

				if !isPackagePathEquivalent(pkg, ip) {
					continue
				}

				// Using a named import.
				if imp.Name != nil {
					return imp.Name.Name
				}

				// Unnamed imports use the package's name.
				return pkg.Name()
			}

			// The generated code needs a package (pkg) to be imported to form the qualifier, but it wasn't imported
			// by the user already and it isn't in this package (f.Package)
			if !isPackagePathEquivalent(pkg, f.Package.Types.Path()) {
				return printImportAlias(pkg.Path(), pkg.Name(), addImports, aliases)
			}

			// The type is defined in the same package
			return ""
		})
	}
}

func (g *generator) printTypeHash(t types.Type) string {
	return strconv.Itoa(g.typeID(t))
}

func (g *generator) printPredicateHash(p *predicate) string {
	return strconv.Itoa(g.predID(p))
}

func (g *generator) predID(p *predicate) int {
	t := p.SentinelOutput
	if i := g.predIDs.At(t); i != nil {
		return i.(int)
	}

	id := g.nextPredID
	g.nextPredID++
	g.predIDs.Set(t, id)
	return id
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

func (g *generator) isPredicate(t types.Type) bool {
	_, ok := t.(*predicateOutput)
	return ok
}

type flowTemplateData struct {
	Flow *flow
}

// Parses bindata-packaged templates by name in-order.
func parseTemplates(funcs template.FuncMap, paths ...string) *template.Template {
	var t *template.Template
	for _, path := range paths {
		contents := templates.MustAssetString(path)
		name := filepath.Base(path)

		// The first template is the root template and all others are
		// associated with it via .New. Only the first template needs
		// the functions attached to it.
		var tmpl *template.Template
		if t == nil {
			t = template.New(name).Funcs(funcs)
			tmpl = t
		} else {
			tmpl = t.New(name)
		}

		template.Must(tmpl.Parse(contents))
		// We can ignore the return value because the template is
		// associated with the tree.
	}

	return t
}

// printImportAlias processes the importPath and returns the alias that should be used for it while
// updating addImports and aliases maps.
func printImportAlias(importPath, alias string, addImports map[string]string, aliases map[string]struct{}) string {
	if name, ok := addImports[importPath]; ok {
		// importPath already exists in addImports, infer alias and return.
		if name == "" {
			name = filepath.Base(importPath)
		}
		return name
	}

	for {
		if _, ok := aliases[alias]; !ok {
			// alias is unique.
			aliases[alias] = struct{}{}
			if alias == filepath.Base(importPath) {
				addImports[importPath] = ""
			} else {
				addImports[importPath] = alias
			}
			return alias
		}
		// alias is already used, mangle a new one by prepending
		// an "_" until it is unique.
		alias = "_" + alias
	}
}

type genParams struct {
	generator  *generator
	file       *file
	writer     io.Writer
	addImports map[string]string
	aliases    map[string]struct{}
}

// directiveGenerator generates code for top-level CFF constructs.
type directiveGenerator interface {
	ast.Node
	// generate produces CFF code with side effects.
	generate(p genParams) error
}

type flowGenerator struct {
	flow *flow
}

// generate adapts cff.Flow code generation.
func (g flowGenerator) generate(p genParams) error {
	return p.generator.generateFlow(p.file, g.flow, p.writer, p.addImports, p.aliases)
}

func (g flowGenerator) End() token.Pos {
	return g.flow.End()
}

func (g flowGenerator) Pos() token.Pos {
	return g.flow.Pos()
}
