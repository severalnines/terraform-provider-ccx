package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
	"github.com/severalnines/terraform-provider-ccx/internal/lib"
)

func (svc *DatastoreService) SetMaintenanceSettings(ctx context.Context, storeID string, settings ccx.MaintenanceSettings) error {
	ur := updateRequest{
		Maintenance: &maintenance{
			DayOfWeek: uint32(settings.DayOfWeek),
			StartHour: uint64(settings.StartHour),
			EndHour:   uint64(settings.EndHour),
		},
	}

	res, err := svc.client.Do(ctx, http.MethodPatch, "/api/prov/api/v2/cluster/"+storeID, ur)
	if err != nil {
		return errors.Join(ccx.RequestSendingErr, err)
	}

	if res.StatusCode == http.StatusBadRequest {
		return lib.ErrorFromErrorResponse(res.Body)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: status = %d", ccx.ResponseStatusFailedErr, res.StatusCode)
	}

	return nil
}
