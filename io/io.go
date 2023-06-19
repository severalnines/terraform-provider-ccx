package io

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

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

// Dump will dump a value in a text file. Useful for debugging. The name is used as prefix to the file.
// An attempt will be made to encode the value to json, if not possible, a golang string will be output.
func Dump(name string, value any) {
	filename := name + "_" + time.Now().Format("2006_01_02-15_04_05_999999999Z07_00") + ".txt"
	var b []byte

	if j, err := json.MarshalIndent(value, "", "    "); err == nil {
		b = j
	} else {
		b = []byte(fmt.Sprintf("%T\n%+v", value, value))
	}

	_ = os.WriteFile(filename, b, 0644)
}
