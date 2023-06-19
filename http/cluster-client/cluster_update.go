package cluster_client

import (
	"context"

	ccxprov "github.com/severalnines/terraform-provider-ccx"
)

func (cli *Client) Update(_ context.Context, _ ccxprov.Cluster) (*ccxprov.Cluster, error) {
	return nil, ccxprov.UpdateNotSupportedErr
}
