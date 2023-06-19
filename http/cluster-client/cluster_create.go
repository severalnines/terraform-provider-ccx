package cluster_client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	ccxprov "github.com/severalnines/terraform-provider-ccx"
	chttp "github.com/severalnines/terraform-provider-ccx/http"
	"github.com/severalnines/terraform-provider-ccx/pointers"
)

type CreateRequestGeneral struct {
	ClusterName string   `json:"cluster_name"`
	ClusterSize int64    `json:"cluster_size"`
	DBVendor    string   `json:"db_vendor"`
	DBVersion   string   `json:"db_version"`
	ClusterType string   `json:"cluster_type"`
	Tags        []string `json:"tags"`
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

func CreateRequestFromCluster(c ccxprov.Cluster) CreateRequest {
	general := CreateRequestGeneral{
		ClusterName: c.ClusterName,
		ClusterSize: c.ClusterSize,
		DBVendor:    c.DBVendor,
		DBVersion:   c.DBVersion,
		ClusterType: c.ClusterType,
		Tags:        c.Tags,
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

type ClusterResponse struct {
	ClusterUUID             string   `json:"uuid"`
	ClusterName             string   `json:"cluster_name"`
	ClusterType             string   `json:"cluster_type"`
	ClusterRegion           string   `json:"region"`
	CloudProvider           string   `json:"cloud_provider"`
	ClusterSize             int64    `json:"cluster_size"`
	ClusterDbVendor         string   `json:"database_vendor"`
	ClusterDbVersion        string   `json:"database_version"`
	ClusterInstanceSize     string   `json:"instance_size"`
	ClusterInstanceDiskType *string  `json:"cluster_instance_disk_type"`
	ClusterInstanceIOPS     *uint64  `json:"iops"`
	ClusterInstanceDiskSize *uint64  `json:"disk_size"`
	HighAvailability        bool     `json:"high_availability"`
	VpcUUID                 *string  `json:"vpc_uuid"`
	Tags                    []string `json:"tags"`
	AZS                     []string `json:"azs"`
	UsePublicIPs            bool     `json:"use_public_ips"`
}

func ClusterFromResponse(r ClusterResponse) ccxprov.Cluster {
	return ccxprov.Cluster{
		ID:                r.ClusterUUID,
		ClusterName:       r.ClusterName,
		ClusterSize:       r.ClusterSize,
		DBVendor:          r.ClusterDbVendor,
		DBVersion:         r.ClusterDbVersion,
		ClusterType:       r.ClusterType,
		Tags:              r.Tags,
		CloudProvider:     r.CloudProvider,
		CloudRegion:       r.ClusterRegion,
		InstanceSize:      r.ClusterInstanceSize,
		VolumeType:        pointers.String(r.ClusterInstanceDiskType),
		VolumeSize:        int64(pointers.Uint64(r.ClusterInstanceDiskSize)),
		VolumeIOPS:        int64(pointers.Uint64(r.ClusterInstanceIOPS)),
		NetworkType:       "", // todo
		HAEnabled:         r.HighAvailability,
		VpcUUID:           pointers.String(r.VpcUUID),
		AvailabilityZones: r.AZS,
	}
}

// Create a new clusters
func (cli *Client) Create(ctx context.Context, c ccxprov.Cluster) (*ccxprov.Cluster, error) {
	cr := CreateRequestFromCluster(c)

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(cr); err != nil {
		return nil, errors.Join(ccxprov.RequestEncodingErr, err)
	}

	url := cli.conn.BaseURL + "/api/prov/api/v2/cluster"
	req, err := http.NewRequest(http.MethodPost, url, &b)
	if err != nil {
		return nil, errors.Join(ccxprov.RequestInitializationErr, err)
	}

	token, err := cli.auth.Auth(ctx)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", token)
	client := &http.Client{Timeout: cli.conn.Timeout}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Join(ccxprov.RequestSendingErr, err)
	}

	if res.StatusCode == http.StatusBadRequest {
		return nil, chttp.ErrorFromErrorResponse(res.Body)
	}

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%w: status = %d", chttp.ErrorFromErrorResponse(res.Body), res.StatusCode)
	}

	var rs ClusterResponse
	if err := chttp.DecodeJsonInto(res.Body, &rs); err != nil {
		return nil, err
	}

	newCluster := ClusterFromResponse(rs)

	if err := cli.LoadAll(ctx); err != nil {
		return nil, errors.Join(ccxprov.ResourcesLoadFailedErr, err)
	}

	return &newCluster, nil
}
