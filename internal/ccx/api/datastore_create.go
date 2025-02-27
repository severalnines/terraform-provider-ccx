package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
	"github.com/severalnines/terraform-provider-ccx/internal/lib"
)

type createStoreGeneral struct {
	Name      string   `json:"cluster_name"`
	Size      int64    `json:"cluster_size"`
	DBVendor  string   `json:"db_vendor"`
	DBVersion string   `json:"db_version"`
	Type      string   `json:"cluster_type"`
	Tags      []string `json:"tags"`
}

type createStoreCloud struct {
	CloudProvider string `json:"cloud_provider"`
	CloudRegion   string `json:"cloud_region"`
}

type createStoreInstance struct {
	InstanceSize string `json:"instance_size"` // "Tiny" ... "2X-Large"
	VolumeType   string `json:"volume_type"`
	VolumeSize   uint64 `json:"volume_size"`
	VolumeIOPS   uint64 `json:"volume_iops"`
}

type createStoreNetwork struct {
	NetworkType       string   `json:"network_type"` // public/private
	HAEnabled         bool     `json:"ha_enabled"`
	VpcUUID           string   `json:"vpc_uuid"`
	AvailabilityZones []string `json:"availability_zones"`
}

type createStoreRequest struct {
	General       createStoreGeneral  `json:"general"`
	Cloud         createStoreCloud    `json:"cloud"`
	Instance      createStoreInstance `json:"instance"`
	Network       createStoreNetwork  `json:"network"`
	Notifications notifications       `json:"notifications"`
}

func createRequestFromDatastore(c ccx.Datastore) createStoreRequest {
	general := createStoreGeneral{
		Name:      c.Name,
		Size:      c.Size,
		DBVendor:  c.DBVendor,
		DBVersion: c.DBVersion,
		Type:      c.Type,
		Tags:      c.Tags,
	}

	cloud := createStoreCloud{
		CloudProvider: c.CloudProvider,
		CloudRegion:   c.CloudRegion,
	}

	var volumeSize uint64
	if c.VolumeSize == 0 {
		volumeSize = 80
	} else {
		volumeSize = c.VolumeSize
	}

	instance := createStoreInstance{
		InstanceSize: c.InstanceSize,
		VolumeType:   c.VolumeType,
		VolumeSize:   volumeSize,
		VolumeIOPS:   c.VolumeIOPS,
	}

	networkType := "public"

	if c.VpcUUID != "" {
		networkType = "private"
	}

	network := createStoreNetwork{
		NetworkType:       networkType,
		HAEnabled:         c.HAEnabled,
		VpcUUID:           c.VpcUUID,
		AvailabilityZones: c.AvailabilityZones,
	}

	notifs := notifications{
		Enabled: c.Notifications.Enabled,
		Emails:  c.Notifications.Emails,
	}

	return createStoreRequest{
		General:       general,
		Cloud:         cloud,
		Instance:      instance,
		Network:       network,
		Notifications: notifs,
	}
}

func allocateAzs(allAzs, existing []string, n int) []string {
	if n <= 0 {
		return nil
	}

	m := make(map[string]int, len(allAzs))

	for _, a := range allAzs {
		m[a] = 0
	}

	for _, e := range existing {
		if _, ok := m[e]; ok { // skip those not in the list, might be older ones no longer available
			m[e] += 1
		}
	}

	azs := make([]lib.CountedItem, 0, len(m))
	for name, count := range m {
		azs = append(azs, lib.CountedItem{
			Name:  name,
			Count: count,
		})
	}

	ls := lib.AllocateN(azs, n)

	return ls
}

type datastoreResponse struct {
	UUID             string   `json:"uuid"`
	Name             string   `json:"cluster_name"`
	Type             string   `json:"cluster_type"`
	Region           string   `json:"region"`
	CloudProvider    string   `json:"cloud_provider"`
	Size             int64    `json:"cluster_size"`
	DbVendor         string   `json:"database_vendor"`
	DbVersion        string   `json:"database_version"`
	InstanceSize     string   `json:"instance_size"`
	DiskType         *string  `json:"cluster_instance_disk_type"`
	IOPS             *uint64  `json:"iops"`
	DiskSize         *uint64  `json:"disk_size"`
	HighAvailability bool     `json:"high_availability"`
	VpcUUID          *string  `json:"vpc_uuid"`
	Tags             []string `json:"tags"`
	AZS              []string `json:"azs"`
}

func (svc *DatastoreService) Create(ctx context.Context, c ccx.Datastore) (*ccx.Datastore, error) {
	cr := createRequestFromDatastore(c)

	if n, h := len(c.AvailabilityZones), int(c.Size); c.VpcUUID == "" && n < h { // allocate AZs if public and need is less than have
		allAzs, err := svc.contentSvc.AvailabilityZones(ctx, c.CloudProvider, c.CloudRegion)
		if err != nil {
			return nil, fmt.Errorf("creating datastore: %w: %w", ccx.AllocatingAZsErr, err)
		}

		c.AvailabilityZones = allocateAzs(allAzs, nil, h-n)
	}

	res, err := svc.client.Do(ctx, http.MethodPost, "/api/prov/api/v2/cluster", cr)
	if err != nil {
		return nil, fmt.Errorf("creating datastore: %w", err)
	}

	var rs datastoreResponse
	if err := lib.DecodeJsonInto(res.Body, &rs); err != nil {
		return nil, fmt.Errorf("creating datastore: %w", err)
	}

	partialDatastore := &ccx.Datastore{ID: rs.UUID}

	status, err := svc.jobs.Await(ctx, rs.UUID, ccx.DeployStoreJob)
	if err != nil {
		return partialDatastore, fmt.Errorf("%w: awaiting deploy job: %w", ccx.CreateFailedReadErr, err)
	} else if status != ccx.JobStatusFinished {
		return partialDatastore, fmt.Errorf("%w: deploy job failed: %s", ccx.CreateFailedReadErr, status)
	}

	newDatastore, err := svc.Read(ctx, rs.UUID)
	if err != nil {
		return partialDatastore, fmt.Errorf("%w: %w", ccx.CreateFailedReadErr, err)
	}

	if c.ParameterGroupID != "" {
		err = svc.ApplyParameterGroup(ctx, rs.UUID, c.ParameterGroupID)
		if err != nil {
			return newDatastore, fmt.Errorf("%w: %w", ccx.ApplyDbParametersFailedErr, err)
		}
	}

	return newDatastore, nil
}
