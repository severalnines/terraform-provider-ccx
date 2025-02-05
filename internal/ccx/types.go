package ccx

import (
	"context"
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
	HAEnabled         bool
	VpcUUID           string
	ParameterGroupID  string
	AvailabilityZones []string

	FirewallRules []FirewallRule
	Hosts         []Host

	Notifications       Notifications
	MaintenanceSettings *MaintenanceSettings

	PrimaryUrl string
	PrimaryDsn string
	ReplicaUrl string
	ReplicaDsn string
	Username   string
	Password   string
	DbName     string
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
	Port          int
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
	return fmt.Sprintf(`{"id": "%s", "name": "%s"}`, c.ID, c.Name)
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
	SetFirewallRules(ctx context.Context, storeID string, firewalls []FirewallRule) error
	SetMaintenanceSettings(ctx context.Context, storeID string, settings MaintenanceSettings) error
	ApplyParameterGroup(ctx context.Context, id, group string) error
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
	return fmt.Sprintf(`{"id": "%s", "name": "%s"}`, v.ID, v.Name)
}

type VPCService interface {
	Create(ctx context.Context, vpc VPC) (*VPC, error)
	Read(ctx context.Context, id string) (*VPC, error)
	Update(ctx context.Context, vpc VPC) (*VPC, error)
	Delete(ctx context.Context, id string) error
}

type InstanceSize struct {
	Code string `json:"code"`
	Type string `json:"type"`
}

type ContentService interface {
	InstanceSizes(ctx context.Context) (map[string][]InstanceSize, error)
	AvailabilityZones(ctx context.Context, provider, region string) ([]string, error)
}

type ParameterGroup struct {
	ID              string            `json:"uuid"`
	Name            string            `json:"name"`
	DatabaseVendor  string            `json:"database_vendor"`
	DatabaseVersion string            `json:"database_version"`
	DatabaseType    string            `json:"database_type"`
	DbParameters    map[string]string `json:"db_parameters"`
}

type ParameterGroupService interface {
	Create(ctx context.Context, p ParameterGroup) (*ParameterGroup, error)
	Read(ctx context.Context, id string) (*ParameterGroup, error)
	Update(ctx context.Context, p ParameterGroup) error
	Delete(ctx context.Context, id string) error
}

type JobType string

const (
	DeployStoreJob    JobType = "JOB_TYPE_DEPLOY_DATASTORE"
	ModifyDbConfigJob JobType = "JOB_TYPE_MODIFYDBCONFIG"
	DestroyStoreJob   JobType = "JOB_TYPE_DESTROY_DATASTORE"
	AddNodeJob        JobType = "JOB_TYPE_ADD_NODE"
	RemoveNodeJob     JobType = "JOB_TYPE_REMOVE_NODE"
)

type JobStatus string

const (
	JobStatusUnknown  JobStatus = "JOB_STATUS_UNKNOWN"
	JobStatusRunning  JobStatus = "JOB_STATUS_RUNNING"
	JobStatusFinished JobStatus = "JOB_STATUS_FINISHED"
	JobStatusErrored  JobStatus = "JOB_STATUS_ERRORED"
)

type JobService interface {
	Await(ctx context.Context, storeID string, job JobType) (JobStatus, error)
	GetStatus(_ context.Context, storeID string, job JobType) (JobStatus, error)
}
