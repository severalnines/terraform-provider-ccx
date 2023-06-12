package vpc_client

import (
	"context"

	ccxprov "github.com/severalnines/terraform-provider-ccx"
)

func (cli *Client) Update(_ context.Context, _ ccxprov.VPC) (*ccxprov.VPC, error) {
	return nil, ccxprov.UpdateNotSupportedErr
}
