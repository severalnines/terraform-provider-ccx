package datastore_client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/severalnines/terraform-provider-ccx/ccx"
	chttp "github.com/severalnines/terraform-provider-ccx/http"
	"github.com/severalnines/terraform-provider-ccx/pointers"
)

type CreateRequestGeneral struct {
	Name      string   `json:"cluster_name"`
	Size      int64    `json:"cluster_size"`
	DBVendor  string   `json:"db_vendor"`
	DBVersion string   `json:"db_version"`
	Type      string   `json:"cluster_type"`
	Tags      []string `json:"tags"`
}

type CreateRequestCloud struct {
	CloudProvider string `json:"cloud_provider"`
	CloudRegion   string `json:"cloud_region"`
}

type CreateRequestInstance struct {
	InstanceSize string `json:"instance_size"` // "Tiny" ... "2X-Large"
	VolumeType   string `json:"volume_type"`
	VolumeSize   int64  `json:"volume_size"`
	VolumeIOPS   int64  `json:"volume_iops"`
}

type CreateRequestNetwork struct {
	NetworkType       string   `json:"network_type"` // public/private
	HAEnabled         bool     `json:"ha_enabled"`
	VpcUUID           string   `json:"vpc_uuid"`
	AvailabilityZones []string `json:"availability_zones"`
}

type CreateRequest struct {
	General  CreateRequestGeneral  `json:"general"`
	Cloud    CreateRequestCloud    `json:"cloud"`
	Instance CreateRequestInstance `json:"instance"`
	Network  CreateRequestNetwork  `json:"network"`
}

func CreateRequestFromDatastore(c ccx.Datastore) CreateRequest {
	general := CreateRequestGeneral{
		Name:      c.Name,
		Size:      c.Size,
		DBVendor:  c.DBVendor,
		DBVersion: c.DBVersion,
		Type:      c.Type,
		Tags:      c.Tags,
	}

	cloud := CreateRequestCloud{
		CloudProvider: c.CloudProvider,
		CloudRegion:   c.CloudRegion,
	}

	var volumeSize int64
	if c.VolumeSize == 0 {
		volumeSize = 80
	} else {
		volumeSize = c.VolumeSize
	}

	instance := CreateRequestInstance{
		InstanceSize: c.InstanceSize,
		VolumeType:   c.VolumeType,
		VolumeSize:   volumeSize,
		VolumeIOPS:   c.VolumeIOPS,
	}

	network := CreateRequestNetwork{
		NetworkType:       c.NetworkType,
		HAEnabled:         c.HAEnabled,
		VpcUUID:           c.VpcUUID,
		AvailabilityZones: c.AvailabilityZones,
	}

	return CreateRequest{
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
		VolumeType:        pointers.String(r.DiskType),
		VolumeSize:        int64(pointers.Uint64(r.DiskSize)),
		VolumeIOPS:        int64(pointers.Uint64(r.IOPS)),
		NetworkType:       "", // todo
		HAEnabled:         r.HighAvailability,
		VpcUUID:           pointers.String(r.VpcUUID),
		AvailabilityZones: r.AZS,
	}
}

// Create a new datastore
func (cli *Client) Create(ctx context.Context, c ccx.Datastore) (*ccx.Datastore, error) {
	cr := CreateRequestFromDatastore(c)

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(cr); err != nil {
		return nil, errors.Join(ccx.RequestEncodingErr, err)
	}

	url := cli.conn.BaseURL + "/api/prov/api/v2/cluster"
	req, err := http.NewRequest(http.MethodPost, url, &b)
	if err != nil {
		return nil, errors.Join(ccx.RequestInitializationErr, err)
	}

	token, err := cli.auth.Auth(ctx)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", token)
	client := &http.Client{Timeout: cli.conn.Timeout}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Join(ccx.RequestSendingErr, err)
	}

	if res.StatusCode == http.StatusBadRequest {
		return nil, chttp.ErrorFromErrorResponse(res.Body)
	}

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%w: status = %d", chttp.ErrorFromErrorResponse(res.Body), res.StatusCode)
	}

	var rs DatastoreResponse
	if err := chttp.DecodeJsonInto(res.Body, &rs); err != nil {
		return nil, err
	}

	newDatastore := DatastoreFromResponse(rs)

	if err := cli.LoadAll(ctx); err != nil {
		return nil, errors.Join(ccx.ResourcesLoadFailedErr, err)
	}

	return &newDatastore, nil
}
