package main

import (
	"errors"
	"fmt"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/cff/internal"
	"github.com/jessevdk/go-flags"
	"go.uber.org/multierr"
	"golang.org/x/tools/go/packages"
)

type options struct {
	Archives           []string `long:"archive" value-name:"IMPORTPATHS=IMPORTMAP=FILE=EXPORT"`
	Files              []file   `long:"file" value-name:"FILE[=OUTPUT]"`
	InstrumentAllTasks bool     `long:"instrument-all-tasks"`
	Sources            []string `long:"source"`
	StdlibRoot         string   `long:"stdlibroot"`
	Args               struct {
		ImportPath string `positional-arg-name:"importPath"`
	} `positional-args:"yes" required:"yes"`
	Quiet bool `long:"quiet"`

	// Temporary flag to gradually onboard users to online scheduling.
	NoOnlineScheduling bool `long:"no-online-scheduling"`
}

// file is the value of the --file option.
//
// Two forms are supported:
//
//  --file=NAME
//  --file=NAME=OUTPUT
//
// For example,
//
//  --file=foo.go=_gen/foo.go --file=bar.go=_gen/bar.go
type file struct {
	Name       string // NAME portion of the argument
	OutputPath string // OUTPUT portion of the argument
}

func (f *file) String() string {
	if len(f.OutputPath) == 0 {
		return f.Name
	}
	return f.Name + "=" + f.OutputPath
}

func (f *file) UnmarshalFlag(name string) error {
	var output string
	if i := strings.IndexByte(name, '='); i >= 0 {
		name, output = name[:i], name[i+1:]
	}

	if len(name) == 0 {
		return errors.New("file name cannot be empty")
	}

	f.Name = name
	f.OutputPath = output
	return nil
}

func newCLIParser() (*flags.Parser, *options) {
	var opts options
	parser := flags.NewParser(&opts, flags.HelpFlag)
	parser.Name = "cff"

	// This is more readable than embedding the descriptions in the options
	// above.
	parser.FindOptionByLongName("archive").Description =
		"Use the given archive FILE for import path IMPORTMAP when parsing the " +
			"source files. IMPORTPATHS is a colon-separated list of import paths; " +
			"IMPORTMAP is the actual import path of the library this archive " +
			"holds; FILE is the path to the archive file; EXPORT is the path to " +
			"the corresponding export file. Currently, IMPORTPATHS and EXPORT " +
			"arguments are ignored."
	parser.FindOptionByLongName("file").Description =
		"Process only the file named NAME inside the package. All other files " +
			"will be ignored. NAME must be the name of the file, not the full path. " +
			"Optionally, OUTPUT may be provided as the path to which the generated " +
			"code for FILE will be written. By default, this defaults to adding a " +
			"_gen suffix to the file name."
	parser.FindOptionByLongName("instrument-all-tasks").Description =
		"Infer a name for tasks that do not specify cff.Instrument and opt-in " +
			"to instrumentation by default."
	parser.FindOptionByLongName("source").Description =
		"When using archives to parse the source code, specifies the filepaths to " +
			"all Go code in the package, so that CFF can parse the entire " +
			"package."
	parser.FindOptionByLongName("stdlibroot").Description =
		"When using archives to parse the source code, specifies the path containing " +
			"archive files for the Go standard library."
	parser.FindOptionByLongName("no-online-scheduling").Description =
		"If set, CFF2 will use an offline scheduling algorithm instead of " +
			"online scheduling at compile time. This will be deleted in the future."

	parser.Args()[0].Description = "Import path of a package containing CFF flows."

	return parser, &opts
}

func main() {
	log.SetFlags(0) // don't include timestamps
	if err := run(os.Args[1:]); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run(args []string) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("You've encountered a CFFv2 bug! Please report this http://t.uber.com/cff-bug")
			panic(err)
		}
	}()

	parser, f := newCLIParser()
	if _, err := parser.ParseArgs(args); err != nil {
		return err
	}

	// For each --file, this is a mapping from FILE to OUTPUT.
	outputs := make(map[string]string)
	for _, file := range f.Files {
		if _, ok := outputs[file.Name]; ok {
			return fmt.Errorf(
				"invalid argument --file=%v: file already specified before", file)
		}
		outputs[file.Name] = file.OutputPath
	}

	archives := make([]internal.Archive, len(f.Archives))
	for i, archive := range f.Archives {
		a, err := parseArchive(archive)
		if err != nil {
			return fmt.Errorf("invalid argument --archive=%q: %v", archive, err)
		}
		archives[i] = a
	}

	fset := token.NewFileSet()
	pkgs, err := loadPackages(internal.LoadParams{
		Fset:       fset,
		ImportPath: f.Args.ImportPath,
		Srcs:       f.Sources,
		StdlibRoot: f.StdlibRoot,
		Archives:   archives,
	})
	if err != nil {
		return err
	}

	processor := internal.Processor{
		Fset:               fset,
		InstrumentAllTasks: f.InstrumentAllTasks,
		OnlineScheduling:   !f.NoOnlineScheduling,
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

func loadPackages(p internal.LoadParams) ([]*internal.Package, error) {
	if len(p.Archives) > 0 {
		return internal.PackagesArchive(p)
	}

	mode := packages.NeedName |
		packages.NeedFiles |
		packages.NeedCompiledGoFiles |
		packages.NeedImports |
		packages.NeedDeps |
		packages.NeedTypes |
		packages.NeedSyntax |
		packages.NeedTypesInfo |
		packages.NeedTypesSizes
	pkgs, err := packages.Load(&packages.Config{
		Mode:       mode,
		Fset:       p.Fset,
		BuildFlags: []string{"-tags=cff"},
	}, p.ImportPath)

	if err != nil {
		return nil, fmt.Errorf("could not load packages: %v", err)
	}

	if len(pkgs) == 0 {
		return nil, errors.New("no packages found")
	}

	ipkgs := make([]*internal.Package, 0, len(pkgs))
	for _, pkg := range pkgs {
		for _, e := range pkg.Errors {
			err = multierr.Append(err, e)
		}
		if err != nil {
			return nil, err
		}
		ipkgs = append(ipkgs, internal.NewPackage(pkg))
	}
	return ipkgs, nil
}

// parseArchive parses the archive string to the internal.Archive type.
//
// The following is the flag format:
//
//  --archive=IMPORTPATHS=IMPORTMAP=FILE=EXPORT
//
// For example,
//
//  --archive=github.com/foo/bar:github.com/foo/baz=github.com/foo/bar=bar.go=bar_export.go
//
// The flag is structured in this format to closely follow https://github.com/bazelbuild/rules_go/blob/8ea79bbd5e6ea09dc611c245d1dc09ef7ab7118a/go/private/actions/compile.bzl#L20;
// however, the IMPORTPATHS and EXPORT elements are ignored. There may be future
// work involved in resolving import aliases, using IMPORTPATHS.
func parseArchive(archive string) (internal.Archive, error) {
	args := strings.Split(archive, "=")
	if len(args) != 4 {
		return internal.Archive{}, fmt.Errorf("expected 4 elements, got %d", len(args))
	}

	// Currently, we ignore the IMPORTPATHS and EXPORT elements.
	return internal.Archive{
		ImportMap: args[1],
		File:      args[2],
	}, nil
}
