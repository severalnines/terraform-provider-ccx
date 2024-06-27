package api

import (
	"context"
	"fmt"

	"github.com/severalnines/terraform-provider-ccx/ccx"
)

type loadAllResponse []struct {
	UUID             string  `json:"uuid"`
	Region           string  `json:"region"`
	CloudProvider    string  `json:"cloud_provider"`
	InstanceSize     string  `json:"instance_size"`
	InstanceIOPS     *uint64 `json:"iops"`
	DiskSize         *uint64 `json:"disk_size"`
	DiskType         *string `json:"disk_type"`
	DbVendor         string  `json:"database_vendor"`
	DbVersion        string  `json:"database_version"`
	Name             string  `json:"cluster_name"`
	Type             string  `json:"cluster_type"`
	Size             int64   `json:"cluster_size"`
	HighAvailability bool    `json:"high_availability"`
	Vpc              *struct {
		VpcUUID string `json:"vpc_uuid"`
	} `json:"vpc"`
	Tags []string `json:"tags"`
	AZS  []string `json:"azs"`
}

func fromLoadAllResponse(r loadAllResponse) map[string]ccx.Datastore {
	c := make(map[string]ccx.Datastore)

	for _, info := range r {
		var vpcUUID string
		if info.Vpc != nil {
			vpcUUID = info.Vpc.VpcUUID
		}

		c[info.UUID] = ccx.Datastore{
			ID:                info.UUID,
			Name:              info.Name,
			Size:              info.Size,
			DBVendor:          info.DbVendor,
			DBVersion:         info.DbVersion,
			Type:              info.Type,
			Tags:              info.Tags,
			CloudProvider:     info.CloudProvider,
			CloudRegion:       info.Region,
			InstanceSize:      info.InstanceSize,
			VolumeType:        pString(info.DiskType),
			VolumeSize:        int64(pUint64(info.DiskSize)),
			VolumeIOPS:        int64(pUint64(info.InstanceIOPS)),
			NetworkType:       "", // todo
			HAEnabled:         info.HighAvailability,
			VpcUUID:           vpcUUID,
			AvailabilityZones: info.AZS,
		}
	}

	return c
}

func (svc *DatastoreService) Read(ctx context.Context, id string) (*ccx.Datastore, error) {
	defer svc.mut.Unlock()
	svc.mut.Lock()

	c, ok := svc.stores[id]
	if !ok {
		return nil, ccx.ResourceNotFoundErr
	}

	if p, err := svc.GetParameters(ctx, id); err == nil {
		c.DbParams = p
	} else {
		return nil, fmt.Errorf("getting parameters: %w", err)
	}

	if fw, err := svc.GetFirewallRules(ctx, id); err == nil {
		c.FirewallRules = fw
	} else {
		return nil, fmt.Errorf("getting firewall rules: %w", err)
	}

	return &c, nil
}

func (svc *DatastoreService) LoadAll(ctx context.Context) error {
	url := svc.baseURL + "/api/deployment/api/v1/deployments"

	var rs loadAllResponse

	if err := httpGet(ctx, svc.auth, url, &rs); err != nil {
		return err
	}

	stores := fromLoadAllResponse(rs)

	svc.mut.Lock()
	defer svc.mut.Unlock()

	svc.stores = stores

	return nil
}
