package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/severalnines/terraform-provider-ccx/ccx"
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
	VolumeSize   int64  `json:"volume_size"`
	VolumeIOPS   int64  `json:"volume_iops"`
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

func CreateRequestFromDatastore(c ccx.Datastore) createStoreRequest {
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

	var volumeSize int64
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

type DatastoreResponse struct {
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

func DatastoreFromResponse(r DatastoreResponse) ccx.Datastore {
	return ccx.Datastore{
		ID:                r.UUID,
		Name:              r.Name,
		Size:              r.Size,
		DBVendor:          r.DbVendor,
		DBVersion:         r.DbVersion,
		Type:              r.Type,
		Tags:              r.Tags,
		CloudProvider:     r.CloudProvider,
		CloudRegion:       r.Region,
		InstanceSize:      r.InstanceSize,
		VolumeType:        pString(r.DiskType),
		VolumeSize:        int64(pUint64(r.DiskSize)),
		VolumeIOPS:        int64(pUint64(r.IOPS)),
		NetworkType:       "", // todo
		HAEnabled:         r.HighAvailability,
		VpcUUID:           pString(r.VpcUUID),
		AvailabilityZones: r.AZS,
	}
}

func (svc *DatastoreService) Create(ctx context.Context, c ccx.Datastore) (*ccx.Datastore, error) {
	cr := CreateRequestFromDatastore(c)

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(cr); err != nil {
		return nil, errors.Join(ccx.RequestEncodingErr, err)
	}

	url := svc.baseURL + "/api/prov/api/v2/cluster"
	req, err := http.NewRequest(http.MethodPost, url, &b)
	if err != nil {
		return nil, errors.Join(ccx.CreateFailedErr, ccx.RequestInitializationErr, err)
	}

	token, err := svc.auth.Auth(ctx)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", token)
	client := &http.Client{Timeout: ccx.DefaultTimeout}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Join(ccx.RequestSendingErr, err)
	}

	if res.StatusCode == http.StatusBadRequest {
		return nil, fmt.Errorf("%w: %w", ccx.CreateFailedErr, lib.ErrorFromErrorResponse(res.Body))
	}

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf(":%w :%w: status = %d", ccx.CreateFailedErr, lib.ErrorFromErrorResponse(res.Body), res.StatusCode)
	}

	var rs DatastoreResponse
	if err := lib.DecodeJsonInto(res.Body, &rs); err != nil {
		return nil, fmt.Errorf("%w: %w", ccx.CreateFailedErr, err)
	}

	newDatastore := DatastoreFromResponse(rs)

	status, err := svc.jobs.Await(ctx, rs.UUID, deployStoreJob, svc.timeout)
	if err != nil {
		return nil, fmt.Errorf("%w: awaiting deploy job: %w", ccx.CreateFailedErr, err)
	} else if status != jobStatusFinished {
		return nil, fmt.Errorf("%w: deploy job failed: %s", ccx.CreateFailedErr, status)
	}

	return &newDatastore, nil
}
