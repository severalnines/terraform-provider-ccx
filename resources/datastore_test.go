package resources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
	"github.com/stretchr/testify/mock"
)

func expectDefaultContent(m mockServices) {
	m.content.EXPECT().InstanceSizes(mock.Anything).Return(map[string][]ccx.InstanceSize{
		"aws": {
			{Code: "small", Type: "m5.large"},
			{Code: "medium", Type: "m5.xlarge"},
		},
	}, nil)

	m.content.EXPECT().DBVendors(mock.Anything).Return([]ccx.DBVendorInfo{
		{
			Name: "MariaDB",
			Code: "mariadb",
			Types: []ccx.DBVendorInfoType{
				{Name: "Multi-master", Code: "galera"},
				{Name: "Master/replicas", Code: "replication"},
			},
			DefaultVersion: "11.4",
			Versions:       []string{"10.11", "11.4"},
			NumNodes:       []int{1, 2, 3},
		},
		{
			Name: "MySQL",
			Code: "percona",
			Types: []ccx.DBVendorInfoType{
				{Name: "Multi-master", Code: "galera"},
				{Name: "Master/replicas", Code: "replication"},
			},
			DefaultVersion: "8",
			Versions:       []string{"8"},
			NumNodes:       []int{1, 2, 3},
		},
		{
			Name: "PostgreSQL",
			Code: "postgres",
			Types: []ccx.DBVendorInfoType{
				{Name: "Streaming Replication", Code: "postgres_streaming"},
			},
			DefaultVersion: "16",
			Versions:       []string{"14", "15", "16"},
			NumNodes:       []int{1, 2, 3},
		},
		{
			Name:           "Microsoft SQL Server",
			Code:           "microsoft",
			DefaultVersion: "2022",
			Versions:       []string{"2019", "2022"},
			Types: []ccx.DBVendorInfoType{
				{Name: "Single server", Code: "mssql_single"},
				{Name: "Always On (async commit mode)", Code: "mssql_ao_async"},
			},
			NumNodes: []int{1, 2},
		},
	}, nil)

	m.content.EXPECT().VolumeTypes(mock.Anything, "aws").Return([]string{"gp2"}, nil)

}

func TestDatastore_Create(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		m, p := mockProvider(t)

		expectDefaultContent(m)

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
			AvailabilityZones: nil,
			FirewallRules:     []ccx.FirewallRule{},
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
			PrimaryUrl: "00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			PrimaryDsn: "postgres://user:secret@00000000-0000-0000-0000-000000000001.app.mydbservice.net:5432",
			ReplicaUrl: "replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			ReplicaDsn: "postgres://user:secret@replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net:5432",
			Username:   "user",
			Password:   "secret",
			DbName:     "mydb",
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
			PrimaryUrl: "00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			PrimaryDsn: "postgres://user:secret@00000000-0000-0000-0000-000000000001.app.mydbservice.net:5432",
			ReplicaUrl: "replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			ReplicaDsn: "postgres://user:secret@replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net:5432",
			Username:   "user",
			Password:   "secret",
			DbName:     "mydb",
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
						resource.TestCheckResourceAttr("ccx_datastore.luna", "network_vpc_uuid", ""),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "network_ha_enabled", "false"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "notifications_enabled", "false"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "notifications_emails.#", "1"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "notifications_emails.0", "user@getccx.com"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "maintenance_day_of_week", "1"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "maintenance_start_hour", "0"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "maintenance_end_hour", "2"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "primary_url", "00000000-0000-0000-0000-000000000001.app.mydbservice.net"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "primary_dsn", "postgres://user:secret@00000000-0000-0000-0000-000000000001.app.mydbservice.net:5432"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "replica_url", "replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "replica_dsn", "postgres://user:secret@replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net:5432"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "username", "user"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "password", "secret"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "dbname", "mydb"),
					),
				},
			},
		})

		m.AssertExpectations(t)
	})

	t.Run("basic: update name", func(t *testing.T) {
		m, p := mockProvider(t)

		expectDefaultContent(m)

		create := ccx.Datastore{
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
			AvailabilityZones: nil,
			FirewallRules:     []ccx.FirewallRule{},
			Notifications: ccx.Notifications{
				Enabled: false,
				Emails:  []string{},
			},
		}

		created := create

		created.ID = "datastore-1"
		created.DBVersion = "1.2"
		created.Tags = []string{"new", "test", "tag1", "tag2"}
		created.VolumeIOPS = 0
		created.MaintenanceSettings = &ccx.MaintenanceSettings{
			DayOfWeek: 1,
			StartHour: 0,
			EndHour:   2,
		}
		created.PrimaryUrl = "datastore-1.app.mydbservice.net"
		created.PrimaryDsn = "postgres://user:secret@datastore-1.app.mydbservice.net:1234/mydb"
		created.ReplicaUrl = "replica.datastore-1.app.mydbservice.net"
		created.ReplicaDsn = "postgres://user:secret@replica.datastore-1.app.mydbservice.net:1234/mydb"
		created.Username = "user"
		created.Password = "secret"
		created.DbName = "mydb"

		updated := created
		updated.Name = "lunar"
		updated.Tags = create.Tags
		updated.PrimaryUrl = ""
		updated.PrimaryDsn = ""
		updated.ReplicaUrl = ""
		updated.ReplicaDsn = ""
		updated.Username = ""
		updated.Password = ""
		updated.DbName = ""

		latest := created

		m.datastore.EXPECT().Create(mock.Anything, create).Return(&created, nil)
		m.datastore.EXPECT().Update(mock.Anything, created, updated).RunAndReturn(func(_ context.Context, _ ccx.Datastore, n ccx.Datastore) (*ccx.Datastore, error) {
			latest = n
			return &latest, nil
		})

		m.datastore.EXPECT().Read(mock.Anything, "datastore-1").RunAndReturn(func(_ context.Context, _ string) (*ccx.Datastore, error) {
			return &latest, nil
		})

		m.datastore.EXPECT().Delete(mock.Anything, "datastore-1").Return(nil)

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
}
`,
					Check: checkTFDatastore("luna", created),
				},
				{
					Config: `
resource "ccx_datastore" "luna" {
  name           = "lunar"
  size           = 1
  db_vendor      = "postgres"
  tags           = ["new", "test"]
  cloud_provider = "aws"
  cloud_region   = "eu-north-1"
  instance_size  = "m5.large"
  volume_size    = 80
  volume_type    = "gp2"
}
`,
					Check: checkTFDatastore("luna", updated),
				},
			},
		})
	})

	t.Run("basic: update maintenance settings", func(t *testing.T) {
		m, p := mockProvider(t)

		expectDefaultContent(m)

		create := ccx.Datastore{
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
			AvailabilityZones: nil,
			FirewallRules:     []ccx.FirewallRule{},
			Notifications: ccx.Notifications{
				Enabled: false,
				Emails:  []string{},
			},
		}

		created := create

		created.ID = "datastore-1"
		created.DBVersion = "1.2"
		created.Tags = []string{"new", "test", "tag1", "tag2"}
		created.VolumeIOPS = 0
		created.MaintenanceSettings = &ccx.MaintenanceSettings{
			DayOfWeek: 1,
			StartHour: 0,
			EndHour:   2,
		}
		created.PrimaryUrl = "datastore-1.app.mydbservice.net"
		created.PrimaryDsn = "postgres://user:secret@datastore-1.app.mydbservice.net:1234/mydb"
		created.ReplicaUrl = "replica.datastore-1.app.mydbservice.net"
		created.ReplicaDsn = "postgres://user:secret@replica.datastore-1.app.mydbservice.net:1234/mydb"
		created.Username = "user"
		created.Password = "secret"
		created.DbName = "mydb"

		updated := created
		updated.Name = "luna"
		updated.Tags = create.Tags
		updated.PrimaryUrl = ""
		updated.PrimaryDsn = ""
		updated.ReplicaUrl = ""
		updated.ReplicaDsn = ""
		updated.Username = ""
		updated.Password = ""
		updated.DbName = ""
		updated.MaintenanceSettings = &ccx.MaintenanceSettings{
			DayOfWeek: 1,
			StartHour: 22,
			EndHour:   0,
		}

		latest := created

		m.datastore.EXPECT().Create(mock.Anything, create).Return(&created, nil)
		m.datastore.EXPECT().Update(mock.Anything, created, updated).RunAndReturn(func(_ context.Context, _ ccx.Datastore, n ccx.Datastore) (*ccx.Datastore, error) {
			latest = n
			return &latest, nil
		})

		m.datastore.EXPECT().Read(mock.Anything, "datastore-1").RunAndReturn(func(_ context.Context, _ string) (*ccx.Datastore, error) {
			return &latest, nil
		})

		m.datastore.EXPECT().Delete(mock.Anything, "datastore-1").Return(nil)

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
}
`,
					Check: checkTFDatastore("luna", created),
				},
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
  maintenance_day_of_week = 1
  maintenance_start_hour = 22
  maintenance_end_hour = 0
}
`,
					Check: checkTFDatastore("luna", updated),
				},
			},
		})
	})

	t.Run("scaling, with db_type replication", func(t *testing.T) {
		// the api payload to create the datastore uses the db_type field: replication
		// whereas the server returns the type field: Replication
		// this should not trigger a diff
		//
		// also testing the scaling part

		m, p := mockProvider(t)

		expectDefaultContent(m)

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
			PrimaryUrl: "00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			PrimaryDsn: "mysql://user:secret@00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306",
			ReplicaUrl: "replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			ReplicaDsn: "mysql://user:secret@replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306",
			Username:   "user",
			Password:   "secret",
			DbName:     "mydb",
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
			AvailabilityZones: nil,
			FirewallRules:     []ccx.FirewallRule{},
			MaintenanceSettings: &ccx.MaintenanceSettings{
				DayOfWeek: 1,
				StartHour: 0,
				EndHour:   2,
			},
			PrimaryUrl: "00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			PrimaryDsn: "mysql://user:secret@00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306",
			ReplicaUrl: "replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			ReplicaDsn: "mysql://user:secret@replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306",
			Username:   "user",
			Password:   "secret",
			DbName:     "mydb",
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
			AvailabilityZones: nil,
			FirewallRules:     []ccx.FirewallRule{},
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
				PrimaryUrl: "00000000-0000-0000-0000-000000000001.app.mydbservice.net",
				PrimaryDsn: "mysql://user:secret@00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306",
				ReplicaUrl: "replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net",
				ReplicaDsn: "mysql://user:secret@replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306",
				Username:   "user",
				Password:   "secret",
				DbName:     "mydb",
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
				MaintenanceSettings: &ccx.MaintenanceSettings{
					DayOfWeek: 1,
					StartHour: 0,
					EndHour:   2,
				},
				AvailabilityZones: nil,
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
						resource.TestCheckResourceAttr("ccx_datastore.luna", "network_vpc_uuid", ""),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "network_ha_enabled", "false"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "notifications_enabled", "false"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "notifications_emails.#", "1"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "notifications_emails.0", "user@getccx.com"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "maintenance_day_of_week", "1"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "maintenance_start_hour", "0"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "maintenance_end_hour", "2"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "primary_url", "00000000-0000-0000-0000-000000000001.app.mydbservice.net"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "primary_dsn", "mysql://user:secret@00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "replica_url", "replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "replica_dsn", "mysql://user:secret@replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "username", "user"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "password", "secret"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "dbname", "mydb"),
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
						resource.TestCheckResourceAttr("ccx_datastore.luna", "network_vpc_uuid", ""),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "network_ha_enabled", "false"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "notifications_enabled", "false"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "notifications_emails.#", "1"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "notifications_emails.0", "user@getccx.com"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "maintenance_day_of_week", "1"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "maintenance_start_hour", "0"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "maintenance_end_hour", "2"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "primary_url", "00000000-0000-0000-0000-000000000001.app.mydbservice.net"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "primary_dsn", "mysql://user:secret@00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "replica_url", "replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "replica_dsn", "mysql://user:secret@replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "username", "user"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "password", "secret"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "dbname", "mydb"),
					),
				},
			},
		})

		m.AssertExpectations(t)
	})

	t.Run("with parameter group", func(t *testing.T) {
		m, p := mockProvider(t)

		expectDefaultContent(m)

		pgCreated := ccx.ParameterGroup{
			ID:              "parameter-group-id",
			Name:            "asteroid",
			DatabaseVendor:  "mariadb",
			DatabaseVersion: "10.11",
			DatabaseType:    "replication",
			DbParameters: map[string]string{
				"max_connections": "100",
				"sql_mode":        "STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION",
			},
		}

		m.parameterGroup.EXPECT().Create(mock.Anything, ccx.ParameterGroup{
			Name:            "asteroid",
			DatabaseVendor:  "mariadb",
			DatabaseVersion: "10.11",
			DatabaseType:    "replication",
			DbParameters: map[string]string{
				"max_connections": "100",
				"sql_mode":        "STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION",
			},
		}).Return(&pgCreated, nil)

		m.parameterGroup.EXPECT().Read(mock.Anything, "parameter-group-id").Return(&pgCreated, nil)

		m.parameterGroup.EXPECT().Delete(mock.Anything, "parameter-group-id").Return(nil)

		m.datastore.EXPECT().Create(mock.Anything, ccx.Datastore{
			Name:              "luna",
			Size:              1,
			DBVendor:          "mariadb",
			Type:              "replication",
			Tags:              []string{"new", "test"},
			CloudProvider:     "aws",
			CloudRegion:       "eu-north-1",
			InstanceSize:      "m5.large",
			VolumeType:        "gp2",
			VolumeSize:        80,
			AvailabilityZones: nil,
			ParameterGroupID:  "parameter-group-id",
			FirewallRules:     []ccx.FirewallRule{},
			Notifications: ccx.Notifications{
				Enabled: false,
				Emails:  []string{},
			},
		}).Return(&ccx.Datastore{
			ID:               "datastore-id",
			Name:             "luna",
			Size:             1,
			DBVendor:         "mariadb",
			DBVersion:        "10.11",
			Type:             "replication",
			Tags:             []string{"new", "test", "mariadb", "10.11", "replication", "aws", "eu-north-1"},
			CloudProvider:    "aws",
			CloudRegion:      "eu-north-1",
			InstanceSize:     "m5.large",
			VolumeType:       "gp2",
			VolumeSize:       80,
			VolumeIOPS:       0,
			ParameterGroupID: "parameter-group-id",
			Notifications: ccx.Notifications{
				Enabled: false,
				Emails:  []string{"user@getccx.com"},
			},
			MaintenanceSettings: &ccx.MaintenanceSettings{
				DayOfWeek: 1,
				StartHour: 0,
				EndHour:   2,
			},
			PrimaryUrl: "00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			PrimaryDsn: "mysql://user:secret@00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306",
			ReplicaUrl: "replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			ReplicaDsn: "mysql://user:secret@replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306",
			Username:   "user",
			Password:   "secret",
			DbName:     "mydb",
		}, nil)

		m.datastore.EXPECT().Read(mock.Anything, "datastore-id").Return(&ccx.Datastore{
			ID:               "datastore-id",
			Name:             "luna",
			Size:             1,
			DBVendor:         "mariadb",
			DBVersion:        "10.11",
			Type:             "replication",
			Tags:             []string{"new", "test", "mariadb", "10.11", "replication", "aws", "eu-north-1"},
			CloudProvider:    "aws",
			CloudRegion:      "eu-north-1",
			InstanceSize:     "m5.large",
			VolumeType:       "gp2",
			VolumeSize:       80,
			VolumeIOPS:       0,
			ParameterGroupID: "parameter-group-id",
			Notifications: ccx.Notifications{
				Enabled: false,
				Emails:  []string{"user@getccx.com"},
			},
			MaintenanceSettings: &ccx.MaintenanceSettings{
				DayOfWeek: 1,
				StartHour: 0,
				EndHour:   2,
			},
			PrimaryUrl: "00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			PrimaryDsn: "mysql://user:secret@00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306",
			ReplicaUrl: "replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			ReplicaDsn: "mysql://user:secret@replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306",
			Username:   "user",
			Password:   "secret",
			DbName:     "mydb",
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
resource "ccx_parameter_group" "asteroid" {
    name = "asteroid"
    database_vendor = "mariadb"
    database_version = "10.11"
    database_type = "replication"

    parameters = {
      max_connections = 100
      sql_mode = "STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION"
    }
}

resource "ccx_datastore" "luna" {
  name            = "luna"
  size            = 1
  db_vendor       = "mariadb"
  tags            = ["new", "test"]
  cloud_provider  = "aws"
  cloud_region    = "eu-north-1"
  instance_size   = "m5.large"
  volume_size     = 80
  volume_type     = "gp2"
  parameter_group = ccx_parameter_group.asteroid.id
}
`,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("ccx_datastore.luna", "id", "datastore-id"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "size", "1"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "db_vendor", "mariadb"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "db_version", "10.11"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "type", "replication"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "tags.#", "2"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "tags.0", "new"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "tags.1", "test"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "cloud_provider", "aws"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "cloud_region", "eu-north-1"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "instance_size", "m5.large"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "volume_type", "gp2"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "volume_size", "80"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "volume_iops", "0"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "network_vpc_uuid", ""),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "network_ha_enabled", "false"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "notifications_enabled", "false"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "notifications_emails.#", "1"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "notifications_emails.0", "user@getccx.com"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "maintenance_day_of_week", "1"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "maintenance_start_hour", "0"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "maintenance_end_hour", "2"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "parameter_group", "parameter-group-id"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "primary_url", "00000000-0000-0000-0000-000000000001.app.mydbservice.net"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "primary_dsn", "mysql://user:secret@00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "replica_url", "replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "replica_dsn", "mysql://user:secret@replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "username", "user"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "password", "secret"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "dbname", "mydb"),
					),
				},
			},
		})

		m.AssertExpectations(t)
	})

	t.Run("basic with firewalls", func(t *testing.T) {
		m, p := mockProvider(t)

		expectDefaultContent(m)

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
			AvailabilityZones: nil,
			FirewallRules: []ccx.FirewallRule{
				{
					Source:      "2.2.2.0/24",
					Description: "One",
				},
				{
					Source:      "2.2.2.1/32",
					Description: "Two",
				},
			},
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
			PrimaryUrl: "00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			PrimaryDsn: "postgres://user:secret@00000000-0000-0000-0000-000000000001.app.mydbservice.net:5432",
			ReplicaUrl: "replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			ReplicaDsn: "postgres://user:secret@replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net:5432",
			Username:   "user",
			Password:   "secret",
			DbName:     "mydb",
		}, nil)

		m.datastore.EXPECT().SetFirewallRules(mock.Anything, "datastore-id", []ccx.FirewallRule{
			{
				Source:      "2.2.2.0/24",
				Description: "One",
			},
			{
				Source:      "2.2.2.1/32",
				Description: "Two",
			},
		}).Return(nil)

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
			FirewallRules: []ccx.FirewallRule{
				{
					Source:      "2.2.2.1/32",
					Description: "Two",
				},
				{
					Source:      "2.2.2.0/24",
					Description: "One",
				},
			},
			PrimaryUrl: "00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			PrimaryDsn: "postgres://user:secret@00000000-0000-0000-0000-000000000001.app.mydbservice.net:5432",
			ReplicaUrl: "replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			ReplicaDsn: "postgres://user:secret@replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net:5432",
			Username:   "user",
			Password:   "secret",
			DbName:     "mydb",
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

  firewall {
    source = "2.2.2.1/32"
    description = "Two"
  }

  firewall {
    source = "2.2.2.0/24"
    description = "One"
  }
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

						resource.TestCheckResourceAttr("ccx_datastore.luna", "network_vpc_uuid", ""),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "network_ha_enabled", "false"),

						resource.TestCheckResourceAttr("ccx_datastore.luna", "notifications_enabled", "false"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "notifications_emails.#", "1"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "notifications_emails.0", "user@getccx.com"),

						resource.TestCheckResourceAttr("ccx_datastore.luna", "maintenance_day_of_week", "1"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "maintenance_start_hour", "0"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "maintenance_end_hour", "2"),

						resource.TestCheckResourceAttr("ccx_datastore.luna", "firewall.#", "2"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "firewall.0.source", "2.2.2.0/24"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "firewall.0.description", "One"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "firewall.1.source", "2.2.2.1/32"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "firewall.1.description", "Two"),

						resource.TestCheckResourceAttr("ccx_datastore.luna", "primary_url", "00000000-0000-0000-0000-000000000001.app.mydbservice.net"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "primary_dsn", "postgres://user:secret@00000000-0000-0000-0000-000000000001.app.mydbservice.net:5432"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "replica_url", "replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "replica_dsn", "postgres://user:secret@replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net:5432"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "username", "user"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "password", "secret"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "dbname", "mydb"),
					),
				},
			},
		})

		m.AssertExpectations(t)
	})

	t.Run("basic with vendor alias mysql", func(t *testing.T) {
		m, p := mockProvider(t)

		expectDefaultContent(m)

		m.datastore.EXPECT().Create(mock.Anything, ccx.Datastore{
			Name:              "luna",
			Size:              1,
			DBVendor:          "percona",
			DBVersion:         "",
			Type:              "replication",
			Tags:              []string{"new", "test"},
			CloudProvider:     "aws",
			CloudRegion:       "eu-north-1",
			InstanceSize:      "m5.large",
			VolumeType:        "gp2",
			VolumeSize:        80,
			AvailabilityZones: nil,
			FirewallRules:     []ccx.FirewallRule{},
			Notifications: ccx.Notifications{
				Enabled: false,
				Emails:  []string{},
			},
		}).Return(&ccx.Datastore{
			ID:            "datastore-id",
			Name:          "luna",
			Size:          1,
			DBVendor:      "percona",
			DBVersion:     "8",
			Type:          "replication",
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
			PrimaryUrl: "00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			PrimaryDsn: "mysql://user:secret@00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306",
			ReplicaUrl: "replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			ReplicaDsn: "mysql://user:secret@replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306",
			Username:   "user",
			Password:   "secret",
			DbName:     "mydb",
		}, nil)

		m.datastore.EXPECT().Read(mock.Anything, "datastore-id").Return(&ccx.Datastore{
			ID:            "datastore-id",
			Name:          "luna",
			Size:          1,
			DBVendor:      "percona",
			DBVersion:     "8",
			Type:          "replication",
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
			PrimaryUrl: "00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			PrimaryDsn: "mysql://user:secret@00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306",
			ReplicaUrl: "replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			ReplicaDsn: "mysql://user:secret@replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306",
			Username:   "user",
			Password:   "secret",
			DbName:     "mydb",
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
  db_vendor      = "mysql"
  tags           = ["new", "test"]
  cloud_provider = "aws"
  cloud_region   = "eu-north-1"
  instance_size  = "m5.large"
  volume_size    = 80
  volume_type    = "gp2"
}
`,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("ccx_datastore.luna", "id", "datastore-id"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "size", "1"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "db_vendor", "percona"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "db_version", "8"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "type", "replication"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "tags.#", "2"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "tags.0", "new"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "tags.1", "test"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "cloud_provider", "aws"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "cloud_region", "eu-north-1"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "instance_size", "m5.large"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "volume_type", "gp2"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "volume_size", "80"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "volume_iops", "0"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "network_vpc_uuid", ""),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "network_ha_enabled", "false"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "notifications_enabled", "false"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "notifications_emails.#", "1"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "notifications_emails.0", "user@getccx.com"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "maintenance_day_of_week", "1"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "maintenance_start_hour", "0"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "maintenance_end_hour", "2"),

						resource.TestCheckResourceAttr("ccx_datastore.luna", "primary_url", "00000000-0000-0000-0000-000000000001.app.mydbservice.net"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "primary_dsn", "mysql://user:secret@00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "replica_url", "replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "replica_dsn", "mysql://user:secret@replica.00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "username", "user"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "password", "secret"),
						resource.TestCheckResourceAttr("ccx_datastore.luna", "dbname", "mydb"),
					),
				},
			},
		})

		m.AssertExpectations(t)
	})
}
