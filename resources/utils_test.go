package resources

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
)

func checkTFDatastore(name string, d ccx.Datastore) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		root := state.RootModule()
		if root == nil {
			return errors.New("root module is nil")
		}

		name = "ccx_datastore." + name

		ds, ok := root.Resources[name]
		if !ok {
			return fmt.Errorf("resource %q not found", name)
		}

		if ds.Type != "ccx_datastore" {
			return fmt.Errorf("expected resource %q to be of type ccx_datastore, got %q", name, ds.Type)
		}

		if ds.Primary.Attributes["id"] != d.ID {
			return fmt.Errorf("expected ID %q, got %q", d.ID, ds.Primary.Attributes["id"])
		}

		if err := testMapElementEqual(ds.Primary.Attributes, "id", d.ID); err != nil {
			return err
		}

		if err := testMapElementEqual(ds.Primary.Attributes, "name", d.Name); err != nil {
			return err
		}

		if err := testMapElementEqual(ds.Primary.Attributes, "size", strconv.FormatInt(d.Size, 10)); err != nil {
			return err
		}

		if err := testMapElementEqual(ds.Primary.Attributes, "db_vendor", d.DBVendor); err != nil {
			return err
		}

		if err := testMapElementEqual(ds.Primary.Attributes, "db_version", d.DBVersion); err != nil {
			return err
		}

		if err := testMapElementEqual(ds.Primary.Attributes, "type", d.Type); err != nil {
			return err
		}

		if err := testMapElementEqualSlice(ds.Primary.Attributes, "tags", d.Tags); err != nil {
			return err
		}

		if err := testMapElementEqual(ds.Primary.Attributes, "cloud_provider", d.CloudProvider); err != nil {
			return err
		}

		if err := testMapElementEqual(ds.Primary.Attributes, "cloud_region", d.CloudRegion); err != nil {
			return err
		}

		if err := testMapElementEqual(ds.Primary.Attributes, "instance_size", d.InstanceSize); err != nil {
			return err
		}

		if err := testMapElementEqual(ds.Primary.Attributes, "volume_type", d.VolumeType); err != nil {
			return err
		}

		if err := testMapElementEqual(ds.Primary.Attributes, "volume_size", strconv.FormatUint(d.VolumeSize, 10)); err != nil {
			return err
		}

		if err := testMapElementEqual(ds.Primary.Attributes, "volume_iops", strconv.FormatUint(d.VolumeIOPS, 10)); err != nil {
			return err
		}

		if err := testMapElementEqual(ds.Primary.Attributes, "network_vpc_uuid", d.VpcUUID); err != nil {
			return err
		}

		if err := testMapElementEqual(ds.Primary.Attributes, "network_ha_enabled", strconv.FormatBool(d.HAEnabled)); err != nil {
			return err
		}

		if err := testMapElementEqual(ds.Primary.Attributes, "notifications_enabled", strconv.FormatBool(d.Notifications.Enabled)); err != nil {
			return err
		}

		if err := testMapElementEqualSlice(ds.Primary.Attributes, "notifications_emails", d.Notifications.Emails); err != nil {
			return err
		}

		if err := testMapElementEqualMaintenanceSettings(ds.Primary.Attributes, d.MaintenanceSettings); err != nil {
			return err
		}

		if err := testMapElementEqual(ds.Primary.Attributes, "primary_url", d.PrimaryUrl); err != nil {
			return err
		}

		if err := testMapElementEqual(ds.Primary.Attributes, "primary_dsn", d.PrimaryDsn); err != nil {
			return err
		}

		if err := testMapElementEqual(ds.Primary.Attributes, "replica_url", d.ReplicaUrl); err != nil {
			return err
		}

		if err := testMapElementEqual(ds.Primary.Attributes, "replica_dsn", d.ReplicaDsn); err != nil {
			return err
		}

		if err := testMapElementEqual(ds.Primary.Attributes, "username", d.Username); err != nil {
			return err
		}

		if err := testMapElementEqual(ds.Primary.Attributes, "password", d.Password); err != nil {
			return err
		}

		if err := testMapElementEqual(ds.Primary.Attributes, "dbname", d.DbName); err != nil {
			return err
		}

		return nil
	}
}

func testMapElementAbsent(attributes map[string]string, key string) error {
	_, ok := attributes[key]
	if ok {
		return fmt.Errorf("attribute %q should be absent", key)
	}

	return nil
}

func testMapElementEqual(attributes map[string]string, key, value string) error {
	v, ok := attributes[key]
	if !ok {
		return fmt.Errorf("attribute %q not found", key)
	}

	if v != value {
		return fmt.Errorf("attribute %q expected value %q, got %q", key, value, v)
	}

	return nil
}

func testMapElementEqualSlice(attributes map[string]string, key string, items []string) error {
	if len(items) == 0 {
		return testMapElementEqual(attributes, key+".#", "0")
	}

	var n int

	if c, ok := attributes[key+".#"]; !ok {
		return fmt.Errorf("slice count %q not found", key+".#")
	} else if v, err := strconv.Atoi(c); err != nil {
		return fmt.Errorf("slice count %q is not a number: %w", key+".#", err)
	} else {
		n = v
	}

	var stateItems []string

	for i := 0; i < n; i++ {
		if v, ok := attributes[key+"."+strconv.Itoa(i)]; !ok {
			return fmt.Errorf("slice %q item %d not found", key, i)
		} else {
			stateItems = append(stateItems, v)
		}
	}

	if !isSubsetOf(stateItems, items) {
		return fmt.Errorf("slice %q items are not correct have %+v not a subsetof %+v", key, stateItems, items)
	}

	return nil
}

func testMapElementEqualMaintenanceSettings(attributes map[string]string, m *ccx.MaintenanceSettings) error {
	if m == nil {
		if err := testMapElementAbsent(attributes, "maintenance_day_of_week"); err != nil {
			return err
		}

		if err := testMapElementAbsent(attributes, "maintenance_start_hour"); err != nil {
			return err
		}

		if err := testMapElementAbsent(attributes, "maintenance_end_hour"); err != nil {
			return err
		}

		return nil
	}

	if err := testMapElementEqual(attributes, "maintenance_day_of_week", strconv.Itoa(int(m.DayOfWeek))); err != nil {
		return err
	}

	if err := testMapElementEqual(attributes, "maintenance_start_hour", strconv.Itoa(m.StartHour)); err != nil {
		return err
	}

	if err := testMapElementEqual(attributes, "maintenance_end_hour", strconv.Itoa(m.EndHour)); err != nil {
		return err
	}

	return nil
}
