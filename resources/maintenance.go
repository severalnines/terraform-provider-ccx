package resources

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
)

func getMaintenanceSettings(d *schema.ResourceData) *ccx.MaintenanceSettings {
	if !allKeysSet(d, "maintenance_day_of_week", "maintenance_start_hour", "maintenance_end_hour") {
		return nil
	}

	return &ccx.MaintenanceSettings{
		DayOfWeek: int32(getInt(d, "maintenance_day_of_week")),
		StartHour: int(getInt(d, "maintenance_start_hour")),
		EndHour:   int(getInt(d, "maintenance_end_hour")),
	}
}

func setMaintenanceSettings(d *schema.ResourceData, m ccx.MaintenanceSettings) error {
	if err := d.Set("maintenance_day_of_week", int(m.DayOfWeek)); err != nil {
		return fmt.Errorf("setting maintenance_day_of_week: %w", err)
	}

	if err := d.Set("maintenance_start_hour", m.StartHour); err != nil {
		return fmt.Errorf("setting maintenance_start_hour: %w", err)
	}

	if err := d.Set("maintenance_end_hour", m.EndHour); err != nil {
		return fmt.Errorf("setting maintenance_end_hour: %w", err)
	}

	return nil
}
