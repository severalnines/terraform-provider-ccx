package terraform

import (
	"context"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/severalnines/terraform-provider-ccx/ccx"
)

func ToConfiguration(d *schema.ResourceData) ccx.TerraformConfiguration {
	return ccx.TerraformConfiguration{
		ClientID:     GetString(d, "client_id"),
		ClientSecret: GetString(d, "client_secret"),
		BaseURL:      GetString(d, "base_url"),
	}
}

func New(r ...ccx.TerraformResource) *Provider {
	return &Provider{resources: r}
}

type Provider struct {
	resources []ccx.TerraformResource
	Config    ccx.TerraformConfiguration
}

func (p *Provider) Resources() terraform.ResourceProvider {
	rsc := map[string]*schema.Resource{}
	for i := range p.resources {
		name := p.resources[i].Name()
		rsc[name] = p.resources[i].Schema()
	}

	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CCX_CLIENT_ID", ""),
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CCX_CLIENT_SECRET", ""),
			},
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CCX_BASE_URL", "https://app.mydbservice.net"),
			},
		},
		ResourcesMap:  rsc,
		ConfigureFunc: p.Configure,
	}
}

func (p *Provider) Configure(d *schema.ResourceData) (interface{}, error) {
	p.Config = ToConfiguration(d)

	for i := range p.resources {
		if err := p.resources[i].Configure(context.Background(), p.Config); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
