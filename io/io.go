package io

type Closable interface {
	Close() error
}

// Close will close a Closable object silently
// useful for deferred closing
func Close(c Closable) {
	if c != nil {
		_ = c.Close()
	}
}
