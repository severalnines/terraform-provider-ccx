package resources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
	"github.com/stretchr/testify/mock"
)

type mockServices struct {
	datastore      *ccx.MockDatastoresService
	vpc            *ccx.MockVPCsService
	parameterGroup *ccx.MockParameterGroupsService
	content        *ccx.MockContentService
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
		datastore:      ccx.NewMockDatastoresService(t),
		vpc:            ccx.NewMockVPCsService(t),
		parameterGroup: ccx.NewMockParameterGroupsService(t),
		content:        ccx.NewMockContentService(t),
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

	return services, makeProvider(configure, datastore, vpc, parameterGroup)
}
