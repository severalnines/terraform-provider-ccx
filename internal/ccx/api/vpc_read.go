package api

import (
	"context"

	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
)

func (svc *VpcService) Read(ctx context.Context, id string) (*ccx.VPC, error) {
	var rs vpcResponse

	if err := svc.httpcli.Get(ctx, "/api/vpc/api/v2/vpcs/"+id, &rs); err != nil {
		return nil, err
	}

	vpc := vpcFromResponse(rs)

	return &vpc, nil
}
