package terraform

import (
	"context"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	ccxprov "github.com/severalnines/terraform-provider-ccx"
)

func ToConfiguration(d *schema.ResourceData) ccxprov.TerraformConfiguration {
	return ccxprov.TerraformConfiguration{
		ClientID:     GetString(d, "client_id"),
		ClientSecret: GetString(d, "client_secret"),
		BaseURL:      GetString(d, "base_url"),
		IsDevMode:    false,
		Mockfile:     "",
	}
}

func New(r ...ccxprov.TerraformResource) *Provider {
	return &Provider{resources: r}
}

type Provider struct {
	resources []ccxprov.TerraformResource
	Config    ccxprov.TerraformConfiguration
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
				DefaultFunc: schema.EnvDefaultFunc("CCX_BASE_URL", ""),
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
			return nil, nil
		}
	}

	return nil, nil
}
