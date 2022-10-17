package flag

import (
	"flag"
	"strings"
)

// GetterPtr is satisfied by types
// that implement flag.Getter on their pointer receiver.
type GetterPtr[T any] interface {
	*T
	flag.Getter
}

// List is a list of values received over the CLI.
//
// Use this for flags that can be accepted multiple times.
type List[T any, P GetterPtr[T]] []T

// AsList builds a flag.List from a pointer to a slice
// of flag.Getter values.
//
//	var list []SomeObject
//	flagSet.Var(flag.AsList(&list), ...)
func AsList[T any, P GetterPtr[T]](vs *[]T) *List[T, P] {
	return (*List[T, P])(vs)
}

var _ flag.Getter = (*List[String, *String])(nil)

func (vl *List[T, P]) String() string {
	var sb strings.Builder
	for i, v := range *vl {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(P(&v).String())
	}
	return sb.String()
}

// Set receives a single value from the command line.
func (vl *List[T, P]) Set(v string) error {
	var t T
	if err := P(&t).Set(v); err != nil {
		return err
	}
	*vl = append(*vl, t)
	return nil
}

// Get returns a list of values in this object.
func (vl *List[T, P]) Get() any {
	result := make([]any, len(*vl))
	for i, v := range *vl {
		result[i] = P(&v).Get()
	}
	return result
}
