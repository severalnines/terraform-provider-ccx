package cluster_client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	ccxprov "github.com/severalnines/terraform-provider-ccx"
	chttp "github.com/severalnines/terraform-provider-ccx/http"
	ccxprovio "github.com/severalnines/terraform-provider-ccx/io"
	"github.com/severalnines/terraform-provider-ccx/pointers"
)

type LoadAllResponse []struct {
	UUID             string  `json:"uuid"`
	Region           string  `json:"region"`
	CloudProvider    string  `json:"cloud_provider"`
	InstanceSize     string  `json:"instance_size"`
	InstanceIOPS     *uint64 `json:"iops"`
	DiskSize         *uint64 `json:"disk_size"`
	DiskType         *string `json:"disk_type"`
	DbVendor         string  `json:"database_vendor"`
	DbVersion        string  `json:"database_version"`
	ClusterName      string  `json:"cluster_name"`
	ClusterType      string  `json:"cluster_type"`
	ClusterSize      int64   `json:"cluster_size"`
	HighAvailability bool    `json:"high_availability"`
	Vpc              *struct {
		VpcUUID string `json:"vpc_uuid"`
	} `json:"vpc"`
	Tags []string `json:"tags"`
	AZS  []string `json:"azs"`
}

func ClustersFromLoadAllResponse(r LoadAllResponse) map[string]ccxprov.Cluster {
	c := make(map[string]ccxprov.Cluster)

	for _, info := range r {
		var vpcUUID string
		if info.Vpc != nil {
			vpcUUID = info.Vpc.VpcUUID
		}

		c[info.UUID] = ccxprov.Cluster{
			ID:                info.UUID,
			ClusterName:       info.ClusterName,
			ClusterSize:       info.ClusterSize,
			DBVendor:          info.DbVendor,
			DBVersion:         info.DbVersion,
			ClusterType:       info.ClusterType,
			Tags:              info.Tags,
			CloudProvider:     info.CloudProvider,
			CloudRegion:       info.Region,
			InstanceSize:      info.InstanceSize,
			VolumeType:        pointers.String(info.DiskType),
			VolumeSize:        int64(pointers.Uint64(info.DiskSize)),
			VolumeIOPS:        int64(pointers.Uint64(info.InstanceIOPS)),
			NetworkType:       "", // todo
			HAEnabled:         info.HighAvailability,
			VpcUUID:           vpcUUID,
			AvailabilityZones: info.AZS,
		}
	}

	return c
}

func (cli *Client) Read(_ context.Context, id string) (*ccxprov.Cluster, error) {
	defer cli.mut.Unlock()
	cli.mut.Lock()

	c, ok := cli.clusters[id]
	if ok {
		return &c, nil
	}

	return &ccxprov.Cluster{}, ccxprov.ResourceNotFoundErr
}

func (cli *Client) LoadAll(ctx context.Context) error {
	url := cli.conn.BaseURL + "/api/deployment/api/v1/deployments"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return errors.Join(ccxprov.RequestInitializationErr, err)
	}

	token, err := cli.auth.Auth(ctx)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", token)
	client := &http.Client{Timeout: cli.conn.Timeout}

	res, err := client.Do(req)
	if err != nil {
		return errors.Join(ccxprov.RequestSendingErr, err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: status = %d", chttp.ErrorFromErrorResponse(res.Body), res.StatusCode)
	}

	defer ccxprovio.Close(res.Body)

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.Join(ccxprov.ResponseReadFailedErr, err)
	}

	var rs LoadAllResponse
	if err := json.Unmarshal(b, &rs); err != nil {
		return errors.Join(ccxprov.ResponseDecodingErr, err)
	}

	clusters := ClustersFromLoadAllResponse(rs)

	cli.mut.Lock()
	defer cli.mut.Unlock()

	cli.clusters = clusters

	return nil
}
