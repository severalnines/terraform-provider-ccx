package api

import (
	"context"
	"errors"
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
	General  createStoreGeneral  `json:"general"`
	Cloud    createStoreCloud    `json:"cloud"`
	Instance createStoreInstance `json:"instance"`
	Network  createStoreNetwork  `json:"network"`
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

	network := createStoreNetwork{
		NetworkType:       c.NetworkType,
		HAEnabled:         c.HAEnabled,
		VpcUUID:           c.VpcUUID,
		AvailabilityZones: c.AvailabilityZones,
	}

	return createStoreRequest{
		General:  general,
		Cloud:    cloud,
		Instance: instance,
		Network:  network,
	}
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

	res, err := svc.httpcli.Do(ctx, http.MethodPost, "/api/prov/api/v2/cluster", cr)
	if err != nil {
		return nil, errors.Join(ccx.RequestSendingErr, err)
	}

	if res.StatusCode == http.StatusBadRequest {
		return nil, fmt.Errorf("%w: %w", ccx.CreateFailedErr, lib.ErrorFromErrorResponse(res.Body))
	}

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%w :%w: status = %d", ccx.CreateFailedErr, lib.ErrorFromErrorResponse(res.Body), res.StatusCode)
	}

	var rs datastoreResponse
	if err := lib.DecodeJsonInto(res.Body, &rs); err != nil {
		return nil, fmt.Errorf("%w: %w", ccx.CreateFailedErr, err)
	}

	status, err := svc.jobs.Await(ctx, rs.UUID, deployStoreJob)
	if err != nil {
		return nil, fmt.Errorf("%w: awaiting deploy job: %w", ccx.CreateFailedErr, err)
	} else if status != jobStatusFinished {
		return nil, fmt.Errorf("%w: deploy job failed: %s", ccx.CreateFailedErr, status)
	}

	newDatastore, err := svc.Read(ctx, rs.UUID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ccx.CreateFailedErr, err)
	}

	return newDatastore, nil
}
