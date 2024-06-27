package resources

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/severalnines/terraform-provider-ccx/ccx"
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

func getFirewalls(d *schema.ResourceData, key string) ([]ccx.FirewallRule, error) {
	var raw []interface{}

	if v, ok := d.Get(key).([]any); ok {
		raw = v
	}

	ls := make([]ccx.FirewallRule, len(raw))

	for i := range raw {
		if v, ok := raw[i].(map[string]any); ok {
			f, err := firewallFromMapAny(v)
			if err != nil {
				return nil, fmt.Errorf("invalid value for firewall[%d]: %w", i, err)
			}

			ls[i] = *f
		}
	}

	return ls, nil
}

func setFirewalls(d *schema.ResourceData, fieldname string, firewalls []ccx.FirewallRule) error {
	value := make([]map[string]any, 0, len(firewalls))

	for _, f := range firewalls {
		value = append(value, map[string]any{
			"source":      f.Source,
			"description": f.Description,
		})
	}

	return d.Set(fieldname, value)
}

func firewallFromMapAny(m map[string]any) (*ccx.FirewallRule, error) {
	var f ccx.FirewallRule

	if v, ok := m["description"].(string); ok {
		f.Description = v
	}

	if v, ok := m["source"].(string); ok {
		f.Source = v
	} else {
		return nil, fmt.Errorf("mandatory field source is missing")
	}

	return &f, nil
}

// nonNewSuppressor suppresses diff for fields when the resource is not new
func nonNewSuppressor(_, _, _ string, d *schema.ResourceData) bool {
	return !d.IsNewResource()
}
