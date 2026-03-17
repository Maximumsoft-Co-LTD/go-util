package slice

// Index iterates from 0 to limit-1 and returns the first index
// where predicate returns true. Returns -1 if no match is found.
func Index(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}
