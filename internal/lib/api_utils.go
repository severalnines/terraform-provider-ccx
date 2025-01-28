package lib

// StringVal returns the value of a string pointer or empty string if pointer is nil
func StringVal(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

// StringP returns a pointer to a string
func StringP(s string) *string {
	return &s
}

// Uint64Val returns the value of an uint64 pointer or empty uint64 if pointer is nil
func Uint64Val(n *uint64) uint64 {
	if n == nil {
		return 0
	}

	return *n
}

// Uint64P returns a pointer to an uint64
func Uint64P(n uint64) *uint64 {
	return &n
}
