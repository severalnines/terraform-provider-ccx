package ccx

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	// DefaultBaseURL to access API services
	DefaultBaseURL = "https://app.mydbservice.net"

	// DefaultTimeout for http requests
	DefaultTimeout = time.Second * 30
)

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
	VolumeSize        uint64
	VolumeIOPS        uint64
	NetworkType       string
	HAEnabled         bool
	VpcUUID           string
	AvailabilityZones []string

	DbParams      map[string]string
	FirewallRules []FirewallRule
	Hosts         []Host

	Notifications       Notifications
	MaintenanceSettings *MaintenanceSettings
}

type Host struct {
	ID            string
	CreatedAt     time.Time
	CloudProvider string
	AZ            string
	InstanceType  string
	DiskType      string
	DiskSize      uint64
	Role          string
	Region        string
}

func (h Host) IsPrimary() bool {
	switch r := strings.ToLower(h.Role); r {
	case "primary", "master":
		return true
	}

	return false
}

// String representation of the Datastore, useful for debugging
func (c Datastore) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		return fmt.Sprintf(`{"id": "%s", "name": "%s"}`, c.ID, c.Name)
	}
	return string(b)
}

type FirewallRule struct {
	Source      string `json:"source"`
	Description string `json:"description"`
}

func (f FirewallRule) String() string {
	return fmt.Sprintf(`{"source": "%s", "description": "%s"}`, f.Source, f.Description)
}

type Notifications struct {
	Enabled bool     `json:"enabled"`
	Emails  []string `json:"emails"`
}

type MaintenanceSettings struct {
	DayOfWeek int32 `json:"day_of_week"`
	StartHour int   `json:"start_hour"`
	EndHour   int   `json:"end_hour"`
}

// DatastoreService is used to manage datastores
type DatastoreService interface {
	Create(ctx context.Context, c Datastore) (*Datastore, error)
	Read(ctx context.Context, id string) (*Datastore, error)
	Update(ctx context.Context, old, next Datastore) (*Datastore, error)
	Delete(ctx context.Context, id string) error
	SetParameters(ctx context.Context, storeID string, parameters map[string]string) error
	SetFirewallRules(ctx context.Context, storeID string, firewalls []FirewallRule) error
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
