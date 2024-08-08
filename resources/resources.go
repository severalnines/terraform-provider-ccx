package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

type TerraformConfiguration struct {
	ClientID     string
	ClientSecret string
	BaseURL      string
	Timeout      time.Duration
}

type TerraformSchema interface {
	Schema() *schema.Resource
}

type TerraformResource interface {
	TerraformSchema
	Configure(ctx context.Context, cfg TerraformConfiguration) error

	Name() string
	Create(*schema.ResourceData, interface{}) error
	Read(*schema.ResourceData, interface{}) error
	Update(*schema.ResourceData, interface{}) error
	Delete(*schema.ResourceData, interface{}) error
	// Exists(*schema.ResourceData, interface{}) (bool, error)
}

func Provider(r ...TerraformResource) *provider {
	return &provider{resources: r}
}

type provider struct {
	resources []TerraformResource
	Config    TerraformConfiguration
}

func (p *provider) Resources() terraform.ResourceProvider {
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
			"timeout": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CCX_TIMEOUT", "15m"),
			},
		},
		ResourcesMap:  rsc,
		ConfigureFunc: p.Configure,
	}
}

func (p *provider) Configure(d *schema.ResourceData) (any, error) {
	ctx := context.Background()

	p.Config = TerraformConfiguration{
		ClientID:     getString(d, "client_id"),
		ClientSecret: getString(d, "client_secret"),
		BaseURL:      getString(d, "base_url"),
	}

	if t, err := time.ParseDuration(getString(d, "timeout")); err == nil {
		p.Config.Timeout = t
	} else {
		return nil, fmt.Errorf("invalid timeout (%s): %w", getString(d, "timeout"), err)
	}

	for i := range p.resources {
		if err := p.resources[i].Configure(ctx, p.Config); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
