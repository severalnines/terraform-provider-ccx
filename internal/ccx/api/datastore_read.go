package api

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
	"github.com/severalnines/terraform-provider-ccx/internal/lib"
)

type getDatastoreResponse struct {
	ID string `json:"uuid"`

	CloudProvider string  `json:"cloud_provider"`
	InstanceSize  string  `json:"instance_size"`
	InstanceIOPS  *uint64 `json:"iops"`
	DiskSize      *uint64 `json:"disk_size"`
	DiskType      *string `json:"disk_type"`
	DbVendor      string  `json:"database_vendor"`
	DbVersion     string  `json:"database_version"`
	Name          string  `json:"cluster_name"`
	Status        string  `json:"cluster_status"`
	StatusText    string  `json:"cluster_status_text"`
	Type          string  `json:"cluster_type"`
	TypeName      string  `json:"cluster_type_name"`
	Size          int64   `json:"cluster_size"`

	SSLEnabled       bool     `json:"ssl_enabled"`
	HighAvailability bool     `json:"high_availability"`
	Tags             []string `json:"tags"`
	AZS              []string `json:"azs"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Vpc *struct {
		VpcUUID string `json:"vpc_uuid"`
	} `json:"vpc"`

	Region struct {
		Code string `json:"code"`
	} `json:"region"`

	MaintenanceSettings *ccx.MaintenanceSettings `json:"maintenance_settings"`
	Notifications       ccx.Notifications        `json:"notifications"`

	PrimaryUrl string `json:"primary_url"`
	ReplicaUrl string `json:"replica_url"`

	DbAccount struct {
		Username   string `json:"database_username"`
		Password   string `json:"database_password"`
		Host       string `json:"database_host"`
		Database   string `json:"database_database"`
		Privileges string `json:"database_privileges"`
	} `json:"db_account"`
}

func getPortFromDatastore(c ccx.Datastore) (int, error) {
	var port int

	for _, n := range c.Hosts {
		if n.Role == "primary" && n.Port != 0 {
			port = n.Port
			break
		}

		if port == 0 && n.Port != 0 {
			port = n.Port
		}
	}

	if port == 0 {
		return 0, errors.New("no port found")
	}

	return port, nil
}

func (svc *DatastoreService) Read(ctx context.Context, id string) (*ccx.Datastore, error) {
	var rs getDatastoreResponse

	err := svc.client.Get(ctx, "/api/deployment/v3/data-stores/"+id, &rs)
	if err != nil {
		return nil, err
	}

	switch rs.Status {
	case "DEPLOY_FAILED",
		"DELETING",
		"DELETE_FAILED",
		"DELETED":
		return nil, ccx.ResourceNotFoundErr
	}

	c := ccx.Datastore{
		ID:                  rs.ID,
		Name:                rs.Name,
		Size:                rs.Size,
		DBVendor:            rs.DbVendor,
		DBVersion:           rs.DbVersion,
		Type:                rs.Type,
		Tags:                rs.Tags,
		CloudProvider:       rs.CloudProvider,
		CloudRegion:         rs.Region.Code,
		InstanceSize:        rs.InstanceSize,
		VolumeType:          lib.StringVal(rs.DiskType),
		VolumeSize:          lib.Uint64Val(rs.DiskSize),
		VolumeIOPS:          lib.Uint64Val(rs.InstanceIOPS),
		HAEnabled:           rs.HighAvailability,
		AvailabilityZones:   rs.AZS,
		Notifications:       rs.Notifications,
		MaintenanceSettings: rs.MaintenanceSettings,
		PrimaryUrl:          rs.PrimaryUrl,
		ReplicaUrl:          rs.ReplicaUrl,
		Username:            rs.DbAccount.Username,
		Password:            rs.DbAccount.Password,
		DbName:              rs.DbAccount.Database,
	}

	if rs.Vpc != nil {
		c.VpcUUID = rs.Vpc.VpcUUID
	}

	if fw, err := svc.GetFirewallRules(ctx, id); err == nil {
		c.FirewallRules = fw
	} else if !errors.Is(err, ccx.ResourceNotFoundErr) {
		return nil, fmt.Errorf("getting firewall rules: %w", err)
	}

	if h, err := svc.GetHosts(ctx, id); err == nil {
		c.Hosts = h
	} else if !errors.Is(err, ccx.ResourceNotFoundErr) {
		return nil, fmt.Errorf("getting hosts: %w", err)
	}

	port, err := getPortFromDatastore(c)
	if err != nil {
		tflog.Warn(ctx, "failed to get port for store, reported dsn might be incorrect", map[string]any{
			"id":  id,
			"err": err.Error(),
		})
	}

	c.PrimaryDsn = dsn(rs.DbVendor, c.PrimaryUrl, port, rs.DbAccount.Username, rs.DbAccount.Password, rs.DbAccount.Database)
	c.ReplicaDsn = dsn(rs.DbVendor, c.ReplicaUrl, port, rs.DbAccount.Username, rs.DbAccount.Password, rs.DbAccount.Database)

	return &c, nil
}

func dsn(vendor string, host string, port int, username, password, dbname string) string {
	var service string

	if !strings.Contains(host, ":") {
		host += ":" + strconv.Itoa(port)
	}

	switch vendor {
	default:
		return ""
	case "mysql", "mariadb", "percona":
		service = "mysql"
	case "postgres", "pgsql":
		service = "postgres"
	case "redis", "valkey":
		service = "rediss"
	case "microsoft":
		return `Data Source=` + host + `;User ID=` + username + `;Password=` + password + `;Database=` + dbname
	}

	return service + "://" + username + ":" + password + "@" + host + "/" + dbname
}
