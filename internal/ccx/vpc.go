package ccx

type VPCsClient struct {
	httpcli HTTPClient
}

// NewVPCsClient creates a new VPC VpcService
func NewVPCsClient(httpcli HTTPClient) VPCsService {
	return &VPCsClient{
		httpcli: httpcli,
	}
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

func vpcFromResponse(r vpcResponse) VPC {
	if r.VPC == nil {
		return VPC{}
	}

	return VPC{
		ID:            r.VPC.ID,
		Name:          r.VPC.Name,
		CloudProvider: r.VPC.CloudProvider,
		CloudSpace:    r.VPC.Cloudspace,
		Region:        r.VPC.Region,
		CidrIpv4Block: r.VPC.CidrIpv4Block,
	}
}
