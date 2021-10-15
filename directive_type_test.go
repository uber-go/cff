package cff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirectiveTypeString(t *testing.T) {
	tests := []struct {
		desc string
		give DirectiveType
		want string
	}{
		{
			desc: "unknown",
			give: UnknownDirective,
		},
		{
			desc: "flow",
			give: FlowDirective,
		},
		{
			desc: "parallel",
			give: ParallelDirective,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.give.String(), tt.desc)
		})
	}
}
