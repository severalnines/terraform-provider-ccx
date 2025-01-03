package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx/mocks"
	"github.com/stretchr/testify/mock"
)

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
