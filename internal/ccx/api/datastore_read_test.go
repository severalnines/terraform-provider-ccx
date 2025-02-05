package api

import (
	"testing"

	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_dsn(t *testing.T) {
	tests := []struct {
		name     string
		vendor   string
		host     string
		port     int
		username string
		password string
		dbname   string
		want     string
	}{
		{
			name:     "mysql",
			vendor:   "mysql",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			port:     3306,
			username: "user",
			password: "pass",
			dbname:   "testdb",
			want:     "mysql://user:pass@00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306/testdb",
		},
		{
			name:     "mysql with host containing port",
			vendor:   "mysql",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306",
			port:     1234,
			username: "user",
			password: "pass",
			dbname:   "testdb",
			want:     "mysql://user:pass@00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306/testdb",
		},
		{
			name:     "percona",
			vendor:   "percona",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			port:     3306,
			username: "user",
			password: "pass",
			dbname:   "testdb",
			want:     "mysql://user:pass@00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306/testdb",
		},
		{
			name:     "percona with host containing port",
			vendor:   "percona",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306",
			port:     1234,
			username: "user",
			password: "pass",
			dbname:   "testdb",
			want:     "mysql://user:pass@00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306/testdb",
		},
		{
			name:     "mariadb",
			vendor:   "mariadb",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			port:     3306,
			username: "user",
			password: "pass",
			dbname:   "testdb",
			want:     "mysql://user:pass@00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306/testdb",
		},
		{
			name:     "mariadb with host containing port",
			vendor:   "mariadb",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306",
			port:     1234,
			username: "user",
			password: "pass",
			dbname:   "testdb",
			want:     "mysql://user:pass@00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306/testdb",
		},
		{
			name:     "postgres",
			vendor:   "postgres",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			port:     5432,
			username: "user",
			password: "pass",
			dbname:   "testdb",
			want:     "postgres://user:pass@00000000-0000-0000-0000-000000000001.app.mydbservice.net:5432/testdb",
		},
		{
			name:     "postgres with host containing port",
			vendor:   "postgres",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net:5432",
			port:     9999,
			username: "user",
			password: "pass",
			dbname:   "testdb",
			want:     "postgres://user:pass@00000000-0000-0000-0000-000000000001.app.mydbservice.net:5432/testdb",
		},
		{
			name:     "redis",
			vendor:   "redis",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			port:     6379,
			username: "user",
			password: "pass",
			dbname:   "testdb",
			want:     "rediss://user:pass@00000000-0000-0000-0000-000000000001.app.mydbservice.net:6379/testdb",
		},
		{
			name:     "redis with host containing port",
			vendor:   "redis",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net:6379",
			port:     8888,
			username: "user",
			password: "pass",
			dbname:   "testdb",
			want:     "rediss://user:pass@00000000-0000-0000-0000-000000000001.app.mydbservice.net:6379/testdb",
		},
		{
			name:     "microsoft",
			vendor:   "microsoft",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			port:     1433,
			username: "user",
			password: "pass",
			dbname:   "testdb",
			want:     "Data Source=00000000-0000-0000-0000-000000000001.app.mydbservice.net:1433;User ID=user;Password=pass;Database=testdb",
		},
		{
			name:     "microsoft with host containing port",
			vendor:   "microsoft",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net:1433",
			port:     9999,
			username: "user",
			password: "pass",
			dbname:   "testdb",
			want:     "Data Source=00000000-0000-0000-0000-000000000001.app.mydbservice.net:1433;User ID=user;Password=pass;Database=testdb",
		},
		{
			name:     "unknown vendor",
			vendor:   "unknown",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net",
			port:     3306,
			username: "user",
			password: "pass",
			dbname:   "testdb",
			want:     "",
		},
		{
			name:     "unknown vendor with host containing port",
			vendor:   "unknown",
			host:     "00000000-0000-0000-0000-000000000001.app.mydbservice.net:3306",
			port:     1234,
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

func Test_getPortFromDatastore(t *testing.T) {
	tests := []struct {
		name    string
		c       ccx.Datastore
		want    int
		wantErr bool
	}{
		{
			name: "single node, success",
			c: ccx.Datastore{
				Hosts: []ccx.Host{
					{
						Port: 3306,
						Role: "primary",
					},
				},
			},
			want:    3306,
			wantErr: false,
		},
		{
			name: "multiple nodes, success",
			c: ccx.Datastore{
				Hosts: []ccx.Host{
					{
						Port: 3306,
						Role: "primary",
					},
				},
			},
			want:    3306,
			wantErr: false,
		},
		{
			name: "multiple nodes, no port in primary, return port of replica",
			c: ccx.Datastore{
				Hosts: []ccx.Host{
					{
						Port: 0,
						Role: "primary",
					},
				},
			},
			want:    3306,
			wantErr: false,
		},
		{
			name: "multiple nodes, no port in primary and port in 1 replica, return port of replica",
			c: ccx.Datastore{
				Hosts: []ccx.Host{
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
			},
			want:    3306,
			wantErr: false,
		},
		{
			name: "multiple nodes, no port at all, return error",
			c: ccx.Datastore{
				Hosts: []ccx.Host{
					{
						Port: 0,
						Role: "replica",
					},
					{
						Port: 0,
						Role: "primary",
					},
				},
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPortFromDatastore(tt.c)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoErrorf(t, err, "should not get error, getting: %v", err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
