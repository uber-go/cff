package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToposort(t *testing.T) {
	tests := []struct {
		desc string

		// Adjacency list definining the graph.
		//
		// give[i] lists dependencies of node i in the graph.
		give [][]int

		want []int
	}{
		{
			desc: "empty graph",
			want: []int{},
		},
		{
			desc: "single node",
			give: [][]int{
				{}, // 0
			},
			want: []int{0},
		},
		{
			desc: "two nodes/no edges",
			give: [][]int{
				{}, // 0
				{}, // 1
			},
			want: []int{0, 1},
		},
		{
			desc: "two nodes/one edge",
			give: [][]int{
				{1}, // 0
				{},  // 1
			},
			want: []int{1, 0},
		},
		{
			desc: "three nodes/one edge",
			give: [][]int{
				{},  // 0
				{0}, // 1
				{},  // 2
			},
			want: []int{0, 1, 2},
		},
		{
			desc: "three nodes/two edges from the same node",
			give: [][]int{
				{},     // 0
				{0, 2}, // 1
				{},     // 2
			},
			want: []int{0, 2, 1},
		},
		{
			desc: "three nodes/two edges to the same node",
			give: [][]int{
				{1}, // 0
				{},  // 1
				{1}, // 2
			},
			want: []int{1, 0, 2},
		},
		{
			desc: "three nodes/linked list",
			give: [][]int{
				{},  // 0
				{0}, // 1
				{1}, // 2
			},
			want: []int{0, 1, 2},
		},
		{
			desc: "complex case",
			give: [][]int{
				{3},       // 0
				{3, 4},    // 1
				{4, 7},    // 2
				{5, 6, 7}, // 3
				{6},       // 4
				{},        // 5
				{},        // 6
				{},        // 7
			},
			want: []int{5, 6, 7, 3, 0, 4, 1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			// sanity check: None of the nodes can refer to a node
			// that doesn't exist. The test is otherwise invalid.
			for i, deps := range tt.give {
				for _, d := range deps {
					if d >= len(tt.give) {
						t.Fatalf("node %d depends on undefined node %d", i, d)
					}
				}
			}

			g := graph{
				Count: len(tt.give),
				Dependencies: func(i int) []int {
					return tt.give[i]
				},
			}

			got := toposort(g)
			assert.Equal(t, tt.want, got)

			// sanity check: verify that walking in topo-order, we
			// never depend on an item before it has already been
			// seen.
			seen := make(map[int]struct{})
			for _, from := range got {
				seen[from] = struct{}{}
				for _, to := range tt.give[from] {
					if _, ok := seen[to]; !ok {
						t.Errorf("node %d depends on %d before it was seen", from, to)
					}
				}
			}
		})
	}
}
