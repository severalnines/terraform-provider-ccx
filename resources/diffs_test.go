package resources

import (
	"context"
	"errors"
	"testing"

	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_checkInstanceSizeEquivalence(t *testing.T) {
	instanceSizes := map[string][]ccx.InstanceSize{
		"aws": {
			{
				Code: "Tiny",
				Type: "t2.micro",
			},
			{
				Code: "Small",
				Type: "t2.small",
			},
			{
				Code: "Medium",
				Type: "t2.medium",
			},
		},
		"openstack": {
			{
				Code: "Alpha",
				Type: "t.alpha",
			},
			{
				Code: "Beta",
				Type: "t.beta",
			},
			{
				Code: "Gamma",
				Type: "t.gamma",
			},
		},
	}

	tests := []struct {
		name          string
		cloudProvider string
		oldValue      string
		newValue      string
		mock          func(svc *mocks.MockContentService)
		want          bool
	}{
		{
			name:          "old and new values are exactly the same",
			cloudProvider: "aws",
			oldValue:      "Tiny",
			newValue:      "Tiny",
			want:          true,
		},
		{
			name:          "old and new values are equivalent",
			cloudProvider: "aws",
			oldValue:      "t2.micro",
			newValue:      "Tiny",
			mock: func(svc *mocks.MockContentService) {
				svc.EXPECT().InstanceSizes(mock.Anything).Return(instanceSizes, nil)
			},
			want: true,
		},
		{
			name:          "old and new values are not equivalent",
			cloudProvider: "aws",
			oldValue:      "t2.micro",
			newValue:      "Small",
			mock: func(svc *mocks.MockContentService) {
				svc.EXPECT().InstanceSizes(mock.Anything).Return(instanceSizes, nil)
			},
			want: false,
		},
		{
			name:          "cloud provider is empty",
			cloudProvider: "",
			oldValue:      "t2.micro",
			newValue:      "Tiny",
			want:          false,
		},
		{
			name:          "cloud provider is not found",
			cloudProvider: "foo",
			oldValue:      "t2.micro",
			newValue:      "Tiny",
			mock: func(svc *mocks.MockContentService) {
				svc.EXPECT().InstanceSizes(mock.Anything).Return(instanceSizes, nil)
			},
			want: false,
		},
		{
			name:          "loading instance sizes failed",
			cloudProvider: "foo",
			oldValue:      "t2.micro",
			newValue:      "Tiny",
			mock: func(svc *mocks.MockContentService) {
				svc.EXPECT().InstanceSizes(mock.Anything).Return(nil, errors.New("something went wrong"))
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc := mocks.NewMockContentService(t)

			if tt.mock != nil {
				tt.mock(svc)
			}

			got := checkInstanceSizeEquivalence(ctx, svc, tt.cloudProvider, tt.oldValue, tt.newValue)
			if got != tt.want {
				t.Errorf("checkInstanceSizeEquivalence() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_vendorSuppressor(t *testing.T) {
	tests := []struct {
		name     string
		oldValue string
		newValue string
		want     bool
	}{
		{
			name:     "old and new values are exactly the same unaliased",
			oldValue: "percona",
			newValue: "percona",
			want:     true,
		},
		{
			name:     "old and new values are exactly the same aliased",
			oldValue: "mysql",
			newValue: "mysql",
			want:     true,
		},
		{
			name:     "old and new values are equivalent, alias + non-alias",
			oldValue: "mysql",
			newValue: "percona",
			want:     true,
		},
		{
			name:     "old and new values are equivalent, non-alias + alias",
			oldValue: "percona",
			newValue: "mysql",
			want:     true,
		},
		{
			name:     "old and new values are not equivalent",
			oldValue: "percona",
			newValue: "mariadb",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := vendorSuppressor("", tt.oldValue, tt.newValue, nil)

			assert.Equal(t, tt.want, got)
		})
	}
}
