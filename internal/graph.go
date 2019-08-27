package internal

type graph struct {
	// Number of nodes in this graph.
	Count int

	// Indexes of the root nodes.
	Roots []int

	// Returns the dependencies of the node at the given index.
	Dependencies func(int) []int
}

func scheduleGraph(g graph) [][]int {
	// distances[i] is distance of node i from root.
	distances := make([]int, g.Count)

	nodes := g.Roots
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
