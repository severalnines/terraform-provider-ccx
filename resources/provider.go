package resources

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
)

type providerConfig struct {
	ClientID     string
	ClientSecret string
	BaseURL      string
	Timeout      time.Duration
}

func Provider() *schema.Provider {
	// make resource managers, so they are ready to be used in schema, but we can't set services into them until configure is called
	datastore := &Datastore{}
	vpc := &VPC{}
	parameterGroup := &ParameterGroup{}

	configure := func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		cfg := providerConfig{
			ClientID:     getString(d, "client_id"),
			ClientSecret: getString(d, "client_secret"),
			BaseURL:      strings.Trim(getString(d, "base_url"), "/"),
		}

		if t, err := time.ParseDuration(getString(d, "timeout")); err == nil {
			cfg.Timeout = t
		} else {
			return nil, diag.Errorf("invalid timeout (%s): %s", getString(d, "timeout"), err)
		}

		httpClient := ccx.NewHTTPClient(cfg.BaseURL, cfg.ClientID, cfg.ClientSecret)

		contentSvc, err := ccx.NewContentClient(httpClient)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		datastoreSvc, err := ccx.NewDatastoresClient(httpClient, cfg.Timeout, contentSvc)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		parameterGroupSvc := ccx.NewParameterGroupsClient(httpClient)

		vpcSvc := ccx.NewVPCsClient(httpClient)

		// set services into resources, now that it is possible

		datastore.svc = datastoreSvc
		datastore.contentSvc = contentSvc
		datastore.pgSvc = parameterGroupSvc

		parameterGroup.svc = parameterGroupSvc
		parameterGroup.contentSvc = contentSvc

		vpc.svc = vpcSvc

		return nil, nil
	}

	return makeProvider(configure, datastore, vpc, parameterGroup)
}

func makeProvider(configure schema.ConfigureContextFunc, datastore *Datastore, vpc *VPC, parameterGroup *ParameterGroup) *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CCX_CLIENT_ID", ""),
				Description: "OAuth client ID, which can be created in the CCX UI.",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CCX_CLIENT_SECRET", ""),
				Description: "OAuth client secret, which can be created in the CCX UI.",
			},
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CCX_BASE_URL", "https://app.mydbservice.net"),
				Description: "If you are using a CCX instance other than the public service provided by Severalnines, set this value. It should be as a URL, e.g. `https://ccx.mycloud.com`.",
			},
			"timeout": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CCX_TIMEOUT", "60m"),
				Description: "Optionally, set a timeout for something. The default is `60m` meaning 60 minutes.",
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
