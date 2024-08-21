package resources

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

type TerraformConfiguration struct {
	ClientID     string
	ClientSecret string
	BaseURL      string
	Timeout      time.Duration
	Logpath      string
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
			"debug_log_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CCX_DEBUG_LOG_PATH", ""),
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
		Logpath:      getString(d, "debug_log_path"),
	}

	if t, err := time.ParseDuration(getString(d, "timeout")); err == nil {
		p.Config.Timeout = t
	} else {
		return nil, fmt.Errorf("invalid timeout (%s): %w", getString(d, "timeout"), err)
	}

	if p.Config.Logpath != "" {
		if err := os.MkdirAll(p.Config.Logpath, 0755); err != nil {
			return nil, fmt.Errorf("creating log directory [%s]: %w", p.Config.Logpath, err)
		}
	}

	for i := range p.resources {
		if err := p.resources[i].Configure(ctx, p.Config); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

// nonNewSuppressor suppresses diff for fields when the resource is not new
func nonNewSuppressor(_, _, _ string, d *schema.ResourceData) bool {
	return !d.IsNewResource()
}
