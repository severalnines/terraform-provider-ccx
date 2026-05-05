package ccx

import (
	"context"
)

func (svc *VPCsClient) Update(_ context.Context, _ VPC) (*VPC, error) {
	return nil, ErrUpdateNotSupported
}
