package resources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
	"github.com/stretchr/testify/mock"
)

func TestDatastore_Create(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		m, p := mockProvider()

		m.datastore.EXPECT().Create(mock.Anything, ccx.Datastore{
			Name:              "luna",
			Size:              1,
			DBVendor:          "postgres",
			Type:              "postgres_streaming",
			Tags:              []string{"new", "test"},
			CloudProvider:     "aws",
			CloudRegion:       "eu-north-1",
			InstanceSize:      "m5.large",
			VolumeType:        "gp2",
			VolumeSize:        80,
			AvailabilityZones: []string{},
			DbParams:          nil,
			FirewallRules:     []ccx.FirewallRule{},
			NetworkType:       "public",
			Notifications: ccx.Notifications{
				Enabled: false,
				Emails:  []string{},
			},
		}).Return(&ccx.Datastore{
			ID:            "datastore-id",
			Name:          "luna",
			Size:          1,
			DBVendor:      "postgres",
			DBVersion:     "15",
			Type:          "postgres_streaming",
			Tags:          []string{"new", "test", "postgres", "15", "postgres_streaming", "aws", "eu-north-1"},
			CloudProvider: "aws",
			CloudRegion:   "eu-north-1",
			InstanceSize:  "m5.large",
			VolumeType:    "gp2",
			VolumeSize:    80,
			VolumeIOPS:    0,
			Notifications: ccx.Notifications{
				Enabled: false,
				Emails:  []string{"user@getccx.com"},
			},
			MaintenanceSettings: &ccx.MaintenanceSettings{
				DayOfWeek: 1,
				StartHour: 0,
				EndHour:   2,
			},
		}, nil)

		m.datastore.EXPECT().Read(mock.Anything, "datastore-id").Return(&ccx.Datastore{
			ID:            "datastore-id",
			Name:          "luna",
			Size:          1,
			DBVendor:      "postgres",
			DBVersion:     "15",
			Type:          "postgres_streaming",
			Tags:          []string{"new", "test", "postgres", "15", "postgres_streaming", "aws", "eu-north-1"},
			CloudProvider: "aws",
			CloudRegion:   "eu-north-1",
			InstanceSize:  "m5.large",
			VolumeType:    "gp2",
			VolumeSize:    80,
			VolumeIOPS:    0,
			Notifications: ccx.Notifications{
				Enabled: false,
				Emails:  []string{"user@getccx.com"},
			},
			MaintenanceSettings: &ccx.MaintenanceSettings{
				DayOfWeek: 1,
				StartHour: 0,
				EndHour:   2,
			},
		}, nil)

		m.datastore.EXPECT().Delete(mock.Anything, "datastore-id").Return(nil)

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			PreCheck: func() {
			},
			ProviderFactories: map[string]func() (*schema.Provider, error){
				"ccx": func() (*schema.Provider, error) {
					return p, nil
				},
			},
			Steps: []resource.TestStep{
				{
					Config: `
resource "ccx_datastore" "luna" {
  name           = "luna"
  size           = 1
  db_vendor      = "postgres"
  tags           = ["new", "test"]
  cloud_provider = "aws"
  cloud_region   = "eu-north-1"
  instance_size  = "m5.large"
  volume_size    = 80
  volume_type    = "gp2"
  network_type   = "public"
}
`,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("ccx_datastore.luna", "id", "datastore-id"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "size", "1"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "db_vendor", "postgres"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "db_version", "15"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "type", "postgres_streaming"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "tags.#", "2"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "tags.0", "new"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "tags.1", "test"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "cloud_provider", "aws"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "cloud_region", "eu-north-1"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "instance_size", "m5.large"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "volume_type", "gp2"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "volume_size", "80"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "volume_iops", "0"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "network_type", "public"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "network_vpc_uuid", ""),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "network_ha_enabled", "false"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "notifications_enabled", "false"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "notifications_emails.#", "1"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "notifications_emails.0", "user@getccx.com"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "maintenance_day_of_week", "1"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "maintenance_start_hour", "0"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "maintenance_end_hour", "2"),
					),
				},
			},
		})

		m.AssertExpectations(t)
	})

	t.Run("scaling, with db_type replication", func(t *testing.T) {
		// the api payload to create the datastore uses the db_type field: replication
		// whereas the server returns the type field: Replication
		// this should not trigger a diff
		//
		// also testing the scaling part

		m, p := mockProvider()

		createdDatastore := &ccx.Datastore{
			ID:            "datastore-id",
			Name:          "luna",
			Size:          1,
			DBVendor:      "percona",
			DBVersion:     "8",
			Type:          "Replication",
			Tags:          []string{"new", "test", "percona", "8", "replication", "aws", "eu-north-1"},
			CloudProvider: "aws",
			CloudRegion:   "eu-north-1",
			InstanceSize:  "m5.large",
			VolumeType:    "gp2",
			VolumeSize:    80,
			VolumeIOPS:    0,
			Notifications: ccx.Notifications{
				Enabled: false,
				Emails:  []string{"user@getccx.com"},
			},
			MaintenanceSettings: &ccx.MaintenanceSettings{
				DayOfWeek: 1,
				StartHour: 0,
				EndHour:   2,
			},
		}

		updatedDatastore := &ccx.Datastore{
			ID:            "datastore-id",
			Name:          "luna",
			Size:          2,
			DBVendor:      "percona",
			DBVersion:     "8",
			Type:          "Replication",
			Tags:          []string{"new", "test"},
			CloudProvider: "aws",
			CloudRegion:   "eu-north-1",
			InstanceSize:  "m5.large",
			VolumeType:    "gp2",
			VolumeSize:    80,
			VolumeIOPS:    0,
			Notifications: ccx.Notifications{
				Enabled: false,
				Emails:  []string{"user@getccx.com"},
			},
			NetworkType:       "public",
			AvailabilityZones: []string{},
			FirewallRules:     []ccx.FirewallRule{},
			MaintenanceSettings: &ccx.MaintenanceSettings{
				DayOfWeek: 1,
				StartHour: 0,
				EndHour:   2,
			},
		}

		m.datastore.EXPECT().Create(mock.Anything, ccx.Datastore{
			Name:              "luna",
			Size:              1,
			DBVendor:          "percona",
			Type:              "replication",
			Tags:              []string{"new", "test"},
			CloudProvider:     "aws",
			CloudRegion:       "eu-north-1",
			InstanceSize:      "m5.large",
			VolumeType:        "gp2",
			VolumeSize:        80,
			AvailabilityZones: []string{},
			DbParams:          nil,
			FirewallRules:     []ccx.FirewallRule{},
			NetworkType:       "public",
			Notifications: ccx.Notifications{
				Enabled: false,
				Emails:  []string{},
			},
		}).Return(createdDatastore, nil).Once()

		var updated bool

		m.datastore.EXPECT().Read(mock.Anything, "datastore-id").
			RunAndReturn(func(_ context.Context, _ string) (*ccx.Datastore, error) {
				if updated {
					return updatedDatastore, nil
				}

				return createdDatastore, nil
			})

		m.datastore.EXPECT().Update(mock.Anything,
			ccx.Datastore{
				ID:            "datastore-id",
				Name:          "luna",
				Size:          1,
				DBVendor:      "percona",
				DBVersion:     "8",
				Type:          "Replication",
				Tags:          []string{"new", "test", "percona", "8", "replication", "aws", "eu-north-1"},
				CloudProvider: "aws",
				CloudRegion:   "eu-north-1",
				InstanceSize:  "m5.large",
				VolumeType:    "gp2",
				VolumeSize:    80,
				VolumeIOPS:    0,
				Notifications: ccx.Notifications{
					Enabled: false,
					Emails:  []string{"user@getccx.com"},
				},
				MaintenanceSettings: &ccx.MaintenanceSettings{
					DayOfWeek: 1,
					StartHour: 0,
					EndHour:   2,
				},
			},
			ccx.Datastore{
				ID:            "datastore-id",
				Name:          "luna",
				Size:          2,
				DBVendor:      "percona",
				DBVersion:     "8",
				Type:          "Replication",
				Tags:          []string{"new", "test"},
				CloudProvider: "aws",
				CloudRegion:   "eu-north-1",
				InstanceSize:  "m5.large",
				VolumeType:    "gp2",
				VolumeSize:    80,
				VolumeIOPS:    0,
				Notifications: ccx.Notifications{
					Enabled: false,
					Emails:  []string{"user@getccx.com"},
				},
				NetworkType:       "public",
				AvailabilityZones: []string{},
				FirewallRules:     []ccx.FirewallRule{},
			}).RunAndReturn(func(_ context.Context, _, _ ccx.Datastore) (*ccx.Datastore, error) {
			updated = true
			return updatedDatastore, nil
		})

		m.datastore.EXPECT().Delete(mock.Anything, "datastore-id").Return(nil).Once()

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			PreCheck: func() {
			},
			ProviderFactories: map[string]func() (*schema.Provider, error){
				"ccx": func() (*schema.Provider, error) {
					return p, nil
				},
			},
			Steps: []resource.TestStep{
				{
					Config: `
resource "ccx_datastore" "luna" {
  name           = "luna"
  size           = 1
  db_vendor      = "percona"
  type           = "replication" 
  tags           = ["new", "test"]
  cloud_provider = "aws"
  cloud_region   = "eu-north-1"
  instance_size  = "m5.large"
  volume_size    = 80
  volume_type    = "gp2"
  network_type   = "public"
}
`,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("ccx_datastore.luna", "id", "datastore-id"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "size", "1"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "db_vendor", "percona"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "db_version", "8"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "type", "Replication"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "tags.#", "2"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "tags.0", "new"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "tags.1", "test"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "cloud_provider", "aws"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "cloud_region", "eu-north-1"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "instance_size", "m5.large"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "volume_type", "gp2"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "volume_size", "80"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "volume_iops", "0"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "network_type", "public"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "network_vpc_uuid", ""),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "network_ha_enabled", "false"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "notifications_enabled", "false"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "notifications_emails.#", "1"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "notifications_emails.0", "user@getccx.com"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "maintenance_day_of_week", "1"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "maintenance_start_hour", "0"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "maintenance_end_hour", "2"),
					),
				},
				{
					Config: `
resource "ccx_datastore" "luna" {
  name           = "luna"
  size           = 2
  db_vendor      = "percona"
  type           = "replication"
  tags           = ["new", "test"]
  cloud_provider = "aws"
  cloud_region   = "eu-north-1"
  instance_size  = "m5.large"
  volume_size    = 80
  volume_type    = "gp2"
  network_type   = "public"
}
`,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("ccx_datastore.luna", "id", "datastore-id"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "size", "2"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "db_vendor", "percona"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "db_version", "8"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "type", "Replication"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "tags.#", "2"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "tags.0", "new"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "tags.1", "test"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "cloud_provider", "aws"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "cloud_region", "eu-north-1"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "instance_size", "m5.large"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "volume_type", "gp2"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "volume_size", "80"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "volume_iops", "0"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "network_type", "public"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "network_vpc_uuid", ""),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "network_ha_enabled", "false"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "notifications_enabled", "false"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "notifications_emails.#", "1"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "notifications_emails.0", "user@getccx.com"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "maintenance_day_of_week", "1"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "maintenance_start_hour", "0"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "maintenance_end_hour", "2"),
					),
				},
			},
		})

		m.AssertExpectations(t)
	})
}
