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
	ProvServiceUrl = "https://ccx-prov-service.s9s-dev.net/api/v1/cluster/"
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
)

func (c *Client) CreateCluster(ClusterName string,
	ClusterType string, CloudProvider string,
	Region string, DbVendor string,
	InstanceSize string, InstanceIops int,
	DbUsername string, DbPassword string,
	DbHost string) (Cluster, error) {
	AccountID := c.userId
	NewCluster := ClusterSpec{}
	NewCluster.AccountID = AccountID
	NewCluster.ClusterName = ClusterName
	NewCluster.ClusterType = ClusterType
	NewCluster.CloudProvider = CloudProvider
	NewCluster.Region = Region
	NewCluster.DbVendor = DbVendor
	NewCluster.InstanceSize = InstanceSize
	NewCluster.InstanceIops = InstanceIops
	NewCluster.DbAccount.DbUsername = DbUsername
	NewCluster.DbAccount.DbPassword = DbPassword
	NewCluster.DbAccount.DbHost = DbHost
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
