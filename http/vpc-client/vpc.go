package vpc_client

import (
	"github.com/severalnines/terraform-provider-ccx/ccx"
	chttp "github.com/severalnines/terraform-provider-ccx/http"
)

var _ ccx.VPCService = &Client{}

type Client struct {
	auth chttp.Authorizer
	conn *chttp.ConnectionParameters
}

// New creates a new VPC Client
func New(authorizer chttp.Authorizer, opts ...chttp.ParameterOption) *Client {
	p := chttp.Parameters(opts...)

	var c = Client{
		auth: authorizer,
		conn: p,
	}

	return &c
}
