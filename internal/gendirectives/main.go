// gendirectives generates a file with the following API:
//
//   package internal
//
//   func IsCodegenDirective(name string) bool
//
// The function returns true if the function with the provided name should be
// considered a code generation directive by CFF.
//
// Code generation directives are defined in cff.go in the root of the CFF
// project.
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"

	"text/template"
)

func main() {
	log.SetFlags(0)
	if err := run(os.Args[1:]); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: %v path_to_cff.go output_file.go", os.Args[0])
	}

	input := args[0]
	output := args[1]

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, input, nil, 0)
	if err != nil {
		return fmt.Errorf("parse %q: %w", input, err)
	}

	var td templateData
	for _, d := range f.Decls {
		decl, ok := d.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// We don't have method directives.
		if decl.Recv != nil {
			continue
		}

		td.Directives = append(td.Directives, decl.Name.Name)
	}

	out, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("create %q: %w", output, err)
	}
	defer out.Close()

	return _tmpl.Execute(out, td)
}

type templateData struct {
	Directives []string
}

var _tmpl = template.Must(template.New("directives.go").Parse(`
package internal

// List of functions in the CFF package that are code generation directives.
var _codegenDirectives = map[string]struct{}{
	{{ range .Directives -}}
		{{ printf "%q" . }}: {},
	{{ end -}}
}

// IsCodegenDirective reports whether the function with the given name in the
// CFF package is a code generation directive.
func IsCodegenDirective(name string) bool {
	_, ok := _codegenDirectives[name]
	return ok
}
`))
