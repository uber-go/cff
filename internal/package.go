package internal

import (
	"go/token"
)

// LoadParams is the arguments for func Load(...).
type LoadParams struct {
	Fset       *token.FileSet
	ImportPath string
}

