package resources

import (
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
)

type firewall struct{}

func (f firewall) Schema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CIDR source for the firewall rule",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of this firewall rule",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func parseRawFirewalls(raw []any) ([]ccx.FirewallRule, error) {
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

	slices.SortStableFunc(ls, func(a, b ccx.FirewallRule) int {
		return strings.Compare(a.Source, b.Source)
	})

	return ls, nil
}

func getFirewallsOldNew(d *schema.ResourceData) (old, nw []ccx.FirewallRule, err error) {
	oldValue, newValue := d.GetChange("firewall")
	if o, ok := oldValue.([]any); !ok {
		return nil, nil, fmt.Errorf("failed to read old firewalls")
	} else if old, err = parseRawFirewalls(o); err != nil {
		return nil, nil, fmt.Errorf("failed to parse old firewalls: %w", err)
	}

	if n, ok := newValue.([]any); !ok {
		return nil, nil, fmt.Errorf("failed to read new firewalls")
	} else if nw, err = parseRawFirewalls(n); err != nil {
		return nil, nil, fmt.Errorf("failed to parse new firewalls: %w", err)
	}

	return
}

func getFirewalls(d *schema.ResourceData) ([]ccx.FirewallRule, error) {
	v, ok := d.Get("firewall").([]any)
	if !ok {
		return nil, fmt.Errorf("failed to read firewalls")
	}

	return parseRawFirewalls(v)
}

func firewallsSame(ls1, ls2 []ccx.FirewallRule) bool {
	if len(ls1) != len(ls2) {
		return false
	}

	haveM := make(map[ccx.FirewallRule]struct{}, len(ls1))
	for _, f := range ls1 {
		haveM[f] = struct{}{}
	}

	for _, f := range ls2 {
		if _, ok := haveM[f]; !ok {
			return false
		}
	}

	return true
}

func setFirewalls(d *schema.ResourceData, firewalls []ccx.FirewallRule) error {
	value := make([]map[string]any, 0, len(firewalls))

	slices.SortStableFunc(firewalls, func(a, b ccx.FirewallRule) int {
		return strings.Compare(a.Source, b.Source)
	})

	for _, f := range firewalls {
		value = append(value, map[string]any{
			"id":          f.Source,
			"source":      f.Source,
			"description": f.Description,
		})
	}

	return d.Set("firewall", value)
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
