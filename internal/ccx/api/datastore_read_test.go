package api

import (
	"context"
	"errors"
	"testing"

	"github.com/severalnines/terraform-provider-ccx/internal/ccx/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDsn(t *testing.T) {
	tests := []struct {
		name     string
		vendor   string
		host     string
		port     string
		username string
		password string
		dbname   string
		want     string
	}{
		{
			name:     "mysql",
			vendor:   "mysql",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			port:     "3306",
			username: "user",
			password: "pass",
			dbname:   "testdb",
			want:     "mysql://user:pass@00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306/testdb",
		},
		{
			name:     "mysql with host containing port",
			vendor:   "mysql",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306",
			port:     "1234",
			username: "user",
			password: "pass",
			dbname:   "testdb",
			want:     "mysql://user:pass@00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306/testdb",
		},
		{
			name:     "percona",
			vendor:   "percona",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			port:     "3306",
			username: "user",
			password: "pass",
			dbname:   "testdb",
			want:     "mysql://user:pass@00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306/testdb",
		},
		{
			name:     "percona with host containing port",
			vendor:   "percona",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306",
			port:     "1234",
			username: "user",
			password: "pass",
			dbname:   "testdb",
			want:     "mysql://user:pass@00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306/testdb",
		},
		{
			name:     "mariadb",
			vendor:   "mariadb",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			port:     "3306",
			username: "user",
			password: "pass",
			dbname:   "testdb",
			want:     "mysql://user:pass@00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306/testdb",
		},
		{
			name:     "mariadb with host containing port",
			vendor:   "mariadb",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306",
			port:     "1234",
			username: "user",
			password: "pass",
			dbname:   "testdb",
			want:     "mysql://user:pass@00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306/testdb",
		},
		{
			name:     "postgres",
			vendor:   "postgres",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			port:     "5432",
			username: "user",
			password: "pass",
			dbname:   "testdb",
			want:     "postgres://user:pass@00000000-0000-0000-0000-000000000001.app.mydbservice.net:5432/testdb",
		},
		{
			name:     "postgres with host containing port",
			vendor:   "postgres",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net:5432",
			port:     "9999",
			username: "user",
			password: "pass",
			dbname:   "testdb",
			want:     "postgres://user:pass@00000000-0000-0000-0000-000000000001.app.mydbservice.net:5432/testdb",
		},
		{
			name:     "redis",
			vendor:   "redis",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			port:     "6379",
			username: "user",
			password: "pass",
			dbname:   "testdb",
			want:     "rediss://user:pass@00000000-0000-0000-0000-000000000001.app.mydbservice.net:6379/testdb",
		},
		{
			name:     "redis with host containing port",
			vendor:   "redis",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net:6379",
			port:     "8888",
			username: "user",
			password: "pass",
			dbname:   "testdb",
			want:     "rediss://user:pass@00000000-0000-0000-0000-000000000001.app.mydbservice.net:6379/testdb",
		},
		{
			name:     "microsoft",
			vendor:   "microsoft",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			port:     "1433",
			username: "user",
			password: "pass",
			dbname:   "testdb",
			want:     "Data Source=00000000-0000-0000-0000-000000000001.app.mydbservice.net:1433;User ID=user;Password=pass;Database=testdb",
		},
		{
			name:     "microsoft with host containing port",
			vendor:   "microsoft",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net:1433",
			port:     "9999",
			username: "user",
			password: "pass",
			dbname:   "testdb",
			want:     "Data Source=00000000-0000-0000-0000-000000000001.app.mydbservice.net:1433;User ID=user;Password=pass;Database=testdb",
		},
		{
			name:     "unknown vendor",
			vendor:   "unknown",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			port:     "3306",
			username: "user",
			password: "pass",
			dbname:   "testdb",
			want:     "",
		},
		{
			name:     "unknown vendor with host containing port",
			vendor:   "unknown",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306",
			port:     "1234",
			username: "user",
			password: "pass",
			dbname:   "testdb",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dsn(tt.vendor, tt.host, tt.port, tt.username, tt.password, tt.dbname); got != tt.want {
				t.Errorf("dsn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockGetDatastoreNodesResponse(t *getDatastoreNodesResponse) func(rs *getDatastoreNodesResponse) bool {
	return func(rs *getDatastoreNodesResponse) bool {
		rs.DatabaseNodes = t.DatabaseNodes
		return true
	}
}

func TestDatastoreService_getPort(t *testing.T) {
	path := "/api/deployment/v2/data-stores/datastore-id/nodes"

	tests := []struct {
		name    string
		id      string
		mock    func(h *mocks.MockHttpClient)
		want    string
		wantErr bool
	}{
		{
			name: "single node, success",
			id:   "datastore-id",
			mock: func(h *mocks.MockHttpClient) {
				mocks.MockHttpClientExpectGet(h, path, mockGetDatastoreNodesResponse(&getDatastoreNodesResponse{
					DatabaseNodes: []struct {
						Port int    `json:"port"`
						Role string `json:"role"`
					}{
						{
							Port: 3306,
							Role: "primary",
						},
					},
				}), nil)
			},
			want:    "3306",
			wantErr: false,
		},
		{
			name: "multiple nodes, success",
			id:   "datastore-id",
			mock: func(h *mocks.MockHttpClient) {
				mocks.MockHttpClientExpectGet(h, path, mockGetDatastoreNodesResponse(&getDatastoreNodesResponse{
					DatabaseNodes: []struct {
						Port int    `json:"port"`
						Role string `json:"role"`
					}{
						{
							Port: 13306,
							Role: "replica",
						},
						{
							Port: 3306,
							Role: "primary",
						},
					},
				}), nil)
			},
			want:    "3306",
			wantErr: false,
		},
		{
			name: "multiple nodes, no port in primary, return port of replica",
			id:   "datastore-id",
			mock: func(h *mocks.MockHttpClient) {
				mocks.MockHttpClientExpectGet(h, path, mockGetDatastoreNodesResponse(&getDatastoreNodesResponse{
					DatabaseNodes: []struct {
						Port int    `json:"port"`
						Role string `json:"role"`
					}{
						{
							Port: 3306,
							Role: "replica",
						},
						{
							Port: 0,
							Role: "primary",
						},
					},
				}), nil)
			},
			want:    "3306",
			wantErr: false,
		},
		{
			name: "multiple nodes, no port in primary and port in 1 replica, return port of replica",
			id:   "datastore-id",
			mock: func(h *mocks.MockHttpClient) {
				mocks.MockHttpClientExpectGet(h, path, mockGetDatastoreNodesResponse(&getDatastoreNodesResponse{
					DatabaseNodes: []struct {
						Port int    `json:"port"`
						Role string `json:"role"`
					}{
						{
							Port: 0,
							Role: "primary",
						},
						{
							Port: 0,
							Role: "replica",
						},
						{
							Port: 3306,
							Role: "replica",
						},
					},
				}), nil)
			},
			want:    "3306",
			wantErr: false,
		},
		{
			name: "multiple nodes, no port at all, return error",
			id:   "datastore-id",
			mock: func(h *mocks.MockHttpClient) {
				mocks.MockHttpClientExpectGet(h, path, mockGetDatastoreNodesResponse(&getDatastoreNodesResponse{
					DatabaseNodes: []struct {
						Port int    `json:"port"`
						Role string `json:"role"`
					}{
						{
							Port: 0,
							Role: "replica",
						},
						{
							Port: 0,
							Role: "primary",
						},
					},
				}), nil)
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "request error, return error",
			id:   "datastore-id",
			mock: func(h *mocks.MockHttpClient) {
				mocks.MockHttpClientExpectGet(h, path, mockGetDatastoreNodesResponse(&getDatastoreNodesResponse{}), errors.New("request error"))
			},
			want:    "",
			wantErr: true,
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
				tt.mock(httpcli)
			}

			got, err := svc.getPort(ctx, tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoErrorf(t, err, "should not get error, getting: %v", err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
