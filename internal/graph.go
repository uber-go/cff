package internal

type graph struct {
	// Number of nodes in this graph.
	Count int

	// Returns the dependencies of the node at the given index.
	Dependencies func(int) []int
}

// Returns a topologically sorted list of nodes.
//
// Does not detect cycles. Make sure that is already done.
func toposort(g graph) []int {
	topo := make([]int, 0, g.Count)

	visited := make(map[int]struct{}, g.Count)
	var visit func(int)
	visit = func(n int) {
		if _, ok := visited[n]; ok {
			return
		}

		for _, d := range g.Dependencies(n) {
			visit(d)
		}

		visited[n] = struct{}{}
		topo = append(topo, n)
	}

	for n := 0; n < g.Count; n++ {
		visit(n)
	}

	return topo
}
