package terraform

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform/helper/schema"
)

func ListToStrings(l types.List) []string {
	var s []string
	el := l.Elements()

	for _, i := range el {
		str, ok := i.(types.String)
		if !ok {
			continue
		}

		s = append(s, str.ValueString())
	}

	return s
}

func StringsToList(l []string) types.List {
	var s []attr.Value
	for _, i := range l {
		v := types.StringValue(i)
		s = append(s, v)
	}

	ls, _ := types.ListValue(types.StringType, s)
	return ls
}

func GetString(d *schema.ResourceData, key string) string {
	if v, ok := d.Get(key).(string); ok {
		return v
	}

	return ""
}

func GetInt(d *schema.ResourceData, key string) int64 {
	if v, ok := d.Get(key).(int64); ok {
		return v
	}

	return 0
}

func GetBool(d *schema.ResourceData, key string) bool {
	if v, ok := d.Get(key).(bool); ok {
		return v
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
