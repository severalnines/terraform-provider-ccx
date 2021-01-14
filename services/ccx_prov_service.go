package services

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	ProvServiceUrl = "https://ccx-prov-service.s9s-dev.net/api/v1/cluster"
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
)

func (c *Client) CreateCluster(ClusterName string,
	ClusterType string, CloudProvider string,
	Region string, DbVendor string,
	InstanceSize string, InstanceIops int,
	DbUsername string, DbPassword string,
	DbHost string) error {
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
	if err != nil {
		return err
	}
	if res.StatusCode != 201 {
		log.Fatalln(res.Status)
	}
	defer res.Body.Close()
	responseBody, _ := ioutil.ReadAll(res.Body)
	var ServiceResponse DeploymentServiceResponse
	json.Unmarshal(responseBody, &ServiceResponse)
	log.Print(ServiceResponse)
	return nil
}
