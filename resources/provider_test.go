package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx/mocks"
	"github.com/stretchr/testify/mock"
)

type mockServices struct {
	datastore      *mocks.MockDatastoreService
	vpc            *mocks.MockVPCService
	parameterGroup *mocks.MockParameterGroupService
}

func (m mockServices) AssertExpectations(t mock.TestingT) {
	m.datastore.AssertExpectations(t)
	m.vpc.AssertExpectations(t)
	m.parameterGroup.AssertExpectations(t)
}

func mockProvider() (mockServices, *schema.Provider) {
	datastore := &Datastore{}
	vpc := &VPC{}
	parameterGroup := &ParameterGroup{}

	services := mockServices{
		datastore:      &mocks.MockDatastoreService{},
		vpc:            &mocks.MockVPCService{},
		parameterGroup: &mocks.MockParameterGroupService{},
	}

	configure := func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		datastore.svc = services.datastore
		vpc.svc = services.vpc
		parameterGroup.svc = services.parameterGroup

		return nil, nil
	}

	return services, provider(configure, datastore, vpc, parameterGroup)
}
