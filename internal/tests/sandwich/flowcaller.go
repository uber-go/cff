package sandwich

// CallFlow exports a function that calls a CFF flow.
func CallFlow() (string, string) {
	s, _ := aFlow()
	t, _ := bFlow()

	return s, t
}
