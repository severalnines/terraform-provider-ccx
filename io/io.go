package io

import (
	"encoding/json"
	"os"

	ccxprov "github.com/severalnines/terraform-provider-ccx"
)

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

// LoadData from disk
func LoadData(path string, target any) error {
	if path == "" {
		return ccxprov.MockPathEmptyErr
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, target)
	if err != nil {
		return err
	}

	return nil
}
