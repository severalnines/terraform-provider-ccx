package strings

// Sames checks if 2 slices of strings are same
func Sames(old, n []string) bool {
	if len(old) != len(n) {
		return false
	}

	o := map[string]any{}
	for i := range old {
		o[old[i]] = nil
	}

	for i := range n {
		_, ok := o[n[i]]
		if !ok {
			return false
		}
	}

	return true
}
