package main

import (
	"errors"
	"fmt"
	"go/token"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/cff/internal"
	"go.uber.org/cff/internal/flag"
	"go.uber.org/cff/internal/pkg"
	"go.uber.org/multierr"
)

type params struct {
	Files          []flag.InOutPair
	AutoInstrument bool
	GenMode        flag.Mode
	Quiet          bool
	ImportPath     string
}

func parseArgs(stderr io.Writer, args []string) (pkg.Loader, *params, error) {
	opts := params{
		GenMode: flag.BaseMode,
	}
	fset := flag.NewSet("cff")
	fset.SetOutput(stderr)
	fset.Usage = func() {
		fmt.Fprintln(fset.Output(), "usage: cff [options] importpath")
		fset.PrintDefaults()
	}

	fset.Var(flag.AsList(&opts.Files), "file", "By default, cff will process all Go files found inside the given Go package.\n"+
		"Pass -file=PATH one or more times to process only specific files and ignore all other files.\n"+
		"When cff processes a file, it generates sibling files with a _gen suffix next to the original files.\n"+
		"Use the form -file=PATH=OUTPUT to specify a different path for the output file.")

	fset.Var(&opts.GenMode, "genmode", "Use the specified CFF code generation mode.\n"+
		"Valid values are: base, modifier, source-map. Defaults to base.")

	fset.BoolVar(&opts.AutoInstrument, "auto-instrument", false,
		"Infer a name for tasks that do not specify cff.Instrument and opt-in "+
			"to instrumentation by default.")

	fset.BoolVar(&opts.Quiet, "quiet", false, "Print less output.")

	loader := _loaderFactory.RegisterFlags(fset)
	if err := fset.Parse(args); err != nil {
		return nil, nil, err
	}
	args = fset.Args()
	switch len(args) {
	case 0:
		return nil, nil, errors.New("please provide an import path")
	case 1:
		opts.ImportPath = args[0]
	default:
		return nil, nil, fmt.Errorf("too many import paths: %q", args)
	}

	return loader, &opts, nil
}

func main() {
	log.SetFlags(0) // don't include timestamps
	if err := run(os.Args[1:]); err != nil && !errors.Is(err, flag.ErrHelp) {
		log.Fatalf("%+v", err)
	}
}

var _loaderFactory pkg.LoaderFactory = &pkg.GoPackagesLoaderFactory{
	BuildFlags: []string{"-tags=cff"},
}

func run(args []string) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("You've encountered a CFFv2 bug! Please report this http://t.uber.com/cff-bug")
			panic(err)
		}
	}()

	loader, f, err := parseArgs(os.Stderr, args)
	if err != nil {
		return err
	}

	// For each --file, this is a mapping from FILE to OUTPUT.
	outputs := make(map[string]string)
	for _, file := range f.Files {
		if _, ok := outputs[file.Input]; ok {
			return fmt.Errorf(
				"invalid argument --file=%v: file already specified before", file)
		}
		outputs[file.Input] = file.Output
	}

	fset := token.NewFileSet()
	pkgs, err := loader.Load(fset, f.ImportPath)
	if err != nil {
		return fmt.Errorf("load packages: %w", err)
	}

	processor := internal.Processor{
		Fset:               fset,
		InstrumentAllTasks: f.AutoInstrument,
		GenMode:            f.GenMode,
	}

	// If --file was provided, only the requested files will be processed.
	// Otherwise all files will be processed.
	hadFiles := len(f.Files) > 0
	var processed, errored int
	for _, pkg := range pkgs {
		for i, path := range pkg.CompiledGoFiles {
			name := filepath.Base(path)
			output, ok := outputs[name]
			if hadFiles && !ok {
				// --file was provided and this file wasn't included.
				continue
			}

			if len(output) == 0 {
				// x/y/foo.go => x/y/foo_gen.go
				output = filepath.Join(
					filepath.Dir(path),
					// foo.go => foo + _gen.go
					strings.TrimSuffix(name, filepath.Ext(name))+"_gen.go",
				)
			}

			processed++
			if perr := processor.Process(pkg, pkg.Syntax[i], output); perr != nil {
				errored++
				err = multierr.Append(err, perr)
			}
		}
	}

	if !f.Quiet {
		log.Printf("Processed %d files with %d errors", processed, errored)
	}
	return err
}
