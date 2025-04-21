package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
)

func (svc *VpcService) Delete(ctx context.Context, id string) error {
	_, err := svc.httpcli.Do(ctx, http.MethodDelete, "/api/vpc/api/v2/vpcs"+"/"+id, nil)
	if errors.Is(err, ccx.ErrResourceNotFound) {
		tflog.Warn(ctx, "deleting vpc: not found", map[string]interface{}{"id": id})
		return nil
	} else if err != nil {
		return fmt.Errorf("deleting vpc: %w", err)
	}

	return nil
}
