package resources

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
)

// caseInsensitiveSuppressor suppresses diff for fields when the values are case-insensitive equal
func caseInsensitiveSuppressor(k, oldValue, newValue string, d *schema.ResourceData) bool {
	return strings.EqualFold(oldValue, newValue)
}

// checkInstanceSizeEquivalence will check if (oldValue or its alias) is equal to (newValue or its alias)
// in the list of sizes fetched from ccx.ContentService for a given cloudProvider
func checkInstanceSizeEquivalence(ctx context.Context, svc ccx.ContentService, cloudProvider, oldValue, newValue string) bool {
	if oldValue != "" && newValue != "" && strings.EqualFold(oldValue, newValue) {
		return true
	}

	if cloudProvider == "" {
		tflog.Error(ctx, "failed to read cloud_provider")
		return false
	}

	instanceSizes, err := svc.InstanceSizes(ctx)
	if err != nil {
		tflog.Error(ctx, "failed to load instance sizes")
		return false
	}

	sizes, ok := instanceSizes[cloudProvider]
	if !ok {
		tflog.Error(ctx, "cloud_provider not found in instance sizes")
		return false
	}

	for _, size := range sizes {
		isOld := strings.EqualFold(size.Code, oldValue) || strings.EqualFold(size.Type, oldValue)
		isNew := strings.EqualFold(size.Code, newValue) || strings.EqualFold(size.Type, newValue)

		if isOld && isNew {
			return true
		}

		if isOld && !isNew || isNew && !isOld {
			return false
		}
	}

	return false
}

func firewallDiffSupressor(_, _, _ string, d *schema.ResourceData) bool {
	old, nw, err := getFirewallsOldNew(d)
	if err != nil {
		tflog.Error(context.Background(), "failed to get old and new firewalls", map[string]any{"err": err.Error()})
		return false
	}

	return firewallsSame(old, nw)
}
