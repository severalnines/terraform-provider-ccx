package resources

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx/api"
	"github.com/severalnines/terraform-provider-ccx/internal/lib"
)

func Provider() *schema.Provider {
	datastore := &Datastore{}
	vpc := &VPC{}
	parameterGroup := &ParameterGroup{}

	configure := func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		cfg := TerraformConfiguration{
			ClientID:     getString(d, "client_id"),
			ClientSecret: getString(d, "client_secret"),
			BaseURL:      strings.Trim(getString(d, "base_url"), "/"),
		}

		if t, err := time.ParseDuration(getString(d, "timeout")); err == nil {
			cfg.Timeout = t
		} else {
			return nil, diag.Errorf("invalid timeout (%s): %s", getString(d, "timeout"), err)
		}

		httpClient := lib.NewHttpClient(cfg.BaseURL, cfg.ClientID, cfg.ClientSecret)

		contentSvc, err := api.Content(httpClient)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		vpcSvc := api.Vpcs(httpClient)
		vpc.svc = vpcSvc

		parameterGroupSvc := api.ParameterGroups(httpClient)
		parameterGroup.svc = parameterGroupSvc
		parameterGroup.contentSvc = contentSvc

		datastoreSvc, err := api.Datastores(httpClient, cfg.Timeout, contentSvc)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		datastore.svc = datastoreSvc
		datastore.contentSvc = contentSvc
		datastore.pgSvc = parameterGroupSvc

		return nil, nil
	}

	return provider(configure, datastore, vpc, parameterGroup)
}

func provider(configure schema.ConfigureContextFunc, datastore *Datastore, vpc *VPC, parameterGroup *ParameterGroup) *schema.Provider {
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
				DefaultFunc: schema.EnvDefaultFunc("CCX_TIMEOUT", "60m"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ccx_datastore":       datastore.Schema(),
			"ccx_vpc":             vpc.Schema(),
			"ccx_parameter_group": parameterGroup.Schema(),
		},
		ConfigureContextFunc: configure,
	}
}
