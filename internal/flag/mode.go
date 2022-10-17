package flag

import (
	"encoding"
	"flag"
	"fmt"
)

// Mode specifies the code generation mode for CFF.
type Mode uint8

const (
	// BaseMode generates CFF code without modification.
	BaseMode Mode = iota + 1

	// SourceMapMode generates CFF code with line directives to remap
	// generated code locations to source.
	SourceMapMode

	// ModifierMode generates CFF code that preserves all original file line
	// locations by generating CFF logic into separate modifier
	// functions.
	ModifierMode
)

var (
	_ encoding.TextUnmarshaler = (*Mode)(nil)
	_ flag.Getter              = (*Mode)(nil)
)

func (m Mode) String() string {
	switch m {
	case SourceMapMode:
		return "source-map"
	case BaseMode:
		return "base"
	case ModifierMode:
		return "modifier"
	default:
		return "unknown"
	}
}

// UnmarshalText unmarshals a Mode.
func (m *Mode) UnmarshalText(text []byte) error {
	return m.Set(string(text))
}

// Get reports the current value of the flag.
func (m *Mode) Get() any {
	return *m
}

// Set receives a flag value from the flag package.
func (m *Mode) Set(value string) error {
	switch value {
	case "source-map":
		*m = SourceMapMode
	case "base":
		*m = BaseMode
	case "modifier":
		*m = ModifierMode
	default:
		return fmt.Errorf("unknown mode %q", value)
	}
	return nil
}
