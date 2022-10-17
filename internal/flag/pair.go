package flag

import (
	"errors"
	"flag"
	"strings"
)

// InOutPair is a flag value that holds a pair of values.
// Two formats are supported:
//
//	--flag=INPUT
//	--flag=INPUT=OUTPUT
//
// For example,
//
//	--flag=foo.go=_gen/foo.go --flag=bar.go
type InOutPair struct {
	Input, Output string
}

var _ flag.Getter = (*InOutPair)(nil)

func (p *InOutPair) String() string {
	if len(p.Output) == 0 {
		return p.Input
	}
	return p.Input + "=" + p.Output
}

// Set receives a flag value.
func (p *InOutPair) Set(name string) error {
	var output string
	if i := strings.IndexByte(name, '='); i >= 0 {
		name, output = name[:i], name[i+1:]
	}

	if len(name) == 0 {
		return errors.New("input cannot be empty")
	}

	p.Input = name
	p.Output = output
	return nil
}

// Get reports the current value of the flag pair.
func (p *InOutPair) Get() any { return p }
