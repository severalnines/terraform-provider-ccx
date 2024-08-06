package resources

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
)

type notifications struct{}

func (n notifications) Schema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable or disable notifications. Default is false",
			},
			"emails": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of email addresses to send notifications to",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func getNotifications(d *schema.ResourceData) ccx.Notifications {
	return ccx.Notifications{
		Enabled: getBool(d, "notifications_enabled"),
		Emails:  getStrings(d, "notifications_emails"),
	}
}

func setNotifications(d *schema.ResourceData, n ccx.Notifications) error {
	emails := make([]any, 0, len(n.Emails))
	for _, email := range n.Emails {
		emails = append(emails, email)
	}

	if err := d.Set("notifications_enabled", n.Enabled); err != nil {
		return fmt.Errorf("setting notifications_enabled: %w", err)
	}

	if err := d.Set("notifications_emails", n.Emails); err != nil {
		return fmt.Errorf("setting notifications_emails: %w", err)
	}

	return nil
}
