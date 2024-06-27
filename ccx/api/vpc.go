package api

import (
	"github.com/severalnines/terraform-provider-ccx/ccx"
)

type VpcService struct {
	auth    authorizer
	baseURL string
}

// Vpcs creates a new VPC VpcService
func Vpcs(baseURL, clientID, clientSecret string) *VpcService {
	a := tokenAuthorizer{
		id:      clientID,
		secret:  clientSecret,
		baseURL: baseURL,
	}

	var c = VpcService{
		auth:    a,
		baseURL: baseURL,
	}

	return &c
}

type vpcResponse struct {
	VPC *struct {
		ID            string `json:"id"`
		Name          string `json:"name"`
		Cloudspace    string `json:"cloudspace"`
		CloudProvider string `json:"cloud"`
		Region        string `json:"region"`
		CidrIpv4Block string `json:"cidr_ipv4_block"`
	} `json:"vpc"`
}

func vpcFromResponse(r vpcResponse) ccx.VPC {
	if r.VPC == nil {
		return ccx.VPC{}
	}

	return ccx.VPC{
		ID:            r.VPC.ID,
		Name:          r.VPC.Name,
		CloudProvider: r.VPC.CloudProvider,
		CloudSpace:    r.VPC.Cloudspace,
		Region:        r.VPC.Region,
		CidrIpv4Block: r.VPC.CidrIpv4Block,
	}
}
