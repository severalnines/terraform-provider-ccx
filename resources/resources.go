package resources

import (
	"context"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx/api"
	"github.com/severalnines/terraform-provider-ccx/internal/lib"
)

type TerraformConfiguration struct {
	ClientID     string
	ClientSecret string
	BaseURL      string
	Timeout      time.Duration
	Logpath      string

	httpClient api.HttpClient
}

type TerraformSchema interface {
	Schema() *schema.Resource
}

type TerraformResource interface {
	TerraformSchema
	Configure(cfg TerraformConfiguration) error

	Name() string
	Create(ctx context.Context, d *schema.ResourceData, _ any) diag.Diagnostics
	Read(ctx context.Context, d *schema.ResourceData, _ any) diag.Diagnostics
	Update(ctx context.Context, d *schema.ResourceData, _ any) diag.Diagnostics
	Delete(ctx context.Context, d *schema.ResourceData, _ any) diag.Diagnostics
	// Exists(*schema.ResourceData, interface{}) (bool, error)
}

func Provider(r ...TerraformResource) *provider {
	return &provider{resources: r}
}

type provider struct {
	resources []TerraformResource
	Config    TerraformConfiguration
}

func (p *provider) Resources() *schema.Provider {
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
		ResourcesMap:         rsc,
		ConfigureContextFunc: p.Configure,
	}
}

func (p *provider) Configure(_ context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
	p.Config = TerraformConfiguration{
		ClientID:     getString(d, "client_id"),
		ClientSecret: getString(d, "client_secret"),
		BaseURL:      getString(d, "base_url"),
		Logpath:      getString(d, "debug_log_path"),
	}

	if t, err := time.ParseDuration(getString(d, "timeout")); err == nil {
		p.Config.Timeout = t
	} else {
		return nil, diag.Errorf("invalid timeout (%s): %s", getString(d, "timeout"), err)
	}

	if p.Config.Logpath != "" {
		if err := os.MkdirAll(p.Config.Logpath, 0755); err != nil {
			return nil, diag.Errorf("creating log directory [%s]: %s", p.Config.Logpath, err)
		}
	}

	p.Config.httpClient = lib.NewHttpClient(p.Config.BaseURL, p.Config.ClientID, p.Config.ClientSecret)

	for i := range p.resources {
		if err := p.resources[i].Configure(p.Config); err != nil {
			return nil, diag.FromErr(err)
		}
	}

	return nil, nil
}

// nonNewSuppressor suppresses diff for fields when the resource is not new
func nonNewSuppressor(_, _, _ string, d *schema.ResourceData) bool {
	return !d.IsNewResource()
}
