package resources

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
)

func getNotifications(d *schema.ResourceData) ccx.Notifications {
	return ccx.Notifications{
		Enabled: getBool(d, "notifications_enabled"),
		Emails:  getStrings(d, "notifications_emails"),
	}
}

func setNotifications(d *schema.ResourceData, n ccx.Notifications) error {
	if err := d.Set("notifications_enabled", n.Enabled); err != nil {
		return fmt.Errorf("setting notifications_enabled: %w", err)
	}

	if err := setStrings(d, "notifications_emails", n.Emails); err != nil {
		return fmt.Errorf("setting notifications_emails: %w", err)
	}

	return nil
}
