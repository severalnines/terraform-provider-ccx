package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/severalnines/terraform-provider-ccx/ccx"
	"github.com/severalnines/terraform-provider-ccx/internal/lib"
)

type createVpcRequest struct {
	Name          string `json:"name"`
	Cloudspace    string `json:"cloudspace"`
	CloudProvider string `json:"cloud"`
	Region        string `json:"region"`
	CidrIpv4Block string `json:"cidr_ipv4_block"`
}

func CreateRequestFromVpc(v ccx.VPC) createVpcRequest {
	return createVpcRequest{
		Name:          v.Name,
		Cloudspace:    v.CloudSpace,
		CloudProvider: v.CloudProvider,
		Region:        v.Region,
		CidrIpv4Block: v.CidrIpv4Block,
	}
}

func (svc *VpcService) Create(ctx context.Context, vpc ccx.VPC) (*ccx.VPC, error) {
	cr := CreateRequestFromVpc(vpc)

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(cr); err != nil {
		return nil, errors.Join(ccx.RequestEncodingErr, err)
	}

	url := svc.baseURL + "/api/vpc/api/v2/vpcs"
	req, err := http.NewRequest(http.MethodPost, url, &b)
	if err != nil {
		return nil, errors.Join(ccx.RequestInitializationErr, err)
	}

	token, err := svc.auth.Auth(ctx)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", token)
	client := &http.Client{Timeout: ccx.DefaultTimeout}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Join(ccx.RequestSendingErr, err)
	}

	if res.StatusCode == http.StatusBadRequest {
		return nil, lib.ErrorFromErrorResponse(res.Body)
	}

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%w: status = %d", ccx.ResponseStatusFailedErr, res.StatusCode)
	}

	var rs vpcResponse
	if err := lib.DecodeJsonInto(res.Body, &rs); err != nil {
		return nil, err
	}

	newVPC := vpcFromResponse(rs)

	return &newVPC, nil
}
