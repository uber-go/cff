package internal

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"text/template"

	"go.uber.org/cff/internal/modifier"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/types/typeutil"
)

type generatorv2 struct {
	fset *token.FileSet

	pkg *types.Package

	typeIDs    *typeutil.Map // map[types.Type]int
	nextTypeID int

	predIDs *typeutil.Map // map[types.Type]int

	// File path to which generated code is written.
	outputPath string
}

func newGeneratorV2(opts generatorOpts) *generatorv2 {
	return &generatorv2{
		fset:       opts.Fset,
		pkg:        opts.Package,
		typeIDs:    new(typeutil.Map),
		predIDs:    new(typeutil.Map),
		nextTypeID: 1,
		outputPath: opts.OutputPath,
	}
}

func (g *generatorv2) GenerateFile(f *file) error {
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

	// This tracks positioning information for the file
	posFile := g.fset.File(f.AST.Pos())

	addImports := make(map[string]string)
	aliases := make(map[string]struct{})

	// Track all aliases that already exist in the file.
	for _, names := range f.Imports {
		aliases[names[0]] = struct{}{}
	}

	fileModifiers := f.modifiers

	// Before generating code, sort the modifiers in order of appearance in the file.
	sort.Slice(fileModifiers, func(i, j int) bool {
		return fileModifiers[i].Expr().Pos() < fileModifiers[j].Expr().Pos()
	})

	// Build tags appear before the package clause.
	// Write those to the output with cff tags inverted.
	lastOff := posFile.Offset(f.AST.Package)
	if err := writeInvertedCffTag(&buff, bs[:lastOff]); err != nil {
		return err
	}

	for _, mod := range fileModifiers {
		if mod.FuncExpr() == "" {
			continue
		}
		// safe to ignore err from these because we're writing to bytes.Buffer.
		buff.Write(bs[lastOff:posFile.Offset(mod.Expr().Pos())])
		// Replace call sites with the function expression.
		buff.Write([]byte(mod.FuncExpr()))

		lastOff = posFile.Offset(mod.Expr().End())
	}

	// Write remaining code as-is.
	if _, err := buff.Write(bs[lastOff:]); err != nil {
		return err
	}

	// At the bottom of the file, generate the type definitions and modifier function
	// bodies.
	for _, mod := range fileModifiers {
		if err := mod.GenImpl(modifier.GenParams{
			Writer:  &buff,
			FuncMap: g.funcMap(f, addImports, aliases),
		}); err != nil {
			return err
		}

		// Insert a newline and space between each modifier generation.
		buff.Write([]byte("\n\n"))
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

	return os.WriteFile(g.outputPath, buff.Bytes(), 0o644)
}

func (g *generatorv2) funcMap(
	file *file,
	addImports map[string]string,
	aliases map[string]struct{},
) template.FuncMap {
	p := &exprPrinterv2{g: g}
	return template.FuncMap{
		"type": g.typePrinter(file, addImports, aliases),
		"typeName": func(t types.Type) string {
			// Report the name of the type without importing it.
			// Useful for comments.
			return types.TypeString(t, nil)
		},
		"quote": strconv.Quote,
		"import": func(importPath string) string {
			if names := file.Imports[importPath]; len(names) > 0 {
				return names[0]
			}
			res := printImportAlias(importPath, filepath.Base(importPath), addImports, aliases)
			return res
		},
		"expr":     p.printExpr,
		"typeHash": g.printTypeHash,
	}
}

func (g *generatorv2) posInfo(n ast.Node) *PosInfo {
	pos := g.fset.Position(n.Pos())
	posInfo := &PosInfo{
		File:   filepath.Join(g.pkg.Path(), filepath.Base(pos.Filename)),
		Line:   pos.Line,
		Column: pos.Column,
	}
	return posInfo
}

// typePrinter returns the qualifier for the type to form an identifier using that type, modifying addImports if the
// type refers to a package that is not already imported
func (g *generatorv2) typePrinter(f *file, addImports map[string]string, aliases map[string]struct{}) func(types.Type) string {
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

func (g *generatorv2) typeID(t types.Type) int {
	if i := g.typeIDs.At(t); i != nil {
		return i.(int)
	}

	id := g.nextTypeID
	g.nextTypeID++
	g.typeIDs.Set(t, id)
	return id
}

func (g *generatorv2) printTypeHash(t types.Type) string {
	return strconv.Itoa(g.typeID(t))
}

type exprPrinterv2 struct {
	g *generatorv2
}

func (p *exprPrinterv2) printExpr(e ast.Expr) string {
	pos := p.g.posInfo(e)
	return fmt.Sprintf("_%d_%d", pos.Line, pos.Column)
}
