package lib

// StringVal returns the value of a string pointer or empty string if pointer is nil
func StringVal(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

// Uint64Val returns the value of a uint64 pointer or empty uint64 if pointer is nil
func Uint64Val(n *uint64) uint64 {
	if n == nil {
		return 0
	}

	return *n
}
