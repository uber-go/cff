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

// Builds a static execution schedule of the nodes defined in the graph
// starting at the provided roots.
func scheduleGraph(nodes []int, g graph) [][]int {
	// distances[i] is distance of node i from root.
	distances := make([]int, g.Count)

	d := 0 // current distance
	for len(nodes) > 0 {
		var newNodes []int

		// We don't need to track visited nodes. If we see a node we've
		// already visited, its distance from root is further than previously
		// thought. For example, suppose root depends on A and B, and A also
		// depends on B. The distance of B from root is not 1 but 2.

		for _, n := range nodes {
			distances[n] = d
			newNodes = append(newNodes, g.Dependencies(n)...)
		}

		d++
		nodes = newNodes
	}

	// Number of concurrently executable groups matches the maximum distance
	// from root.
	sched := make([][]int, d)
	for i, d := range distances {
		// All nodes with the same distance from root belong together. We'll
		// put them into the schedule in the reverse order so that sched[0] is
		// the set of tasks with no dependencies (maximum distance from root).
		s := len(sched) - 1 - d
		if s > -1 {
			sched[s] = append(sched[s], i)
		}
	}
	return sched
}
