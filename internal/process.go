package internal

import (
	"go/ast"
	"go/token"

	"go.uber.org/cff/internal/flag"
	"go.uber.org/cff/internal/pkg"
)

// Processor processes CFF files.
type Processor struct {
	Fset               *token.FileSet
	InstrumentAllTasks bool
	GenMode            flag.Mode
	RequireBuildTag    bool
}

// Process processes a single CFF file.
func (p *Processor) Process(pkg *pkg.Package, file *ast.File, outputPath string) error {
	c := newCompiler(compilerOpts{
		Fset:               p.Fset,
		Info:               pkg.TypesInfo,
		Package:            pkg.Types,
		InstrumentAllTasks: p.InstrumentAllTasks,
		RequireBuildTag:    p.RequireBuildTag,
	})

	f, err := c.CompileFile(file, pkg)
	if err != nil {
		return err
	}

	if p.GenMode == flag.ModifierMode {
		g := newGeneratorV2(generatorOpts{
			Fset:       p.Fset,
			Package:    pkg.Types,
			OutputPath: outputPath,
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
