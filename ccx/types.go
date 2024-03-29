package ccx

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

type Datastore struct {
	ID                string
	Name              string
	Size              int64
	DBVendor          string
	DBVersion         string
	Type              string
	Tags              []string
	CloudProvider     string
	CloudRegion       string
	InstanceSize      string
	VolumeType        string
	VolumeSize        int64
	VolumeIOPS        int64
	NetworkType       string
	HAEnabled         bool
	VpcUUID           string
	AvailabilityZones []string
}

// String representation of the Datastore, useful for debugging
func (c Datastore) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		return fmt.Sprintf(`{"id": "%s", "name": "%s"}`, c.ID, c.Name)
	}
	return string(b)
}

// DatastoreService is used to manage datastores
type DatastoreService interface {
	Create(ctx context.Context, c Datastore) (*Datastore, error)
	Read(ctx context.Context, id string) (*Datastore, error)
	Update(ctx context.Context, c Datastore) (*Datastore, error)
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
