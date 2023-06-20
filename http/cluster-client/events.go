package cluster_client

import (
	"context"

	ccxprov "github.com/severalnines/terraform-provider-ccx"
)

type onLoadCallback func(context.Context, ccxprov.Cluster)

func (o onLoadCallback) OnCalled(ctx context.Context, c ccxprov.Cluster) {
	o(ctx, c)
}
