package terraform

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ccxprov "github.com/severalnines/terraform-provider-ccx"
)

var _ provider.Provider = &Provider{}

func New(r ...resource.Resource) *Provider {
	return &Provider{resources: r}
}

type Provider struct {
	resources []resource.Resource
	Config    *Configuration
}

func (p *Provider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.Version = "1.1.0"
	resp.TypeName = "CCX Provider"
}

func (p *Provider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	// var providers Providers

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Optional: true,
			},
			"client_secret": schema.StringAttribute{
				Optional: true,
			},
			"base_url": schema.StringAttribute{
				Optional: true,
			},
			"dev_mode": schema.BoolAttribute{
				Optional: true,
			},
			"mock_file": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

type ConfigurationModel struct {
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	BaseURL      types.String `tfsdk:"base_url"`
	IsDevMode    types.Bool   `tfsdk:"dev_mode"`
	Mockfile     types.String `tfsdk:"mock_file"`
}

type Configuration struct {
	ClientID     string
	ClientSecret string
	BaseURL      string
	IsDevMode    bool
	Mockfile     string
}

func FromModel(m ConfigurationModel) (*Configuration, error) {
	clientID, err := getValueOrEnv(m.ClientID.ValueString(), "CCX_CLIENT_ID")
	if err != nil {
		return nil, err
	}

	clientSecret, err := getValueOrEnv(m.ClientSecret.ValueString(), "CCX_CLIENT_SECRET")
	if err != nil {
		return nil, err
	}

	return &Configuration{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		BaseURL:      m.BaseURL.ValueString(),
		IsDevMode:    m.IsDevMode.ValueBool(),
		Mockfile:     m.Mockfile.ValueString(),
	}, nil
}

type MockData struct {
	Clusters map[string]ccxprov.Cluster `json:"clusters"`
	VPCS     map[string]ccxprov.VPC     `json:"VPCS"`
}

func (p *Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var (
		cfg ConfigurationModel
		err error
	)
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}

	p.Config, err = FromModel(cfg)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read provider configuration", err.Error())
		return
	}

	for i := range p.resources {
		r, ok := p.resources[i].(Configurable)
		if !ok {
			continue
		}

		if err := r.Configure(ctx, p, resp.Diagnostics); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Failed to configure resource [%d]", i), err.Error())
			return
		}
	}
}

func (p *Provider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (p *Provider) Resources(_ context.Context) []func() resource.Resource {
	var r []func() resource.Resource

	for i := range p.resources {
		r = append(r, resourceFunc(p.resources[i]))
	}

	return r
}

func getValueOrEnv(value string, env string) (string, error) {
	if value == "" {
		value = os.Getenv(env)
	}

	if value == "" {
		return "", ccxprov.MissingParameterErr
	}

	return value, nil
}

type Configurable interface {
	Configure(ctx context.Context, p *Provider, diagnostics diag.Diagnostics) error
}

func resourceFunc(r resource.Resource) func() resource.Resource {
	return func() resource.Resource {
		return r
	}
}
