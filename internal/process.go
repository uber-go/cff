package internal

import (
	"go/ast"
	"go/token"

	"code.uber.internal/go/importer"
)

// Processor processes CFF2 files.
type Processor struct {
	Fset               *token.FileSet
	InstrumentAllTasks bool
}

// Process processes a single CFF2 file.
func (p *Processor) Process(pkg *importer.Package, file *ast.File, outputPath string) error {
	c := newCompiler(compilerOpts{
		Fset:               p.Fset,
		Info:               pkg.TypesInfo,
		Package:            pkg.Types,
		InstrumentAllTasks: p.InstrumentAllTasks,
	})

	f, err := c.CompileFile(file, pkg)
	if err != nil {
		return err
	}

	g := newGenerator(generatorOpts{
		Fset:       p.Fset,
		Package:    pkg.Types,
		OutputPath: outputPath,
	})
	if err := g.GenerateFile(f); err != nil {
		return err
	}

	return nil
}
