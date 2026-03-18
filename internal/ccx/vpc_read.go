package ccx

import (
	"context"
)

func (svc *VPCsClient) Read(ctx context.Context, id string) (*VPC, error) {
	var rs vpcResponse

	if err := svc.httpcli.Get(ctx, "/api/vpc/api/v2/vpcs/"+id, &rs); err != nil {
		return nil, err
	}

	vpc := vpcFromResponse(rs)

	return &vpc, nil
}
