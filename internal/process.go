package internal

import (
	"go/ast"
	"go/token"

	"go.uber.org/cff/mode"
	"code.uber.internal/go/importer"
)

// Processor processes CFF2 files.
type Processor struct {
	Fset               *token.FileSet
	InstrumentAllTasks bool
	GenMode            mode.GenerationMode
	UseV2Gen           bool
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

	if p.UseV2Gen {
		g := newGeneratorV2(generatorOpts{
			Fset:       p.Fset,
			Package:    pkg.Types,
			OutputPath: outputPath,
			GenMode:    p.GenMode,
		})
		if err := g.GenerateFile(f); err != nil {
			return err
		}
	} else {
		g := newGenerator(generatorOpts{
			Fset:       p.Fset,
			Package:    pkg.Types,
			OutputPath: outputPath,
			GenMode:    p.GenMode,
		})
		if err := g.GenerateFile(f); err != nil {
			return err
		}
	}

	return nil
}
