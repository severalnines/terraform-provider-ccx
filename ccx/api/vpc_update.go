package api

import (
	"context"

	"github.com/severalnines/terraform-provider-ccx/ccx"
)

func (svc *VpcService) Update(_ context.Context, _ ccx.VPC) (*ccx.VPC, error) {
	return nil, ccx.UpdateNotSupportedErr
}
