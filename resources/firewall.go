package resources

import (
	"github.com/hashicorp/terraform/helper/schema"
)

type Firewall struct{}

func (f Firewall) Schema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
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
			State: schema.ImportStatePassthrough,
		},
	}
}
