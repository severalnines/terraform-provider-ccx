package resources

import (
	"fmt"

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

func getFirewalls(d *schema.ResourceData) ([]ccx.FirewallRule, error) {
	var raw []interface{}

	if v, ok := d.Get("firewall").([]any); ok {
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

func setFirewalls(d *schema.ResourceData, firewalls []ccx.FirewallRule) error {
	value := make([]map[string]any, 0, len(firewalls))

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
