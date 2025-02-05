package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx/mocks"
	"github.com/severalnines/terraform-provider-ccx/internal/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDatastoreService_Update(t *testing.T) {
	tests := []struct {
		name    string
		old     ccx.Datastore
		next    ccx.Datastore
		mock    func(h *mocks.MockHttpClient, j *mocks.MockJobService, c *mocks.MockContentService)
		want    *ccx.Datastore
		wantErr bool
	}{
		{
			name: "resize, 1 -> 3 with 2 azs total",
			old: ccx.Datastore{
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
				Hosts: []ccx.Host{
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
			next: ccx.Datastore{
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
				Hosts: []ccx.Host{
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
				Notifications: ccx.Notifications{
					Enabled: false,
					Emails:  []string{"user@getccx.com"},
				},
				AvailabilityZones: nil,
				FirewallRules:     []ccx.FirewallRule{},
			},
			mock: func(h *mocks.MockHttpClient, j *mocks.MockJobService, c *mocks.MockContentService) {
				c.EXPECT().AvailabilityZones(mock.Anything, "aws", "eu-north-1").
					Return([]string{"eu-north-1a", "eu-north-1b"}, nil)

				r := httptest.NewRecorder()
				r.Body.WriteString("{}")

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

				j.EXPECT().Await(mock.Anything, "datastore-id", ccx.AddNodeJob).Return(ccx.JobStatusFinished, nil)

				mocks.MockHttpClientExpectGet(h, "/api/deployment/v3/data-stores/datastore-id", getDatastoreResponse{
					ID:            "datastore-id",
					CloudProvider: "aws",
					Region: struct {
						Code string `json:"code"`
					}{
						Code: "eu-north-1",
					},
					InstanceSize:     "m5.large",
					InstanceIOPS:     nil,
					DiskSize:         lib.Uint64P(80),
					DiskType:         lib.StringP("gp2"),
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
					Notifications: ccx.Notifications{
						Enabled: false,
						Emails:  []string{"user@getccx.com"},
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

				mocks.MockHttpClientExpectGet(h, "/api/firewall/api/v1/firewalls/datastore-id", getFirewallsResponse{}, nil)

				mocks.MockHttpClientExpectGet(h, "/api/deployment/v2/data-stores/datastore-id/nodes", getHostsResponse{
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
			want: &ccx.Datastore{
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
				Hosts: []ccx.Host{
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
				Notifications: ccx.Notifications{
					Enabled: false,
					Emails:  []string{"user@getccx.com"},
				},
				AvailabilityZones: []string{
					"eu-north-1a",
					"eu-north-1b",
				},
				FirewallRules: []ccx.FirewallRule{},
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

			httpcli := mocks.NewMockHttpClient(t)
			jobsSvc := mocks.NewMockJobService(t)
			contentSvc := mocks.NewMockContentService(t)

			svc := &DatastoreService{
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
