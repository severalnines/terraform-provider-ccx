package ccx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDatastoreService_Update(t *testing.T) {
	tests := []struct {
		name    string
		old     Datastore
		next    Datastore
		mock    func(h *MockHTTPClient, j *MockJobsService, c *MockContentService)
		want    *Datastore
		wantErr bool
	}{
		{
			name: "resize, 1 -> 3 with 2 azs total",
			old: Datastore{
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
				Hosts: []Host{
					{
						ID:            "host-1",
						CreatedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						CloudProvider: "aws",
						AZ:            "eu-north-1a",
						InstanceType:  "m5.large",
						Role:          "primary",
						DiskSize:      80,
						DiskType:      "gp2",
						Region:        "eu-north-1",
					},
				},
				Notifications: Notifications{
					Enabled: false,
					Emails:  []string{"user@getcom"},
				},
				MaintenanceSettings: &MaintenanceSettings{
					DayOfWeek: 1,
					StartHour: 0,
					EndHour:   2,
				},
			},
			next: Datastore{
				ID:            "datastore-id",
				Name:          "luna",
				Size:          3,
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
				Hosts: []Host{
					{
						ID:            "host-1",
						CreatedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						CloudProvider: "aws",
						AZ:            "eu-north-1a",
						InstanceType:  "m5.large",
						Role:          "primary",
						DiskSize:      80,
						DiskType:      "gp2",
						Region:        "eu-north-1",
					},
				},
				Notifications: Notifications{
					Enabled: false,
					Emails:  []string{"user@getcom"},
				},
				AvailabilityZones: nil,
				FirewallRules:     []FirewallRule{},
			},
			mock: func(h *MockHTTPClient, j *MockJobsService, c *MockContentService) {
				c.EXPECT().AvailabilityZones(mock.Anything, "aws", "eu-north-1").
					Return([]string{"eu-north-1a", "eu-north-1b"}, nil)

				r := httptest.NewRecorder()
				r.Body.WriteString("{}")

				// First call: tag update
				h.EXPECT().Do(mock.Anything, http.MethodPatch, "/api/prov/api/v2/cluster/datastore-id", updateRequest{
					NewName:       "",
					NewVolumeSize: 0,
					Remove:        nil,
					Add:           nil,
					Notifications: nil,
					Maintenance:   nil,
					Tags:          []string{"new", "test"},
				}).Return(fakeHttpResponse(http.StatusOK, ""), nil)

				// Second call: resize update
				h.EXPECT().Do(mock.Anything, http.MethodPatch, "/api/prov/api/v2/cluster/datastore-id", updateRequest{
					NewName:       "",
					NewVolumeSize: 0,
					Remove:        nil,
					Notifications: nil,
					Maintenance:   nil,
					Add: &addHosts{
						Specs: []hostSpecs{
							{
								InstanceSize: "m5.large",
								AZ:           "eu-north-1b",
							},
							{
								InstanceSize: "m5.large",
								AZ:           "eu-north-1b",
							},
						},
					},
				}).Return(fakeHttpResponse(http.StatusOK, ""), nil)

				j.EXPECT().Await(mock.Anything, "datastore-id", AddNodeJob).Return(JobStatusFinished, nil)

				MockHTTPClientExpectGet(h, "/api/deployment/v3/data-stores/datastore-id", getDatastoreResponse{
					ID:            "datastore-id",
					CloudProvider: "aws",
					Region: struct {
						Code string `json:"code"`
					}{
						Code: "eu-north-1",
					},
					InstanceSize:     "m5.large",
					InstanceIOPS:     nil,
					DiskSize:         Uint64P(80),
					DiskType:         StringP("gp2"),
					DbVendor:         "percona",
					DbVersion:        "8",
					Name:             "luna",
					Status:           "STARTED",
					StatusText:       "There are no failed nodes, there are started nodes",
					Type:             "Replication",
					TypeName:         "Streaming Replication",
					Size:             3,
					SSLEnabled:       true,
					HighAvailability: false,
					Tags: []string{
						"percona",
						"8",
						"replication",
						"aws",
						"eu-north-1",
					},
					AZS: []string{
						"eu-north-1a",
						"eu-north-1b",
					},
					Notifications: Notifications{
						Enabled: false,
						Emails:  []string{"user@getcom"},
					},
					CreatedAt:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:  time.Date(2021, 1, 1, 1, 0, 0, 0, time.UTC),
					PrimaryUrl: "cluster.mydbservice.net",
					ReplicaUrl: "replica.mydbservice.net",

					DbAccount: struct {
						Username   string `json:"database_username"`
						Password   string `json:"database_password"`
						Host       string `json:"database_host"`
						Database   string `json:"database_database"`
						Privileges string `json:"database_privileges"`
					}{Username: "ccx", Password: "top-secret", Database: "mydb"},
				}, nil)

				MockHTTPClientExpectGet(h, "/api/firewall/api/v1/firewalls/datastore-id", getFirewallsResponse{}, nil)

				MockHTTPClientExpectGet(h, "/api/deployment/v2/data-stores/datastore-id/nodes", getHostsResponse{
					UUID: "datastore-id",
					Hosts: []struct {
						ID            string    `json:"host_uuid"`
						CreatedAt     time.Time `json:"created_at"`
						CloudProvider string    `json:"cloud_provider"`
						AZ            string    `json:"host_az"`
						InstanceType  string    `json:"instance_type"`
						DiskType      string    `json:"disk_type"`
						DiskSize      uint64    `json:"disk_size"`
						Role          string    `json:"role"`
						Port          int       `json:"port"`
						Region        struct {
							Code string `json:"code"`
						} `json:"region"`
					}{
						{
							ID:            "host-1",
							CreatedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
							CloudProvider: "aws",
							AZ:            "eu-north-1a",
							InstanceType:  "m5.large",
							DiskType:      "gp2",
							DiskSize:      80,
							Role:          "primary",
							Port:          3306,
							Region: struct {
								Code string `json:"code"`
							}{
								Code: "eu-north-1",
							},
						},
						{
							ID:            "host-2",
							CreatedAt:     time.Date(2021, 1, 1, 1, 0, 0, 0, time.UTC),
							CloudProvider: "aws",
							AZ:            "eu-north-1b",
							InstanceType:  "m5.large",
							DiskType:      "gp2",
							DiskSize:      80,
							Role:          "replica",
							Port:          3306,
							Region: struct {
								Code string `json:"code"`
							}{
								Code: "eu-north-1",
							},
						},
						{
							ID:            "host-3",
							CreatedAt:     time.Date(2021, 1, 1, 1, 0, 0, 0, time.UTC),
							CloudProvider: "aws",
							AZ:            "eu-north-1b",
							InstanceType:  "m5.large",
							DiskType:      "gp2",
							DiskSize:      80,
							Role:          "replica",
							Port:          3306,
							Region: struct {
								Code string `json:"code"`
							}{
								Code: "eu-north-1",
							},
						},
					},
				}, nil)
			},
			want: &Datastore{
				ID:            "datastore-id",
				Name:          "luna",
				Size:          3,
				DBVendor:      "percona",
				DBVersion:     "8",
				Type:          "Replication",
				Tags:          []string{"percona", "8", "replication", "aws", "eu-north-1"},
				CloudProvider: "aws",
				CloudRegion:   "eu-north-1",
				InstanceSize:  "m5.large",
				VolumeType:    "gp2",
				VolumeSize:    80,
				VolumeIOPS:    0,
				Hosts: []Host{
					{
						ID:            "host-1",
						CreatedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						CloudProvider: "aws",
						AZ:            "eu-north-1a",
						InstanceType:  "m5.large",
						Role:          "primary",
						DiskSize:      80,
						DiskType:      "gp2",
						Region:        "eu-north-1",
						Port:          3306,
					},
					{
						ID:            "host-2",
						CreatedAt:     time.Date(2021, 1, 1, 1, 0, 0, 0, time.UTC),
						CloudProvider: "aws",
						AZ:            "eu-north-1b",
						InstanceType:  "m5.large",
						Role:          "replica",
						DiskSize:      80,
						DiskType:      "gp2",
						Region:        "eu-north-1",
						Port:          3306,
					},
					{
						ID:            "host-3",
						CreatedAt:     time.Date(2021, 1, 1, 1, 0, 0, 0, time.UTC),
						CloudProvider: "aws",
						AZ:            "eu-north-1b",
						InstanceType:  "m5.large",
						Role:          "replica",
						DiskSize:      80,
						DiskType:      "gp2",
						Region:        "eu-north-1",
						Port:          3306,
					},
				},
				Notifications: Notifications{
					Enabled: false,
					Emails:  []string{"user@getcom"},
				},
				AvailabilityZones: []string{
					"eu-north-1a",
					"eu-north-1b",
				},
				FirewallRules: []FirewallRule{},
				PrimaryUrl:    "cluster.mydbservice.net",
				PrimaryDsn:    "mysql://ccx:top-secret@cluster.mydbservice.net:3306/mydb",
				ReplicaUrl:    "replica.mydbservice.net",
				ReplicaDsn:    "mysql://ccx:top-secret@replica.mydbservice.net:3306/mydb",
				Username:      "ccx",
				Password:      "top-secret",
				DbName:        "mydb",
			},
			wantErr: false,
		},
		{
			name: "update instance size only",
			old: Datastore{
				ID:            "datastore-id",
				Name:          "luna",
				Size:          1,
				DBVendor:      "percona",
				DBVersion:     "8",
				Type:          "Replication",
				Tags:          []string{"tag1", "tag2"},
				CloudProvider: "aws",
				CloudRegion:   "eu-north-1",
				InstanceSize:  "m5.large",
				VolumeType:    "gp2",
				VolumeSize:    80,
				VolumeIOPS:    0,
				Hosts: []Host{
					{
						ID:            "host-1",
						CreatedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						CloudProvider: "aws",
						AZ:            "eu-north-1a",
						InstanceType:  "m5.large",
						Role:          "primary",
						DiskSize:      80,
						DiskType:      "gp2",
						Region:        "eu-north-1",
					},
				},
				Notifications: Notifications{
					Enabled: false,
					Emails:  []string{"user@getcom"},
				},
			},
			next: Datastore{
				ID:            "datastore-id",
				Name:          "luna",
				Size:          1,
				DBVendor:      "percona",
				DBVersion:     "8",
				Type:          "Replication",
				Tags:          []string{"tag1", "tag2"},
				CloudProvider: "aws",
				CloudRegion:   "eu-north-1",
				InstanceSize:  "m5.xlarge", // Changed from m5.large to m5.xlarge
				VolumeType:    "gp2",
				VolumeSize:    80,
				VolumeIOPS:    0,
				Hosts: []Host{
					{
						ID:            "host-1",
						CreatedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						CloudProvider: "aws",
						AZ:            "eu-north-1a",
						InstanceType:  "m5.large",
						Role:          "primary",
						DiskSize:      80,
						DiskType:      "gp2",
						Region:        "eu-north-1",
					},
				},
				Notifications: Notifications{
					Enabled: false,
					Emails:  []string{"user@getcom"},
				},
			},
			mock: func(h *MockHTTPClient, j *MockJobsService, c *MockContentService) {
				r := httptest.NewRecorder()
				r.Body.WriteString("{}")

				// Expect instance size update call
				h.EXPECT().Do(mock.Anything, http.MethodPatch, "/api/prov/api/v2/cluster/datastore-id", updateRequest{
					NewName:         "",
					NewVolumeSize:   0,
					NewInstanceSize: "m5.xlarge",
					Remove:          nil,
					Add:             nil,
					Notifications:   nil,
					Maintenance:     nil,
					Tags:            nil,
				}).Return(fakeHttpResponse(http.StatusOK, ""), nil)

				// Mock the Read call that happens after update
				MockHTTPClientExpectGet(h, "/api/deployment/v3/data-stores/datastore-id", getDatastoreResponse{
					ID:            "datastore-id",
					CloudProvider: "aws",
					Region: struct {
						Code string `json:"code"`
					}{
						Code: "eu-north-1",
					},
					InstanceSize:     "m5.xlarge",
					InstanceIOPS:     nil,
					DiskSize:         Uint64P(80),
					DiskType:         StringP("gp2"),
					DbVendor:         "percona",
					DbVersion:        "8",
					Name:             "luna",
					Status:           "STARTED",
					StatusText:       "There are no failed nodes, there are started nodes",
					Type:             "Replication",
					TypeName:         "Streaming Replication",
					Size:             1,
					SSLEnabled:       true,
					HighAvailability: false,
					Tags:             []string{"tag1", "tag2"},
					AZS:              []string{"eu-north-1a"},
					Notifications: Notifications{
						Enabled: false,
						Emails:  []string{"user@getcom"},
					},
					CreatedAt:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:  time.Date(2021, 1, 1, 1, 0, 0, 0, time.UTC),
					PrimaryUrl: "cluster.mydbservice.net",
					ReplicaUrl: "replica.mydbservice.net",
					DbAccount: struct {
						Username   string `json:"database_username"`
						Password   string `json:"database_password"`
						Host       string `json:"database_host"`
						Database   string `json:"database_database"`
						Privileges string `json:"database_privileges"`
					}{Username: "ccx", Password: "top-secret", Database: "mydb"},
				}, nil)

				MockHTTPClientExpectGet(h, "/api/firewall/api/v1/firewalls/datastore-id", getFirewallsResponse{}, nil)

				MockHTTPClientExpectGet(h, "/api/deployment/v2/data-stores/datastore-id/nodes", getHostsResponse{
					UUID: "datastore-id",
					Hosts: []struct {
						ID            string    `json:"host_uuid"`
						CreatedAt     time.Time `json:"created_at"`
						CloudProvider string    `json:"cloud_provider"`
						AZ            string    `json:"host_az"`
						InstanceType  string    `json:"instance_type"`
						DiskType      string    `json:"disk_type"`
						DiskSize      uint64    `json:"disk_size"`
						Role          string    `json:"role"`
						Port          int       `json:"port"`
						Region        struct {
							Code string `json:"code"`
						} `json:"region"`
					}{
						{
							ID:            "host-1",
							CreatedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
							CloudProvider: "aws",
							AZ:            "eu-north-1a",
							InstanceType:  "m5.xlarge", // Updated instance type
							Role:          "primary",
							DiskSize:      80,
							DiskType:      "gp2",
							Port:          3306,
							Region: struct {
								Code string `json:"code"`
							}{
								Code: "eu-north-1",
							},
						},
					},
				}, nil)
			},
			want: &Datastore{
				ID:            "datastore-id",
				Name:          "luna",
				Size:          1,
				DBVendor:      "percona",
				DBVersion:     "8",
				Type:          "Replication",
				Tags:          []string{"tag1", "tag2"},
				CloudProvider: "aws",
				CloudRegion:   "eu-north-1",
				InstanceSize:  "m5.xlarge",
				VolumeType:    "gp2",
				VolumeSize:    80,
				VolumeIOPS:    0,
				Hosts: []Host{
					{
						ID:            "host-1",
						CreatedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						CloudProvider: "aws",
						AZ:            "eu-north-1a",
						InstanceType:  "m5.xlarge",
						Role:          "primary",
						DiskSize:      80,
						DiskType:      "gp2",
						Region:        "eu-north-1",
						Port:          3306,
					},
				},
				Notifications: Notifications{
					Enabled: false,
					Emails:  []string{"user@getcom"},
				},
				AvailabilityZones: []string{
					"eu-north-1a",
				},
				FirewallRules: []FirewallRule{},
				PrimaryUrl:    "cluster.mydbservice.net",
				PrimaryDsn:    "mysql://ccx:top-secret@cluster.mydbservice.net:3306/mydb",
				ReplicaUrl:    "replica.mydbservice.net",
				ReplicaDsn:    "mysql://ccx:top-secret@replica.mydbservice.net:3306/mydb",
				Username:      "ccx",
				Password:      "top-secret",
				DbName:        "mydb",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			httpcli := NewMockHTTPClient(t)
			jobsSvc := NewMockJobsService(t)
			contentSvc := NewMockContentService(t)

			svc := &DatastoresClient{
				client:     httpcli,
				jobs:       jobsSvc,
				contentSvc: contentSvc,
			}

			if tt.mock != nil {
				tt.mock(httpcli, jobsSvc, contentSvc)
			}

			got, err := svc.Update(ctx, tt.old, tt.next)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoErrorf(t, err, "should not get error, getting: %v", err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
