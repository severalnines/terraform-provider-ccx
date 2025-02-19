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

type dbVendorsInfoResponse struct {
	Database struct {
		Vendors []struct {
			Name     string   `json:"name"`
			Version  string   `json:"version"`
			Versions []string `json:"versions"`
			Code     string   `json:"code"`
			NumNodes []int    `json:"num_nodes"`
			Types    []struct {
				Name string `json:"name"`
				Code string `json:"code"`
			} `json:"types"`
		} `json:"vendors"`
	} `json:"database"`
}

func (svc *ContentService) DBVendors(ctx context.Context) ([]ccx.DBVendorInfo, error) {
	var rs dbVendorsInfoResponse

	err := svc.client.Get(ctx, "/api/content/api/v1/deploy-wizard", &rs)
	if err != nil {
		return nil, err
	}

	vendors := make([]ccx.DBVendorInfo, 0, len(rs.Database.Vendors))

	for _, v := range rs.Database.Vendors {
		vendor := ccx.DBVendorInfo{
			Name:           v.Name,
			Code:           v.Code,
			DefaultVersion: v.Version,
			Versions:       v.Versions,
			NumNodes:       v.NumNodes,
		}

		for _, t := range v.Types {
			vendor.Types = append(vendor.Types, ccx.DBVendorInfoType{
				Name: t.Name,
				Code: t.Code,
			})
		}

		vendors = append(vendors, vendor)
	}

	return vendors, nil
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

type volumeTypesResponse struct {
	Instance struct {
		VolumeTypes map[string][]struct {
			Code string `json:"code"`
		} `json:"volume_types"`
	} `json:"instance"`
}

func (svc *ContentService) VolumeTypes(ctx context.Context, cloud string) ([]string, error) {
	var rs volumeTypesResponse

	err := svc.client.Get(ctx, "/api/content/api/v1/deploy-wizard", &rs)
	if err != nil {
		return nil, err
	}

	vt, ok := rs.Instance.VolumeTypes[cloud]
	if !ok {
		return nil, fmt.Errorf("no volume types found for cloud %q", cloud)
	}

	ls := make([]string, 0, len(vt))
	for _, v := range vt {
		ls = append(ls, v.Code)
	}

	return ls, nil
}
