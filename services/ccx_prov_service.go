package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

type (
	ClusterSpec struct {
		AccountID     string `json:"account_id"`
		ClusterName   string `json:"cluster_name"`
		ClusterType   string `json:"cluster_type"`
		CloudProvider string `json:"cloud_provider"`
		Region        string `json:"region"`
		DbVendor      string `json:"db_vendor"`
		InstanceSize  string `json:"instance_size"`
		InstanceIops  int    `json:"instance_iops"`
		DbAccount     struct {
			DbUsername string `json:"db_username"`
			DbPassword string `json:"db_password"`
			DbHost     string `json:"db_host"`
		} `json:"db_account"`
	}
	DBAccount struct {
		UserName   string `json:"database_username"`
		Password   string `json:"database_password"`
		Host       string `json:"database_host"`
		Database   string `json:"database_database"`
		Privileges string `json:"database_privileges"`
	}
	Cluster struct {
		ClusterUUID             string     `json:"uuid" reform:"cluster_uuid,pk"`
		ControllerID            *string    `json:"controller_id" reform:"controller_uuid"`
		UserID                  string     `json:"account_id" reform:"user_id"`
		ControllerInternalID    int64      `json:"cluster_id,string" reform:"controller_internal_id"`
		ClusterName             string     `json:"cluster_name" reform:"cluster_name"`
		ClusterType             string     `json:"cluster_type" reform:"cluster_type"`
		ClusterRegion           string     `json:"region" reform:"cluster_region"`
		CloudProvider           string     `json:"cloud_provider" reform:"cluster_cloud"`
		ClusterStatus           string     `json:"cluster_status" reform:"cluster_status"`
		ClusterSize             int64      `json:"cluster_size" reform:"cluster_size"`
		ClusterDbVendor         string     `json:"database_vendor" reform:"cluster_db_vendor"`
		ClusterDbVersion        string     `json:"database_version" reform:"cluster_db_version"`
		ClusterDbEndpoint       *string    `json:"database_endpoint" reform:"cluster_db_endpoint"`
		ClusterInstanceSize     string     `json:"instance_size" reform:"cluster_instance_size"`
		ClusterInstanceDiskType *string    `json:"cluster_instance_disk_type" reform:"cluster_instance_disk_type"`
		ClusterInstanceIOPS     *uint64    `json:"iops" reform:"cluster_instance_iops"`
		ClusterInstanceDiskSize *uint64    `json:"disk_size" reform:"cluster_instance_disk_size"`
		HighAvailability        bool       `json:"high_availability" reform:"high_availability"`
		CreatedAt               time.Time  `json:"created,string" reform:"created_at"`
		UpdatedAt               time.Time  `json:"last_updated,string" reform:"updated_at"`
		DeletedAt               *time.Time `json:"deleted_at" reform:"deleted_at"`
		DbAccount               DBAccount  `json:"database_account" reform:"-"`
		Operable                bool       `json:"operable" reform:"-"`
		NotOperableReason       string     `json:"not_operable_reason" reform:"-"`
		VpcUUID                 *string    `json:"vpc_uuid" reform:"vpc_uuid"`
		SubnetUUID              *string    `json:"subnet_uuid" reform:"subnet_uuid"`
		Tags                    []string   `json:"tags" reform:"tags"`
		AZS                     []string   `json:"azs" reform:"azs"`
	}

	CreateClusterRequestV2 struct {
		General struct {
			ClusterName string   `json:"cluster_name"`
			ClusterSize int      `json:"cluster_size"`
			DBVendor    string   `json:"db_vendor"`
			Tags        []string `json:"tags"`
		} `json:"general"`
		Cloud struct {
			CloudProvider string `json:"cloud_provider"`
			CloudRegion   string `json:"cloud_region"`
		} `json:"cloud"`
		Instance struct {
			InstanceSize string `json:"instance_size"` // "Tiny" ... "2X-Large"
			VolumeType   string `json:"volume_type,omitempty"`
			VolumeSize   int    `json:"volume_size,omitempty"`
			VolumeIOPS   string `json:"volume_iops,omitempty"`
		} `json:"instance"`
		Network struct {
			NetworkType string `json:"network_type"` // public/private
			VpcUUID     string `json:"vpc_uuid",omitempty`
		} `json:"network"`
	}
)

func (c *Client) CreateCluster(
	ClusterName string, ClusterSize int, DbVendor string, tags []string,
	CloudRegion string, CloudProvider string, InstanceSize string, volumeType string, volumeSize int,
	volumeIops string, networkType string, vpcUUID string) (*Cluster, error) {
	NewCluster := CreateClusterRequestV2{}
	//general settings
	NewCluster.General.ClusterName = ClusterName
	NewCluster.General.ClusterSize = ClusterSize
	NewCluster.General.DBVendor = DbVendor
	NewCluster.General.Tags = tags
	//Cloud Settings
	NewCluster.Cloud.CloudProvider = CloudProvider
	NewCluster.Cloud.CloudRegion = CloudRegion
	//Instance Settings
	NewCluster.Instance.InstanceSize = InstanceSize
	NewCluster.Instance.VolumeType = volumeType
	NewCluster.Instance.VolumeSize = volumeSize
	//Network Settings
	NewCluster.Network.NetworkType = networkType
	NewCluster.Network.VpcUUID = vpcUUID
	var BaseURLV1 string
	if os.Getenv("ENVIRONMENT") == "dev" {
		BaseURLV1 = ProvServiceUrlDev
	} else if os.Getenv("ENVIRONMENT") == "test" {
		BaseURLV1 = ProvServiceUrlTest
	} else if os.Getenv("ENVIRONMENT") == "prod" {
		BaseURLV1 = ProvServiceUrlProd
	} else {
		BaseURLV1 = ProvServiceUrlProd
	}
	if volumeType == "gp2" || volumeType == "gp3" {
		if volumeIops != "" {
			return nil, errors.New("cannot set iops for volume type gp2|gp3. Please delete the iops parameter and try again")
		}
		NewCluster.Instance.VolumeIOPS = ""
	} else {
		NewCluster.Instance.VolumeIOPS = volumeIops
	}
	if ClusterSize%2 == 0 {
		return nil, fmt.Errorf("cluster size is invalid. Please enter a valid size ( 1 node , 3 nodes , 5 nodes )")
	}
	clusterJSON := new(bytes.Buffer)
	err := json.NewEncoder(clusterJSON).Encode(NewCluster)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}
	req, err := http.NewRequest("POST", BaseURLV1, clusterJSON)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}
	req.AddCookie(c.httpCookie)
	res, err := c.httpClient.Do(req)
	if err != nil || res.StatusCode != 201 {
		dump, _ := httputil.DumpResponse(res, true)
		log.Printf("Received error response from service %s", string(dump))
		return nil, fmt.Errorf("service returned non 200 status code: %s", err)
	}
	defer res.Body.Close()
	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var ServiceResponse Cluster
	if err := json.Unmarshal(responseBody, &ServiceResponse); err != nil {
		return nil, err
	}
	return &ServiceResponse, nil
}
func (c *Client) DeleteCluster(clusterUUID string) error {
	var BaseURLV1 string
	if os.Getenv("ENVIRONMENT") == "dev" {
		BaseURLV1 = ProvServiceUrlDev + "/" + clusterUUID
	} else if os.Getenv("ENVIRONMENT") == "test" {
		BaseURLV1 = ProvServiceUrlTest + "/" + clusterUUID
	} else if os.Getenv("ENVIRONMENT") == "prod" {
		BaseURLV1 = ProvServiceUrlProd + "/" + clusterUUID
	} else {
		BaseURLV1 = ProvServiceUrlProd + "/" + clusterUUID
	}
	req, err := http.NewRequest("DELETE", BaseURLV1, nil)
	if err != nil {
		return err
	}
	req.AddCookie(c.httpCookie)
	res, err := c.httpClient.Do(req)
	if err != nil || res.StatusCode != 200 {
		return fmt.Errorf("error when processing delete request! Please retry!:\t %v", res.Status)
	}
	defer res.Body.Close()
	return nil
}
