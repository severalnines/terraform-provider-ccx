package services

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/lensesio/tableprinter"
)

type DeploymentServiceResponse []struct {
	AccountID         string `json:"account_id"`
	UUID              string `json:"uuid"`
	Region            string `json:"region"`
	CloudProvider     string `json:"cloud_provider"`
	InstanceSize      string `json:"instance_size"`
	InstanceIops      int    `json:"instance_iops"`
	DatabaseVendor    string `json:"database_vendor"`
	DatabaseVersion   string `json:"database_version"`
	DatabaseEndpoint  string `json:"database_endpoint"`
	ClusterName       string `json:"cluster_name"`
	ClusterStatus     string `json:"cluster_status"`
	ClusterStatusText string `json:"cluster_status_text"`
	ClusterType       string `json:"cluster_type"`
	ClusterID         int    `json:"cluster_id"`
	ClusterSize       int    `json:"cluster_size"`
	Operable          bool   `json:"operable"`
	NotOperableReason string `json:"not_operable_reason"`
	SslEnabled        bool   `json:"ssl_enabled"`
	Vpc               struct {
		VpcUUID       string `json:"vpc_uuid"`
		UserID        string `json:"user_id"`
		CloudProvider string `json:"cloud_provider"`
		VpcName       string `json:"vpc_name"`
		VpcData       struct {
			VpcID         string `json:"vpc_id"`
			Type          string `json:"type"`
			Cloud         string `json:"cloud"`
			Region        string `json:"region"`
			CidrIpv4Block string `json:"cidr_ipv4_block"`
			State         string `json:"state"`
			Aws           struct {
				VpcID              string `json:"vpc_id"`
				IgwID              string `json:"igw_id"`
				MainRouteTableID   string `json:"main_route_table_id"`
				PublicRouteTableID string `json:"public_route_table_id"`
				Subnets            []struct {
					SubnetID                string `json:"subnet_id"`
					VpcID                   string `json:"vpc_id"`
					AvailabilityZone        string `json:"availability_zone"`
					AvailabilityZoneID      string `json:"availability_zone_id"`
					AvailableIPAddressCount int    `json:"available_ip_address_count"`
					CidrIpv4Block           string `json:"cidr_ipv4_block"`
					DefaultForAz            bool   `json:"default_for_az"`
					State                   string `json:"state"`
				} `json:"subnets"`
			} `json:"aws"`
		} `json:"vpc_data"`
		CanBeDeleted       bool      `json:"can_be_deleted"`
		CannotDeleteReason string    `json:"cannot_delete_reason"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
	} `json:"vpc"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ClusterDetailResponse struct {
	UUID          string `json:"uuid"`
	DatabaseNodes []struct {
		HostUUID         string      `json:"host_uuid"`
		ClusterUUID      string      `json:"cluster_uuid"`
		CloudProvider    string      `json:"cloud_provider"`
		HostType         interface{} `json:"host_type"`
		InstanceID       string      `json:"instance_id"`
		InstanceType     string      `json:"instance_type"`
		Region           string      `json:"region"`
		PublicIP         interface{} `json:"public_ip"`
		PrivateIP        string      `json:"private_ip"`
		Fqdn             interface{} `json:"fqdn"`
		Disks            string      `json:"disks"`
		Vpc              string      `json:"vpc"`
		Subnet           string      `json:"subnet"`
		SecurityGroup    string      `json:"security_group"`
		CreatedAt        time.Time   `json:"created_at"`
		UpdatedAt        time.Time   `json:"updated_at"`
		ServiceStartedAt time.Time   `json:"service_started_at"`
		UniqueID         int         `json:"unique_id"`
		LastSeenAt       time.Time   `json:"last_seen_at"`
		AdminPort        int         `json:"admin_port"`
		Port             int         `json:"port"`
		ListeningPort    int         `json:"listening_port"`
		ListeningPortRo  int         `json:"listening_port_ro"`
		Hostname         string      `json:"hostname"`
		HostStatus       string      `json:"host_status"`
		Role             string      `json:"role"`
		NodeType         string      `json:"node_type"`
		SslEnabled       bool        `json:"ssl_enabled"`
	} `json:"database_nodes"`
	LoadBalancerNodes []struct {
		HostUUID         string      `json:"host_uuid"`
		ClusterUUID      string      `json:"cluster_uuid"`
		CloudProvider    string      `json:"cloud_provider"`
		HostType         interface{} `json:"host_type"`
		InstanceID       string      `json:"instance_id"`
		InstanceType     string      `json:"instance_type"`
		Region           string      `json:"region"`
		PublicIP         interface{} `json:"public_ip"`
		PrivateIP        string      `json:"private_ip"`
		Fqdn             interface{} `json:"fqdn"`
		Disks            string      `json:"disks"`
		Vpc              string      `json:"vpc"`
		Subnet           string      `json:"subnet"`
		SecurityGroup    string      `json:"security_group"`
		CreatedAt        time.Time   `json:"created_at"`
		UpdatedAt        time.Time   `json:"updated_at"`
		ServiceStartedAt time.Time   `json:"service_started_at"`
		UniqueID         int         `json:"unique_id"`
		LastSeenAt       time.Time   `json:"last_seen_at"`
		AdminPort        int         `json:"admin_port"`
		Port             int         `json:"port"`
		ListeningPort    int         `json:"listening_port"`
		ListeningPortRo  int         `json:"listening_port_ro"`
		Hostname         string      `json:"hostname"`
		HostStatus       string      `json:"host_status"`
		Role             string      `json:"role"`
		NodeType         string      `json:"node_type"`
		SslEnabled       bool        `json:"ssl_enabled"`
	} `json:"load_balancer_nodes"`
}

type ClusterTableHeaders struct {
	Name             string `header:"cluster name"`
	Status           string `header:"status"`
	Uuid             string `header:"UUID"`
	Databasevendor   string `header:"Vendor"`
	Databaseversion  string `header:"Database Version"`
	Databaseendpoint string `header:"Database Endpoint"`
}

type ClusterDetailHeaders struct {
	Uuid      string      `header:"UUID"`
	Hostname  string      `header:"Hostname"`
	Vpc       string      `header:"VPC"`
	PublicIP  interface{} `header:"Public IP"`
	PrivateIP string      `header:"Private IP"`
	Status    string      `header:"Status"`
	Port      int         `header:"Port"`
	Role      string      `header:"Role"`
	NodeType  string      `header:"Type"`
}

func GetClusters(userId string, cookie *http.Cookie) {
	BaseURLV1 := "https://ccx-deployment-service.s9s-dev.net/api/v1/deployments"
	req, _ := http.NewRequest("GET", BaseURLV1, nil)
	req.AddCookie(cookie)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	if res.StatusCode != 200 {
		log.Fatal(res.Status)
	}
	defer res.Body.Close()
	printer := tableprinter.New(os.Stdout)
	responseBody, _ := ioutil.ReadAll(res.Body)
	var ServiceResponse DeploymentServiceResponse
	var table []ClusterTableHeaders
	json.Unmarshal(responseBody, &ServiceResponse)
	for i := range ServiceResponse {
		table = append(table,
			ClusterTableHeaders{ServiceResponse[i].ClusterName,
				ServiceResponse[i].ClusterStatus,
				ServiceResponse[i].UUID,
				ServiceResponse[i].DatabaseVendor,
				ServiceResponse[i].DatabaseVersion,
				ServiceResponse[i].DatabaseEndpoint})
	}
	printer.Print(table)
}

func GetClusterByID(uuid string, cookie *http.Cookie) {
	BaseURLV1 := "https://ccx-deployment-service.s9s-dev.net/api/v1/deployment/" + uuid
	req, _ := http.NewRequest("GET", BaseURLV1, nil)
	req.AddCookie(cookie)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	if res.StatusCode != 200 {
		log.Fatal(res.Status)
	}
	defer res.Body.Close()
	printer := tableprinter.New(os.Stdout)
	responseBody, _ := ioutil.ReadAll(res.Body)
	var ServiceResponse ClusterDetailResponse
	var table []ClusterDetailHeaders
	json.Unmarshal(responseBody, &ServiceResponse)
	for i := range ServiceResponse.DatabaseNodes {
		table = append(table,
			ClusterDetailHeaders{ServiceResponse.UUID,
				ServiceResponse.DatabaseNodes[i].Hostname,
				ServiceResponse.DatabaseNodes[i].Vpc,
				ServiceResponse.DatabaseNodes[i].PublicIP,
				ServiceResponse.DatabaseNodes[i].PrivateIP,
				ServiceResponse.DatabaseNodes[i].HostStatus,
				ServiceResponse.DatabaseNodes[i].Port,
				ServiceResponse.DatabaseNodes[i].Role,
				ServiceResponse.DatabaseNodes[i].NodeType,
			})
	}
	printer.Print(table)
}
