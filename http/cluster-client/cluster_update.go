package cluster_client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	ccxprov "github.com/severalnines/terraform-provider-ccx"
	chttp "github.com/severalnines/terraform-provider-ccx/http"
	cstrings "github.com/severalnines/terraform-provider-ccx/strings"
)

type UpdateRequest struct {
	ClusterName string `json:"cluster_name"`

	// AddNodes      *UpdateAddNodesRequest `json:"add_nodes"`
	NewVolumeSize uint `json:"new_volume_size"`
}

// type UpdateAddNodesRequest struct {
// 	Specs []UpdateAddNodesSpecsRequest `json:"specs"`
// }

// type UpdateAddNodesSpecsRequest struct {
// 	InstanceType string `json:"instance_size"`
// 	AZ           string `json:"availability_zone"`
// }

func (cli *Client) Update(ctx context.Context, c ccxprov.Cluster) (*ccxprov.Cluster, error) {
	old, err := cli.Read(ctx, c.ID)
	if err == ccxprov.ResourceNotFoundErr {
		return nil, ccxprov.ResourceNotFoundErr
	}

	if hasCan, err := HasSupportedChanges(*old, c); err != nil {
		return nil, err
	} else if !hasCan {
		return old, nil
	}

	var ur UpdateRequest

	if old.ClusterName != c.ClusterName {
		ur.ClusterName = c.ClusterName
	}

	// if n := c.ClusterSize - old.ClusterSize; n > 0 {
	// 	ur.AddNodes = &UpdateAddNodesRequest{}
	//
	// 	for n > 0 {
	// 		n -= 1
	// 		ur.AddNodes.Specs = append(ur.AddNodes.Specs, UpdateAddNodesSpecsRequest{})
	// 	}
	// }

	if old.VolumeSize != c.VolumeSize {
		ur.NewVolumeSize = uint(c.VolumeSize)
	}

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(ur); err != nil {
		return nil, errors.Join(ccxprov.RequestEncodingErr, err)
	}

	url := cli.conn.BaseURL + "/api/prov/api/v2/cluster/" + c.ID
	req, err := http.NewRequest(http.MethodPatch, url, &b)
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

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status = %d", ccxprov.ResponseStatusFailedErr, res.StatusCode)
	}

	var rs ClusterResponse
	if err := chttp.DecodeJsonInto(res.Body, &rs); err != nil {
		return nil, err
	}

	updatedCluster := ClusterFromResponse(rs)

	if err := cli.LoadAll(ctx); err != nil {
		return nil, errors.Join(ccxprov.ResourcesLoadFailedErr, err)
	}

	return &updatedCluster, nil
}

func HasSupportedChanges(old, c ccxprov.Cluster) (bool, error) {
	var (
		hasCan, hasCant bool
		fields          []string
	)

	if old.ClusterName != c.ClusterName {
		hasCan = true
	}

	if old.ClusterSize != c.ClusterSize {
		hasCant = true
		fields = append(fields, "cluster_size")
	}

	if old.VolumeSize != c.VolumeSize {
		// hasCan = true
		hasCant = true
		fields = append(fields, "volume_size")
	}

	if old.DBVendor != c.DBVendor {
		hasCant = true
		fields = append(fields, "db_vendor")
	}

	if old.DBVersion != c.DBVersion {
		hasCant = true
		fields = append(fields, "db_version")
	}

	if old.ClusterType != c.ClusterType {
		hasCant = true
		fields = append(fields, "cluster_type")
	}

	// if !cstrings.Sames(old.Tags, c.Tags) {
	// 	hasCant = true
	// 	fields = append(fields, "tags")
	// }

	if old.CloudSpace != c.CloudSpace {
		hasCant = true
		fields = append(fields, "cloud_space")
	}

	if old.CloudProvider != c.CloudProvider {
		hasCant = true
		fields = append(fields, "cloud_provider")
	}

	if old.CloudRegion != c.CloudRegion {
		hasCant = true
		fields = append(fields, "cloud_region")
	}

	if old.InstanceSize != c.InstanceSize {
		hasCant = true
		fields = append(fields, "instance_size")
	}

	if old.VolumeType != c.VolumeType {
		hasCant = true
		fields = append(fields, "volume_type")
	}

	if old.VolumeIOPS != c.VolumeIOPS {
		hasCant = true
		fields = append(fields, "volume_iops")
	}

	// if old.NetworkType != c.NetworkType {
	// 	hasCant = true
	// 	fields = append(fields, "network_type")
	// }

	if old.HAEnabled != c.HAEnabled {
		hasCant = true
		fields = append(fields, "ha_enabled")
	}

	if old.VpcUUID != c.VpcUUID {
		hasCant = true
		fields = append(fields, "vpc_uuid")
	}

	if !cstrings.Sames(old.AvailabilityZones, c.AvailabilityZones) {
		hasCant = true
		fields = append(fields, "availability_zones")
	}

	if hasCant {
		return hasCan, fmt.Errorf("%w: %s", ccxprov.UpdateNotSupportedErr, strings.Join(fields, ", "))
	}

	return hasCan, nil
}
