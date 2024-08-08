package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
)

func (svc *VpcService) Delete(ctx context.Context, id string) error {
	res, err := svc.httpcli.Do(ctx, http.MethodDelete, "/api/vpc/api/v2/vpcs"+"/"+id, nil)
	if err != nil {
		return errors.Join(ccx.RequestSendingErr, err)
	}

	if res.StatusCode != http.StatusAccepted {
		return fmt.Errorf("%w: status = %d", ccx.ResponseStatusFailedErr, res.StatusCode)
	}

	return nil
}
