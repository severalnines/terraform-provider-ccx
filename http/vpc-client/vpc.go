package vpc_client

import (
	ccxprov "github.com/severalnines/terraform-provider-ccx"
	chttp "github.com/severalnines/terraform-provider-ccx/http"
)

var _ ccxprov.VPCService = &Client{}

type Client struct {
	auth chttp.Authorizer
	conn *chttp.ConnectionParameters
}

// New creates a new clusters Client
func New(authorizer chttp.Authorizer, opts ...chttp.ParameterOption) *Client {
	p := chttp.Parameters(opts...)

	var c = Client{
		auth: authorizer,
		conn: p,
	}

	return &c
}
