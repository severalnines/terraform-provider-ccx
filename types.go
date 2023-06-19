package terraform_provider_ccx

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

type TerraformConfiguration struct {
	ClientID     string
	ClientSecret string
	BaseURL      string
	IsDevMode    bool
	Mockfile     string
}

type TerraformResource interface {
	Configure(ctx context.Context, cfg TerraformConfiguration) error

	Name() string
	Schema() *schema.Resource
	Create(*schema.ResourceData, interface{}) error
	Read(*schema.ResourceData, interface{}) error
	Update(*schema.ResourceData, interface{}) error
	Delete(*schema.ResourceData, interface{}) error
	// Exists(*schema.ResourceData, interface{}) (bool, error)
}

// Cluster is a database cluster
type Cluster struct {
	ID                string   `json:"id"`
	ClusterName       string   `json:"cluster_name"`
	ClusterSize       int64    `json:"cluster_size"`
	DBVendor          string   `json:"db_vendor"`
	DBVersion         string   `json:"db_version"`
	ClusterType       string   `json:"cluster_type"`
	Tags              []string `json:"tags"`
	CloudSpace        string   `json:"cloud_space"`
	CloudProvider     string   `json:"cloud_provider"`
	CloudRegion       string   `json:"cloud_region"`
	InstanceSize      string   `json:"instance_size"` // "Tiny" ... "2X-Large"
	VolumeType        string   `json:"volume_type"`
	VolumeSize        int64    `json:"volume_size"`
	VolumeIOPS        int64    `json:"volume_iops"`
	NetworkType       string   `json:"network_type"` // public/private
	HAEnabled         bool     `json:"ha_enabled"`
	VpcUUID           string   `json:"vpc_uuid"`
	AvailabilityZones []string `json:"availability_zones"`
}

// String representation of the Cluster, useful for debugging
func (c Cluster) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		return fmt.Sprintf(`{"id": "%s", "name": "%s"}`, c.ID, c.ClusterName)
	}
	return string(b)
}

// ClusterService is used to manage VPCs
type ClusterService interface {
	Create(ctx context.Context, c Cluster) (*Cluster, error)
	Read(ctx context.Context, id string) (*Cluster, error)
	Update(ctx context.Context, c Cluster) (*Cluster, error)
	Delete(ctx context.Context, id string) error
}

type VPC struct {
	ID            string `json:"ID"`
	Name          string `json:"name"`
	CloudSpace    string `json:"cloudspace"`
	CloudProvider string `json:"cloud"`
	Region        string `json:"region"`
	CidrIpv4Block string `json:"cidr_ipv4_block"`
}

// String representation of the VPC, useful for debugging
func (v VPC) String() string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf(`{"id": "%s", "name": "%s"}`, v.ID, v.Name)
	}
	return string(b)
}

type VPCService interface {
	Create(ctx context.Context, vpc VPC) (*VPC, error)
	Read(ctx context.Context, id string) (*VPC, error)
	Update(ctx context.Context, vpc VPC) (*VPC, error)
	Delete(ctx context.Context, id string) error
}
