package terraform

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
)

func GetString(d *schema.ResourceData, key string) string {
	v, ok := d.GetOkExists(key)

	if !ok {
		return ""
	}

	switch s := v.(type) {
	case string:
		return s
	case *string:
		return *s
	}

	return ""
}

func GetInt(d *schema.ResourceData, key string) int64 {
	v, ok := d.GetOk(key)
	if !ok {
		return 0
	}

	switch i := v.(type) {
	case int64:
		return i
	case *int64:
		return *i
	case int:
		return int64(i)
	case *int:
		return int64(*i)
	case uint64:
		return int64(i)
	case *uint64:
		return int64(*i)
	case uint:
		return int64(i)
	case *uint:
		return int64(*i)
	}

	return 0
}

func GetBool(d *schema.ResourceData, key string) bool {
	v, ok := d.GetOk(key)
	if !ok {
		return false
	}

	switch b := v.(type) {
	case bool:
		return b
	case *bool:
		return *b
	}

	return false
}

// GetStrings from ResourceData
func GetStrings(d *schema.ResourceData, key string) []string {
	var raw []interface{}

	if v, ok := d.Get(key).([]any); ok {
		raw = v
	}

	s := make([]string, len(raw))

	for i := range raw {
		if v, ok := raw[i].(string); ok {
			s[i] = v
		}
	}

	return s
}

func ValidateName(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected name to be string"))
		return warns, errs
	}
	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(value)) {
		errs = append(errs, fmt.Errorf("name cannot contain whitespace. Got %s", value))
		return warns, errs
	}
	return warns, errs
}
