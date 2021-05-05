package services

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	ProvServiceUrl = "https://ccx-prov-service.s9s-dev.net/api/v2/cluster/"
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
		ClusterUUID          string     `json:"uuid" reform:"cluster_uuid,pk"`
		ControllerID         *string    `json:"controller_id" reform:"controller_uuid"`
		UserID               string     `json:"account_id" reform:"user_id"`
		ControllerInternalID int64      `json:"cluster_id,string" reform:"controller_internal_id"`
		ClusterName          string     `json:"cluster_name" reform:"cluster_name"`
		ClusterType          string     `json:"cluster_type" reform:"cluster_type"`
		ClusterRegion        string     `json:"region" reform:"cluster_region"`
		CloudProvider        string     `json:"cloud_provider" reform:"cluster_cloud"`
		ClusterStatus        string     `json:"cluster_status" reform:"cluster_status"`
		ClusterSize          int64      `json:"cluster_size" reform:"cluster_size"`
		ClusterDbVendor      string     `json:"database_vendor" reform:"cluster_db_vendor"`
		ClusterDbVersion     string     `json:"database_version" reform:"cluster_db_version"`
		ClusterDbEndpoint    *string    `json:"database_endpoint" reform:"cluster_db_endpoint"`
		ClusterInstanceSize  string     `json:"instance_size" reform:"cluster_instance_size"`
		ClusterInstanceIOPS  int64      `json:"iops" reform:"cluster_instance_iops"`
		CreatedAt            time.Time  `json:"created,string" reform:"created_at"`
		UpdatedAt            time.Time  `json:"last_updated,string" reform:"updated_at"`
		DeletedAt            *time.Time `json:"deleted_at" reform:"deleted_at"`
		DbAccount            DBAccount  `json:"database_account" reform:"-"`
		Operable             bool       `json:"operable" reform:"-"`
		NotOperableReason    string     `json:"not_operable_reason" reform:"-"`
		VpcUUID              *string    `json:"vpc_uuid" reform:"vpc_uuid"`
		SubnetUUID           *string    `json:"subnet_uuid" reform:"subnet_uuid"`
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
			VolumeType   string `json:"volume_type"`
			VolumeSize   int    `json:"volume_size"`
			VolumeIOPS   int    `json:"volume_iops"`
		} `json:"instance"`
		Network struct {
			NetworkType       string   `json:"network_type"` // public/private
			HAEnabled         bool     `json:"ha_enabled"`
			VpcUUID           string   `json:"vpc_uuid"`
			AvailabilityZones []string `json:"availability_zones"`
		} `json:"network"`
	}
)

func (c *Client) CreateCluster(
	ClusterName string, ClusterSize int, DbVendor string, tags []string,
	CloudRegion string, CloudProvider string, InstanceSize string, volumeType string, volumeSize int,
	volumeIops int, networkType string, networkHAEnabled bool, VPCUuid string,
	vpcAz []string) (Cluster, error) {
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
	NewCluster.Instance.VolumeIOPS = volumeIops
	//Network Settings
	NewCluster.Network.NetworkType = networkType
	NewCluster.Network.HAEnabled = networkHAEnabled
	NewCluster.Network.VpcUUID = VPCUuid
	NewCluster.Network.AvailabilityZones = vpcAz

	clusterJSON := new(bytes.Buffer)
	json.NewEncoder(clusterJSON).Encode(NewCluster)
	req, _ := http.NewRequest("POST", ProvServiceUrl, clusterJSON)
	req.AddCookie(c.httpCookie)
	res, err := c.httpClient.Do(req)
	log.Println("Response done!")
	if err != nil {
		log.Println(err)
	}
	if res.StatusCode != 201 {
		log.Fatalln(res.Status)
	}
	defer res.Body.Close()
	responseBody, _ := ioutil.ReadAll(res.Body)
	var ServiceResponse Cluster
	json.Unmarshal(responseBody, &ServiceResponse)
	log.Println(ServiceResponse)
	return ServiceResponse, nil
}
func (c *Client) DeleteCluster(clusterUUID string) error {
	req, _ := http.NewRequest("DELETE", ProvServiceUrl+clusterUUID, nil)
	req.AddCookie(c.httpCookie)
	res, err := c.httpClient.Do(req)
	log.Println("Response done!")
	if err != nil {
		log.Println(err)
	}
	if res.StatusCode != 201 {
		log.Println(res.Status)
	}
	defer res.Body.Close()
	return nil
}
