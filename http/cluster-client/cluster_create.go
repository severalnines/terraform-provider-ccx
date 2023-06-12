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

type CreateRequest struct {
	General struct {
		ClusterName string   `json:"cluster_name"`
		ClusterSize int64    `json:"cluster_size"`
		DBVendor    string   `json:"db_vendor"`
		DBVersion   string   `json:"db_version"`
		ClusterType string   `json:"cluster_type"`
		Tags        []string `json:"tags"`
	} `json:"general"`
	Cloud struct {
		CloudProvider string `json:"cloud_provider"`
		CloudRegion   string `json:"cloud_region"`
	} `json:"cloud"`
	Instance struct {
		InstanceSize string `json:"instance_size"` // "Tiny" ... "2X-Large"
		VolumeType   string `json:"volume_type"`
		VolumeSize   uint64 `json:"volume_size"`
		VolumeIOPS   uint64 `json:"volume_iops"`
	} `json:"instance"`
	Network struct {
		NetworkType       string   `json:"network_type"` // public/private
		HAEnabled         bool     `json:"ha_enabled"`
		VpcUUID           string   `json:"vpc_uuid"`
		AvailabilityZones []string `json:"availability_zones"`
	} `json:"network"`
}

func CreateRequestFromCluster(c ccxprov.Cluster) CreateRequest {
	return CreateRequest{
		General: struct {
			ClusterName string   `json:"cluster_name"`
			ClusterSize int64    `json:"cluster_size"`
			DBVendor    string   `json:"db_vendor"`
			DBVersion   string   `json:"db_version"`
			ClusterType string   `json:"cluster_type"`
			Tags        []string `json:"tags"`
		}{
			ClusterName: c.ClusterName,
			ClusterSize: c.ClusterSize,
			DBVendor:    c.DBVendor,
			DBVersion:   c.DBVersion,
			ClusterType: c.ClusterType,
			Tags:        c.Tags,
		},
		Cloud: struct {
			CloudProvider string `json:"cloud_provider"`
			CloudRegion   string `json:"cloud_region"`
		}{
			CloudProvider: c.CloudProvider,
			CloudRegion:   c.CloudRegion,
		},
		Instance: struct {
			InstanceSize string `json:"instance_size"`
			VolumeType   string `json:"volume_type"`
			VolumeSize   uint64 `json:"volume_size"`
			VolumeIOPS   uint64 `json:"volume_iops"`
		}{
			InstanceSize: c.InstanceSize,
			VolumeType:   c.VolumeType,
			VolumeSize:   uint64(c.VolumeSize),
			VolumeIOPS:   uint64(c.VolumeIOPS),
		},
		Network: struct {
			NetworkType       string   `json:"network_type"`
			HAEnabled         bool     `json:"ha_enabled"`
			VpcUUID           string   `json:"vpc_uuid"`
			AvailabilityZones []string `json:"availability_zones"`
		}{
			NetworkType:       c.NetworkType,
			HAEnabled:         c.HAEnabled,
			VpcUUID:           c.VpcUUID,
			AvailabilityZones: c.AvailabilityZones,
		},
	}
}

type CreateResponse struct {
	ClusterUUID             string   `json:"uuid" reform:"cluster_uuid,pk"`
	ClusterName             string   `json:"cluster_name" reform:"cluster_name"`
	ClusterType             string   `json:"cluster_type" reform:"cluster_type"`
	ClusterRegion           string   `json:"region" reform:"cluster_region"`
	CloudProvider           string   `json:"cloud_provider" reform:"cluster_cloud"`
	ClusterSize             int64    `json:"cluster_size" reform:"cluster_size"`
	ClusterDbVendor         string   `json:"database_vendor" reform:"cluster_db_vendor"`
	ClusterDbVersion        string   `json:"database_version" reform:"cluster_db_version"`
	ClusterInstanceSize     string   `json:"instance_size" reform:"cluster_instance_size"`
	ClusterInstanceDiskType *string  `json:"cluster_instance_disk_type" reform:"cluster_instance_disk_type"`
	ClusterInstanceIOPS     *uint64  `json:"iops" reform:"cluster_instance_iops"`
	ClusterInstanceDiskSize *uint64  `json:"disk_size" reform:"cluster_instance_disk_size"`
	HighAvailability        bool     `json:"high_availability" reform:"high_availability"`
	VpcUUID                 *string  `json:"vpc_uuid" reform:"vpc_uuid"`
	Tags                    []string `json:"tags" reform:"tags"`
	AZS                     []string `json:"azs" reform:"azs"`
	UsePublicIPs            bool     `json:"use_public_ips" reform:"use_public_ips"`
}

func ClusterFromCreateResponse(r CreateResponse) ccxprov.Cluster {
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
		return nil, fmt.Errorf("%w: status = %d", ccxprov.ResponseStatusFailedErr, res.StatusCode)
	}

	var rs CreateResponse
	if err := chttp.DecodeJsonInto(res.Body, &rs); err != nil {
		return nil, err
	}

	newCluster := ClusterFromCreateResponse(rs)
	// c.ID = newCluster.ID
	newCluster.NetworkType = c.NetworkType // todo should come from server response
	newCluster.Tags = c.Tags
	newCluster.AvailabilityZones = c.AvailabilityZones

	if err := cli.LoadAll(ctx); err != nil {
		return nil, errors.Join(ccxprov.ResourcesLoadFailedErr, err)
	}

	return &newCluster, nil
}
