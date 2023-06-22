package vpc_client

import (
	"context"

	"github.com/severalnines/terraform-provider-ccx/ccx"
)

func (cli *Client) Update(_ context.Context, _ ccx.VPC) (*ccx.VPC, error) {
	return nil, ccx.UpdateNotSupportedErr
}
