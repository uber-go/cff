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
)

type generator struct {
	fset *token.FileSet

	typeIDs    *typeutil.Map // map[types.Type]int
	nextTypeID int

	// File path to which generated code is written.
	outputPath string

	onlineScheduling bool
}

type generatorOpts struct {
	Fset             *token.FileSet
	OutputPath       string
	OnlineScheduling bool
}

func newGenerator(opts generatorOpts) *generator {
	return &generator{
		fset:             opts.Fset,
		typeIDs:          new(typeutil.Map),
		nextTypeID:       1,
		outputPath:       opts.OutputPath,
		onlineScheduling: opts.OnlineScheduling,
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
func (g *generator) generateFlow(file *file, f *flow, w io.Writer, addImports map[string]string) error {
	tmpl := parseTemplates(g.funcMap(file, addImports), _flowTmpl, _taskTmpl)
	return tmpl.ExecuteTemplate(w, _flowTmpl, flowTemplateData{
		Flow:   f,
		Online: g.onlineScheduling,
	})
}

func (g *generator) funcMap(file *file, addImports map[string]string) template.FuncMap {
	return template.FuncMap{
		"type": g.typePrinter(file, addImports),
		"typeName": func(t types.Type) string {
			// Report the name of the type without importing it.
			// Useful for comments.
			return types.TypeString(t, nil)
		},
		"typeHash": g.printTypeHash,
		"expr":     g.printExpr,
		"quote":    strconv.Quote,
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
	}
}

// typePrinter returns the qualifier for the type to form an identifier using that type, modifying addImports if the
// type refers to a package that is not already imported
func (g *generator) typePrinter(f *file, addImports map[string]string) func(types.Type) string {
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
				addImports[pkg.Path()] = pkg.Name()
				return pkg.Name()
			}

			// The type is defined in the same package
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
	Flow   *flow
	Online bool
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
