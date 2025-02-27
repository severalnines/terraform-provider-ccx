package resources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx/mocks"
	"github.com/stretchr/testify/mock"
)

type mockServices struct {
	datastore      *mocks.MockDatastoreService
	vpc            *mocks.MockVPCService
	parameterGroup *mocks.MockParameterGroupService
	content        *mocks.MockContentService
}

func (m mockServices) AssertExpectations(t mock.TestingT) {
	m.datastore.AssertExpectations(t)
	m.vpc.AssertExpectations(t)
	m.parameterGroup.AssertExpectations(t)
}

func mockProvider(t *testing.T) (mockServices, *schema.Provider) {
	datastore := &Datastore{}
	vpc := &VPC{}
	parameterGroup := &ParameterGroup{}

	services := mockServices{
		datastore:      mocks.NewMockDatastoreService(t),
		vpc:            mocks.NewMockVPCService(t),
		parameterGroup: mocks.NewMockParameterGroupService(t),
		content:        mocks.NewMockContentService(t),
	}

	configure := func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		vpc.svc = services.vpc
		parameterGroup.svc = services.parameterGroup
		parameterGroup.contentSvc = services.content
		datastore.svc = services.datastore
		datastore.contentSvc = services.content
		datastore.pgSvc = services.parameterGroup

		return nil, nil
	}

	return services, provider(configure, datastore, vpc, parameterGroup)
}
