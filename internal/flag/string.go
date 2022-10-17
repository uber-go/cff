package flag

import "flag"

// String is a string received as a command line argument.
//
// We define our own to use it with ListValue.
type String string

var _ flag.Getter = (*String)(nil)

func (s *String) String() string { return string(*s) }

// Get returns the current string value.
func (s *String) Get() any { return s.String() }

// Set receives a command line value.
func (s *String) Set(v string) error {
	*s = String(v)
	return nil
}
