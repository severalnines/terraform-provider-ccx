package resources

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func getString(d *schema.ResourceData, key string) string {
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

func getInt(d *schema.ResourceData, key string) int64 {
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

func getBool(d *schema.ResourceData, key string) bool {
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

// getStrings from ResourceData
func getStrings(d *schema.ResourceData, key string) []string {
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

func getMapString(d *schema.ResourceData, key string) map[string]string {
	raw, ok := d.GetOk(key)
	if !ok {
		return nil
	}

	m := make(map[string]string)
	for k, v := range raw.(map[string]interface{}) {
		s, ok := v.(string)
		if ok {
			m[k] = s
		}
	}

	return m
}

func allKeysSet(d *schema.ResourceData, keys ...string) bool {
	for _, key := range keys {
		if _, ok := d.GetOk(key); !ok {
			return false
		}
	}

	return true
}
