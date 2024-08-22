package api

import (
	"context"
	"time"

	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
)

type getHostsResponse struct {
	UUID  string
	Hosts []struct {
		ID            string    `json:"host_uuid"`
		CreatedAt     time.Time `json:"created_at"`
		CloudProvider string    `json:"cloud_provider"`
		AZ            string    `json:"host_az"`
		InstanceType  string    `json:"instance_type"`
		DiskType      string    `json:"disk_type"`
		DiskSize      uint64    `json:"disk_size"`
		Role          string    `json:"role"`
		Region        struct {
			Code string `json:"code"`
		} `json:"region"`
	} `json:"database_nodes"`
}

func (svc *DatastoreService) GetHosts(ctx context.Context, clusterID string) ([]ccx.Host, error) {
	var rs getHostsResponse

	err := svc.client.Get(ctx, "/api/deployment/v2/data-stores/"+clusterID+"/nodes", &rs)
	if err != nil {
		return nil, err
	}

	ls := make([]ccx.Host, 0, len(rs.Hosts))

	for i := range rs.Hosts {
		h := ccx.Host{
			ID:            rs.Hosts[i].ID,
			CreatedAt:     rs.Hosts[i].CreatedAt,
			CloudProvider: rs.Hosts[i].CloudProvider,
			AZ:            rs.Hosts[i].AZ,
			InstanceType:  rs.Hosts[i].InstanceType,
			DiskType:      rs.Hosts[i].DiskType,
			DiskSize:      rs.Hosts[i].DiskSize,
			Role:          rs.Hosts[i].Role,
			Region:        rs.Hosts[i].Region.Code,
		}

		ls = append(ls, h)
	}

	return ls, nil
}
