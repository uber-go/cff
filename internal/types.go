package internal

import "go/types"

func isContext(t types.Type) bool {
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}

	o := named.Obj()
	return o.Pkg().Path() == "context" && o.Name() == "Context"
}

func isError(t types.Type) bool {
	n, ok := t.(*types.Named)
	if !ok {
		return false
	}
	o := n.Obj()
	return o.Pkg() == nil && o.Name() == "error"
}
