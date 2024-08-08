package api

import (
	"context"
	"fmt"
	"time"

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
}

func (svc *DatastoreService) Read(ctx context.Context, id string) (*ccx.Datastore, error) {
	var rs getDatastoreResponse

	err := svc.httpcli.Get(ctx, "/api/deployment/v3/data-stores/"+id, &rs)
	if err != nil {
		return nil, err
	}

	c := ccx.Datastore{
		ID:                rs.ID,
		Name:              rs.Name,
		Size:              rs.Size,
		DBVendor:          rs.DbVendor,
		DBVersion:         rs.DbVersion,
		Type:              rs.Type,
		Tags:              rs.Tags,
		CloudProvider:     rs.CloudProvider,
		CloudRegion:       rs.Region.Code,
		InstanceSize:      rs.InstanceSize,
		VolumeType:        lib.StringVal(rs.DiskType),
		VolumeSize:        lib.Uint64Val(rs.DiskSize),
		VolumeIOPS:        lib.Uint64Val(rs.InstanceIOPS),
		HAEnabled:         rs.HighAvailability,
		AvailabilityZones: rs.AZS,
	}

	if rs.Vpc != nil {
		c.VpcUUID = rs.Vpc.VpcUUID
	}

	if p, err := svc.GetParameters(ctx, id); err == nil {
		c.DbParams = p
	} else {
		return nil, fmt.Errorf("getting parameters: %w", err)
	}

	if fw, err := svc.GetFirewallRules(ctx, id); err == nil {
		c.FirewallRules = fw
	} else {
		return nil, fmt.Errorf("getting firewall rules: %w", err)
	}

	if h, err := svc.GetHosts(ctx, id); err == nil {
		c.Hosts = h
	} else {
		return nil, fmt.Errorf("getting hosts: %w", err)
	}

	return &c, nil
}
