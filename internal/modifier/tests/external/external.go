package external

// An A is a test type.
type A int

// B is a test type.
type B bool

// Run is a test function.
func Run(a A) B {
	if a > 0 {
		return true
	}
	return false
}
