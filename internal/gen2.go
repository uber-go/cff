package internal

import (
	"bytes"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"sort"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/types/typeutil"
)

type generatorv2 struct {
	fset *token.FileSet

	pkg *types.Package

	typeIDs    *typeutil.Map // map[types.Type]int
	nextTypeID int

	predIDs    *typeutil.Map // map[types.Type]int
	nextPredID int

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

	modifiers := f.modifiers

	// Before generating code, sort the modifiers in order of appearance in the file.
	sort.Slice(modifiers, func(i, j int) bool {
		return modifiers[i].Node().Pos() < modifiers[j].Node().Pos()
	})

	var lastOff int
	for _, mod := range modifiers {
		// safe to ignore err from these because we're writing to bytes.Buffer.
		buff.Write(bs[lastOff:posFile.Offset(mod.Node().Pos())])
		// Replace call sites with the function expression.
		buff.Write([]byte(mod.FuncExpr()))

		lastOff = posFile.Offset(mod.Node().End())
	}

	// Write remaining code as-is.
	if _, err := buff.Write(bs[lastOff:]); err != nil {
		return err
	}

	// At the bottom of the file, generate the type definitions and modifier function
	// bodies.
	for _, mod := range modifiers {
		if err := mod.GenImpl(genParams{
			generatorv2: g,
			file:        f,
			writer:      &buff,
			addImports:  addImports,
			aliases:     aliases,
		}); err != nil {
			return err
		}

		// Insert a newline between each modifier generations.
		buff.Write([]byte("\n"))
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

	return os.WriteFile(g.outputPath, buff.Bytes(), 0644)
}
