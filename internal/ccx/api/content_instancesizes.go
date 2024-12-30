package api

import (
	"context"
	"fmt"

	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
)

type instanceSizesResponse struct {
	Instance struct {
		InstanceSizes map[string][]ccx.InstanceSize `json:"instance_sizes"`
	} `json:"instance"`
}

func (svc *ContentService) InstanceSizes(ctx context.Context) (map[string][]ccx.InstanceSize, error) {
	var rs instanceSizesResponse

	err := svc.client.Get(ctx, "/api/content/api/v1/deploy-wizard", &rs)
	if err != nil {
		return nil, err
	}

	return rs.Instance.InstanceSizes, nil
}

type availabilityZonesResponse struct {
	Network struct {
		AvailabilityZones map[string]map[string][]struct {
			Code string `json:"code"`
		} `json:"availability_zones"`
	} `json:"network"`
}

func (svc *ContentService) AvailabilityZones(ctx context.Context, provider, region string) ([]string, error) {
	var rs availabilityZonesResponse

	err := svc.client.Get(ctx, "/api/content/api/v1/deploy-wizard", &rs)
	if err != nil {
		return nil, err
	}

	p, ok := rs.Network.AvailabilityZones[provider]
	if !ok {
		return nil, fmt.Errorf("no availability zones found for provider %q", provider)
	}

	r, ok := p[region]
	if !ok {
		return nil, fmt.Errorf("no availability zones found for provider %q in region %q", provider, region)
	}

	ls := make([]string, 0, len(r))
	for _, az := range r {
		ls = append(ls, az.Code)
	}

	return ls, nil
}
