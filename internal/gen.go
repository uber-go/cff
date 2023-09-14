package internal

import (
	"bytes"
	"embed"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"text/template"

	"go.uber.org/cff/internal/flag"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/types/typeutil"
)

const (
	_flowRootTmpl    = "flow.go.tmpl"
	_flowTmplDir     = "templates/flow/*"
	_paramExprTmpl   = "param_expr.go.tmpl"
	_prologueTmplDir = "templates/prologue/*"
	_sharedTmplDir   = "templates/shared/*"
)

//go:embed templates/*
var tmplFS embed.FS

type generator struct {
	fset *token.FileSet

	pkg *types.Package

	typeIDs    *typeutil.Map // map[types.Type]int
	nextTypeID int

	predIDs    *typeutil.Map // map[types.Type]int
	nextPredID int

	// File path to which generated code is written.
	outputPath string

	// magic token to "reset" line directives.
	magic string

	sourceMapped bool
}

type generatorOpts struct {
	Fset       *token.FileSet
	Package    *types.Package
	OutputPath string
	RandSrc    rand.Source
	GenMode    flag.Mode
}

func newGenerator(opts generatorOpts) *generator {
	if opts.RandSrc == nil {
		opts.RandSrc = rand.NewSource(rand.Int63())
	}
	return &generator{
		fset:         opts.Fset,
		pkg:          opts.Package,
		typeIDs:      new(typeutil.Map),
		predIDs:      new(typeutil.Map),
		nextTypeID:   1,
		outputPath:   opts.OutputPath,
		magic:        fmt.Sprintf("CFF_MAGIC_TOKEN=%d\n", rand.New(opts.RandSrc).Int()),
		sourceMapped: opts.GenMode == flag.SourceMapMode,
	}
}

func (g *generator) GenerateFile(f *file) error {
	if len(f.Generators) == 0 {
		// Don't regenerate files that don't have directiveGenerators.
		return nil
	}

	bs, err := os.ReadFile(f.Filepath)
	if err != nil {
		return err
	}

	// Output buffer
	var buff bytes.Buffer

	// This tracks positioning information for the file.
	posFile := g.fset.File(f.AST.Pos())

	// Annotate with line directive to original source on top of the file.
	if g.sourceMapped {
		fmt.Fprintf(&buff, "//line %v:1\n", filepath.Base(posFile.Name()))
	}

	addImports := make(map[string]string) // import path -> name or empty for implicit name
	aliases := make(map[string]struct{})  // map to record aliases that have been used during in code gen

	// Track all aliases that already exist in the file.
	for _, names := range f.Imports {
		aliases[names[0]] = struct{}{}
	}

	// Build tags appear before the package clause.
	// Write those to the output with cff tags inverted.
	lastOff := posFile.Offset(f.AST.Package)
	if err := writeInvertedCffTag(&buff, bs[:lastOff]); err != nil {
		return err
	}

	for _, gen := range f.Generators {
		// Everything from previous position up to this cff generator call.
		if _, err := buff.Write(bs[lastOff:posFile.Offset(gen.Pos())]); err != nil {
			return err
		}

		// Generate code for top-level cff constructs and update the
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
		tmpFile, tmpErr := os.CreateTemp("", "*.go")
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

	// Add the newly added imports to the file first.
	for _, importPath := range newImports {
		astutil.AddNamedImport(fset, file, addImports[importPath], importPath)
	}

	buff.Reset()
	// Format the node and write it to the buffer.
	if err := format.Node(&buff, fset, file); err != nil {
		return err
	}

	if g.sourceMapped {
		// get new file content with replaced magic tokens.
		var newBuff bytes.Buffer
		if err := g.resetMagicTokens(&newBuff, &buff); err != nil {
			return err
		}
		return os.WriteFile(g.outputPath, newBuff.Bytes(), 0o644)
	}
	return os.WriteFile(g.outputPath, buff.Bytes(), 0o644)
}

func (g *generator) resetMagicTokens(w io.Writer, buff *bytes.Buffer) error {
	// After formatting, we search the final output for the magic token to replace it with line
	// directives. We do this after the formatting because the line numbers of the generated
	// code can change due to formatting, so any line directives that point to the same file
	// can break after formatting.

	// First write the file to FS.
	bb := buff.Bytes()
	if err := os.WriteFile(g.outputPath, bb, 0o644); err != nil {
		return err
	}

	// Then we need to re-parse this to search for the magic token and
	// replace it with the line directives to "reset" the line directives
	// that point back to the original cff source.
	// Without these, all the generated code will point to arbitrary and/or
	// non-existent locations in the original source.
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, g.outputPath, bb, parser.ParseComments)
	if err != nil {
		return err
	}

	var magicList []*ast.CommentGroup
	for _, cg := range file.Comments {
		if cg.Text() == g.magic {
			magicList = append(magicList, cg)
		}
	}
	var offset int
	outputPath := filepath.Base(g.outputPath)
	for _, magic := range magicList {
		pos := fset.PositionFor(magic.Pos(), false)
		if _, err := w.Write(bb[offset:pos.Offset]); err != nil {
			return err
		}
		fmt.Fprintf(w, "/*line %v:%d*/", outputPath, pos.Line+1)
		offset = fset.PositionFor(magic.End(), false).Offset
	}
	// Write remaining code as-is.
	_, err = w.Write(bb[offset:])
	return err
}

// generateFlow runs the cff template for the given flow and writes it to w, modifying addImports if the template
// requires additional imports to be added.
func (g *generator) generateFlow(file *file, f *flow, w io.Writer, addImports map[string]string, aliases map[string]struct{}) error {
	// Tracks user-provided expressions used in the generated code.
	// We use this to ensure that the expressions are evaluated
	// in the order they were specified by the user.
	exprs := make(map[ast.Expr]struct{})

	fnMap := g.funcMap(file, addImports, aliases, exprs)
	t := template.New(_flowRootTmpl).Funcs(fnMap)
	tmpl, err := t.ParseFS(tmplFS, _flowTmplDir, _sharedTmplDir)
	if err != nil {
		return err
	}

	var b bytes.Buffer
	// Render the template to a staging buffer to understand how user provided
	// expressions are applied in generated code before writing the final
	// result.
	if err := tmpl.ExecuteTemplate(&b, _flowRootTmpl, flowTemplateData{
		Flow: f,
	}); err != nil {
		return err
	}
	if _, err := io.WriteString(w, "func() (err error) {\n"); err != nil {
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
	if g.sourceMapped {
		// Annotate with line directives after we're done generating code.
		// Get the expression's End position and find the associated line.
		endPos := g.fset.Position(f.End())
		// -1 because this is a line above the closing }().
		fmt.Fprintf(w, "/*line %v:%d*/", filepath.Base(f.PosInfo.File), endPos.Line-1)
	}

	if _, err := io.WriteString(w, "}()"); err != nil {
		return err
	}
	return nil
}

func (g *generator) funcMap(
	file *file,
	addImports map[string]string,
	aliases map[string]struct{},
	exprs map[ast.Expr]struct{},
) template.FuncMap {
	p := &exprPrinter{exprs: exprs, g: g, sourceMapped: g.sourceMapped}
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
		"expr":        p.printExpr,
		"rawExpr":     g.printRawExpr,
		"lineDir":     p.printLineDir,
		"magic":       g.printMagic,
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

func (g *generator) printMagic() string {
	if !g.sourceMapped {
		return ""
	}
	return fmt.Sprintf("\n// %v", g.magic)
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

func (g *generator) posInfo(n ast.Node) *PosInfo {
	pos := g.fset.Position(n.Pos())
	posInfo := &PosInfo{
		File:   filepath.Join(g.pkg.Path(), filepath.Base(pos.Filename)),
		Line:   pos.Line,
		Column: pos.Column,
	}
	return posInfo
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

// printRawExpr prints an ast.Expr.
//
// printExpr cannot be used in place of printRawExpr because printExpr prints
// a variable for user provided expressions. In the prologue template, the
// initial variables must have raw expressions assigned to them.
func (g *generator) printRawExpr(e ast.Expr) string {
	var buff bytes.Buffer
	if err := format.Node(&buff, g.fset, e); err != nil {
		// format.Node can fail with an in-memory buffer
		// only if the ast.Expr is invalid.
		// We are certain that that's not possible here
		// because the node was already parsed and type-checked.
		// So if this fails, we should panic.
		panic(fmt.Sprintf("failed to format node: %v", err))
	}
	return buff.String()
}

// exprPrinter is a convenience type to (1) avoid declaring a non-trival
// function inside funcMap's scope and (2) pass through state from funcMap
// to printExpr to avoid tracking unecessary state on the global generator
// object.
type exprPrinter struct {
	exprs        map[ast.Expr]struct{}
	g            *generator
	sourceMapped bool
}

// printExpr called on a user provided expression returns a variable name hash
// for the expression.
//
// Expressions called by printExpr are recorded on the generator and assigned
// to variables at the start of cff.Flow generated code to preserve the order
// in which they were provided to the flow.
//
// When called on a non-user provided expression the expression itself is
// printed.
//
// This is necessary as code generation updates to prevent variable shadowing
// in cff.Flow changed the order in which user provided expressions were
// invoked (GO-1098).
func (p *exprPrinter) printExpr(e ast.Expr) string {
	if ident, ok := e.(*ast.Ident); ok && ident.Name == "nil" {
		// In the generated code, we cannot do the following,
		// because nil is untyped by default.
		//
		//   _12_34 := nil
		//
		// In lieu of trying to specify a type for it,
		// we'll use nil expressions as-is.
		return ident.Name
	}
	if !e.Pos().IsValid() {
		// This expression was not user provided (e.g. implied instrument
		// names) there is no assigned variable that can be used, print
		// the expression directly.
		return p.g.printRawExpr(e)
	}
	p.exprs[e] = struct{}{}
	pos := p.g.posInfo(e)
	return fmt.Sprintf("_%d_%d", pos.Line, pos.Column)
}

// printLineDir prints a line directive for an ast.Expr. It looks up the
// position of the expression in the original source, so that the generated
// code can be mapped back to the original source.
func (p *exprPrinter) printLineDir(e ast.Expr) string {
	if !p.sourceMapped {
		return ""
	}
	pos := p.g.posInfo(e)
	return fmt.Sprintf("/*line %v:%d:%d*/", filepath.Base(pos.File), pos.Line, pos.Column)
}

// paramExprs returns a slice of user provided expressions sorted such that
// earlier positioned expression appear first.
func paramExprs(provided map[ast.Expr]struct{}) []ast.Expr {
	var exprs []ast.Expr

	for expr := range provided {
		// A user provided expr cannot have a 0 Line nor 0 Column, so filter
		// these non user provided entries out. cff implied instrumentation
		// creates expressions that are not user provided and must be handled.
		if !expr.Pos().IsValid() {
			continue
		}
		exprs = append(exprs, expr)
	}

	sort.Slice(exprs, func(i, j int) bool {
		return exprs[i].Pos() < exprs[j].Pos()
	})
	return exprs
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

// directiveGenerator generates code for top-level cff constructs.
type directiveGenerator interface {
	ast.Node
	// generate produces cff code with side effects.
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

type parallelGenerator struct {
	parallel *parallel
}

func (g parallelGenerator) generate(p genParams) error {
	return p.generator.generateParallel(p.file, g.parallel, p.writer, p.addImports, p.aliases)
}

func (g parallelGenerator) End() token.Pos {
	return g.parallel.End()
}

func (g parallelGenerator) Pos() token.Pos {
	return g.parallel.Pos()
}
