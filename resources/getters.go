package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getString(d *schema.ResourceData, key string) string {
	v := d.Get(key)

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
	v := d.Get(key)

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
	var raw []any

	val := d.Get(key)

	switch v := val.(type) {
	case *schema.Set:
		raw = v.List()
	case []any:
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

func isSubsetOf[T comparable](a, b []T) bool {
	m := make(map[T]any, len(b))
	for _, v := range b {
		m[v] = nil
	}

	for _, v := range a {
		if _, ok := m[v]; !ok {
			return false
		}
	}

	return true
}

// setTags
// the server adds additional tags
// terraform marks the field as dirty
// this is a hack to avoid that
// if the tags defined in terraform are a subset of the tags returned by the server,
// we keep the tags as is
func setTags(d *schema.ResourceData, key string, tags []string) error {
	old := getStrings(d, key)

	ok := isSubsetOf(old, tags)
	if ok {
		return setStrings(d, key, old)
	}

	return setStrings(d, key, tags)
}

func setStrings(d *schema.ResourceData, key string, values []string) error {
	var raw []any

	for _, v := range values {
		raw = append(raw, v)
	}

	return d.Set(key, raw)
}

func allKeysSet(d *schema.ResourceData, keys ...string) bool {
	for _, key := range keys {
		if _, ok := d.GetOk(key); !ok {
			return false
		}
	}

	return true
}

func getAzs(d *schema.ResourceData) ([]string, bool) {
	if _, ok := d.GetOk("network_az"); ok {
		azs := getStrings(d, "network_az")
		return azs, true
	}

	return nil, false
}
