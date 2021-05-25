package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
)

const (
	VpcServiceUrl = "https://ccx.s9s-dev.net/api/vpc/api/v2/vpc"
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
	json.NewEncoder(vpcJSON).Encode(NewVPC)
	req, _ := http.NewRequest("POST", VpcServiceUrl, vpcJSON)
	req.AddCookie(c.httpCookie)
	res, _ := c.httpClient.Do(req)
	if res.StatusCode != 201 {
		dump, _ := httputil.DumpResponse(res, true)
		return nil, fmt.Errorf("service returned non 200 status code: %s", string(dump))
	}
	defer res.Body.Close()
	responseBody, _ := ioutil.ReadAll(res.Body)
	log.Printf(string(responseBody))
	var ServiceResponse CreateVpcResponse
	dump, _ := httputil.DumpResponse(res, true)
	log.Printf(string(dump))
	json.Unmarshal(responseBody, &ServiceResponse)
	return &ServiceResponse, nil
}
func (c *Client) GetVPCbyUUID(uuid string) error {
	BaseURLV1 := VpcServiceUrl + uuid
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
	BaseURLV1 := VpcServiceUrl + uuid
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
