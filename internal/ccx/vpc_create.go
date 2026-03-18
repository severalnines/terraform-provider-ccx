package ccx

import (
	"context"
	"fmt"
	"net/http"
)

type createVpcRequest struct {
	Name          string `json:"name"`
	Cloudspace    string `json:"cloudspace"`
	CloudProvider string `json:"cloud"`
	Region        string `json:"region"`
	CidrIpv4Block string `json:"cidr_ipv4_block"`
}

func CreateRequestFromVpc(v VPC) createVpcRequest {
	return createVpcRequest{
		Name:          v.Name,
		Cloudspace:    v.CloudSpace,
		CloudProvider: v.CloudProvider,
		Region:        v.Region,
		CidrIpv4Block: v.CidrIpv4Block,
	}
}

func (svc *VPCsClient) Create(ctx context.Context, vpc VPC) (*VPC, error) {
	cr := CreateRequestFromVpc(vpc)

	res, err := svc.httpcli.Do(ctx, http.MethodPost, "/api/vpc/api/v2/vpcs", cr)
	if err != nil {
		return nil, fmt.Errorf("creating vpc: %w", err)
	}

	var rs vpcResponse
	if err := DecodeJsonInto(res.Body, &rs); err != nil {
		return nil, err
	}

	newVPC := vpcFromResponse(rs)

	return &newVPC, nil
}
