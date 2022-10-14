package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"go/build/constraint"
	"io"
)

// invertCFFConstraint takes the address of a parsed
// build tag (// +build) or constraint (//go:build)
// and modifies it in-place to invert all instances of the "cff" tag.
func invertCFFConstraint(exp *constraint.Expr) {
	switch ex := (*exp).(type) {
	case *constraint.AndExpr:
		invertCFFConstraint(&ex.X)
		invertCFFConstraint(&ex.Y)
	case *constraint.NotExpr:
		// Special-case: If "X" in "!X" is "cff",
		// just remove the "!".
		if t, ok := ex.X.(*constraint.TagExpr); ok && t.Tag == "cff" {
			*exp = ex.X
			return
		}
		invertCFFConstraint(&ex.X)
	case *constraint.OrExpr:
		invertCFFConstraint(&ex.X)
		invertCFFConstraint(&ex.Y)
	case *constraint.TagExpr:
		if ex.Tag == "cff" {
			*exp = &constraint.NotExpr{X: ex}
		}
	}
}

var _newline = []byte("\n")

// writeInvertedCFFTag writes the provided byte slice to the io.Writer,
// accounting for any "cff" build constraints in the byte slice
// by inverting them.
func writeInvertedCFFTag(w io.Writer, bs []byte) error {
	// This reduces the unnecessary if err != nil statements below.
	errw := stickyErrWriter{W: w}
	w = &errw

	// For each line before the package clause,
	// if it's a build constraint that contains "cff", invert "cff".
	scan := bufio.NewScanner(bytes.NewReader(bs))
	for scan.Scan() {
		line := scan.Text()
		isGoBuild := constraint.IsGoBuild(line)
		isPlusBuild := constraint.IsPlusBuild(line)
		if !isGoBuild && !isPlusBuild {
			fmt.Fprintln(w, line)
			continue
		}

		expr, err := constraint.Parse(line)
		if err != nil {
			// Leave invalid constraints unchanged.
			fmt.Fprintln(w, line)
			continue
		}
		invertCFFConstraint(&expr)

		if isGoBuild {
			fmt.Fprintf(w, "//go:build %v\n", expr.String())
			continue
		}

		lines, err := constraint.PlusBuildLines(expr)
		if err != nil {
			// This is not possible in production cases
			// because if we could parse it from build tags,
			// we can turn it back into build tags.
			// So we won't really see coverage here.
			fmt.Fprintln(w, line)
			continue
		}
		for _, tagline := range lines {
			fmt.Fprintln(w, tagline)
		}
	}
	return errw.Err
}

// stickyErrWriter is a Writer that never returns an error,
// but records it internally.
// We use this to reduce boilerplate while writing.
type stickyErrWriter struct {
	W   io.Writer
	Err error
}

func (w *stickyErrWriter) Write(bs []byte) (int, error) {
	// Already failed. Pretend this worked.
	if w.Err != nil {
		return len(bs), nil
	}

	n, err := w.W.Write(bs)
	if err != nil {
		w.Err = err
		return len(bs), nil
	}

	return n, nil
}
