package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

// CreateVpcRequest struct.
type (
	CreateVpcRequest struct {
		CloudProvider string `json:"cloud"`
		Region        string `json:"region"`
		CidrIpv4Block string `json:"cidr_ipv4_block,omitempty"`
		VpcName       string `json:"vpc_name"`
	}
	CreateVpcResponse struct {
		Cloud            string `json:"cloud"`
		Region           string `json:"region"`
		CidrIpv4Block    string `json:"cidr_ipv4_block"`
		AvailabilityZone string `json:"availability_zone"`
		VpcUUID          string `json:"vpc_uuid"`
		VpcName          string `json:"vpc_name"`
	}
)

func (c *Client) CreateVpc(VpcName string, VpcCloudProvider string, VpcRegion string, VpcCidrIpv4Block string) (*CreateVpcResponse, error) {
	log.Printf("Start create vpc")
	NewVPC := &CreateVpcRequest{}
	NewVPC.VpcName = VpcName
	NewVPC.CloudProvider = VpcCloudProvider
	NewVPC.Region = VpcRegion
	NewVPC.CidrIpv4Block = VpcCidrIpv4Block
	vpcJSON := new(bytes.Buffer)
	log.Printf("NewVPC")
	var BaseURLV1 string
	json.NewEncoder(vpcJSON).Encode(NewVPC)
	if os.Getenv("ENVIRONMENT") == "dev" {
		BaseURLV1 = VpcServiceUrlDev
	} else if os.Getenv("ENVIRONMENT") == "test" {
		BaseURLV1 = VpcServiceUrlTest
	} else if os.Getenv("ENVIRONMENT") == "prod" {
		BaseURLV1 = VpcServiceUrlProd
	} else {
		BaseURLV1 = VpcServiceUrlProd
	}
	log.Println(BaseURLV1)
	req, _ := http.NewRequest("POST", BaseURLV1, vpcJSON)
	req.AddCookie(c.httpCookie)
	res, _ := c.httpClient.Do(req)
	if res.StatusCode != 201 {
		dump, _ := httputil.DumpResponse(res, true)
		return nil, fmt.Errorf("service returned non 200 status code: %s", string(dump))
	}
	defer res.Body.Close()
	responseBody, _ := ioutil.ReadAll(res.Body)
	var ServiceResponse CreateVpcResponse
	dump, _ := httputil.DumpResponse(res, true)
	log.Println(string(dump))
	json.Unmarshal(responseBody, &ServiceResponse)
	return &ServiceResponse, nil
}
func (c *Client) GetVPCbyUUID(uuid string) error {
	var BaseURLV1 string
	if os.Getenv("ENVIRONMENT") == "dev" {
		BaseURLV1 = VpcServiceUrlDev + "/" + uuid
	} else if os.Getenv("ENVIRONMENT") == "test" {
		BaseURLV1 = VpcServiceUrlTest + "/" + uuid
	} else if os.Getenv("ENVIRONMENT") == "prod" {
		BaseURLV1 = VpcServiceUrlProd + "/" + uuid
	} else {
		BaseURLV1 = VpcServiceUrlProd + "/" + uuid
	}
	req, _ := http.NewRequest("GET", BaseURLV1, nil)
	req.AddCookie(c.httpCookie)
	res, err := c.httpClient.Do(req)
	dump, _ := httputil.DumpResponse(res, true)
	log.Printf(string(dump))
	if err != nil {
		log.Fatal("CCX_VPC_SERVICE: Error!")
	}
	if res.StatusCode != 200 {
		log.Printf("CCX_VPC_SERVICE: Error! %v", res.StatusCode)
	}
	defer res.Body.Close()
	responseBody, _ := ioutil.ReadAll(res.Body)
	var ServiceResponse ClusterDetailResponse
	json.Unmarshal(responseBody, &ServiceResponse)
	log.Println(ServiceResponse)
	return nil
}
func (c *Client) DeleteVPCbyUUID(uuid string) error {
	var BaseURLV1 string
	if os.Getenv("ENVIRONMENT") == "dev" {
		BaseURLV1 = VpcServiceUrlDev + "/" + uuid
	} else if os.Getenv("ENVIRONMENT") == "test" {
		BaseURLV1 = VpcServiceUrlTest + "/" + uuid
	} else if os.Getenv("ENVIRONMENT") == "prod" {
		BaseURLV1 = VpcServiceUrlProd + "/" + uuid
	} else {
		BaseURLV1 = VpcServiceUrlProd + "/" + uuid
	}
	req, _ := http.NewRequest("GET", BaseURLV1, nil)
	req.AddCookie(c.httpCookie)
	res, err := c.httpClient.Do(req)
	dump, _ := httputil.DumpResponse(res, true)
	log.Printf(string(dump))
	if err != nil {
		log.Fatal("CCX_VPC_SERVICE: Error!")
	}
	if res.StatusCode != 200 {
		log.Printf("CCX_VPC_SERVICE: Error! %v", res.StatusCode)
	}
	defer res.Body.Close()
	responseBody, _ := ioutil.ReadAll(res.Body)
	var ServiceResponse ClusterDetailResponse
	json.Unmarshal(responseBody, &ServiceResponse)
	log.Println(ServiceResponse)
	return nil
}
