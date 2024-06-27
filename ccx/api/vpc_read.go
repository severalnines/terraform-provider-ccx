package api

import (
	"context"

	"github.com/severalnines/terraform-provider-ccx/ccx"
)

func (svc *VpcService) Read(ctx context.Context, id string) (*ccx.VPC, error) {
	url := svc.baseURL + "/api/vpc/api/v2/vpcs/" + id

	var rs vpcResponse

	if err := httpGet(ctx, svc.auth, url, &rs); err != nil {
		return nil, err
	}

	vpc := vpcFromResponse(rs)

	return &vpc, nil
}
