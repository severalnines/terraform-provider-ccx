package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
)

func (svc *DatastoreService) SetMaintenanceSettings(ctx context.Context, storeID string, settings ccx.MaintenanceSettings) error {
	ur := updateRequest{
		Maintenance: &maintenance{
			DayOfWeek: uint32(settings.DayOfWeek),
			StartHour: uint64(settings.StartHour),
			EndHour:   uint64(settings.EndHour),
		},
	}

	_, err := svc.client.Do(ctx, http.MethodPatch, "/api/prov/api/v2/cluster/"+storeID, ur)
	if err != nil {
		return fmt.Errorf("setting maintenance settings: %w", err)
	}

	return nil
}
