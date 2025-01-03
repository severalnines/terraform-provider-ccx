package resources

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx/api"
	"github.com/severalnines/terraform-provider-ccx/internal/lib"
)

func Provider() *schema.Provider {
	datastore := &Datastore{}
	vpc := &VPC{}

	configure := func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		cfg := TerraformConfiguration{
			ClientID:     getString(d, "client_id"),
			ClientSecret: getString(d, "client_secret"),
			BaseURL:      getString(d, "base_url"),
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

		datastoreSvc, err := api.Datastores(httpClient, cfg.Timeout, contentSvc)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		datastore.svc = datastoreSvc
		datastore.contentSvc = contentSvc

		vpcSvc := api.Vpcs(httpClient)
		vpc.svc = vpcSvc

		return nil, nil
	}

	return provider(configure, datastore, vpc)
}

func provider(configure schema.ConfigureContextFunc, datastore *Datastore, vpc *VPC) *schema.Provider {
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
			"ccx_datastore": datastore.Schema(),
			"ccx_vpc":       vpc.Schema(),
		},
		ConfigureContextFunc: configure,
	}
}