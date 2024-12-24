package resources

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx/api"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx/mocks"
	"github.com/severalnines/terraform-provider-ccx/internal/lib"
	"github.com/stretchr/testify/mock"
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

		datastoreSvc, err := api.Datastores(httpClient, cfg.Timeout)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		contentSvc, err := api.Content(httpClient)
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

type mockServices struct {
	datastore *mocks.MockDatastoreService
	vpc       *mocks.MockVPCService
}

func (m mockServices) AssertExpectations(t mock.TestingT) {
	m.datastore.AssertExpectations(t)
	m.vpc.AssertExpectations(t)
}

func mockProvider() (mockServices, *schema.Provider) {
	datastore := &Datastore{}
	vpc := &VPC{}

	services := mockServices{
		datastore: &mocks.MockDatastoreService{},
		vpc:       &mocks.MockVPCService{},
	}

	configure := func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		datastore.svc = services.datastore
		vpc.svc = services.vpc

		return nil, nil
	}

	return services, provider(configure, datastore, vpc)
}
