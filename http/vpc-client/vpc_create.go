package vpc_client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/severalnines/terraform-provider-ccx/ccx"
	chttp "github.com/severalnines/terraform-provider-ccx/http"
)

type CreateRequest struct {
	Name          string `json:"name"`
	Cloudspace    string `json:"cloudspace"`
	CloudProvider string `json:"cloud"`
	Region        string `json:"region"`
	CidrIpv4Block string `json:"cidr_ipv4_block"`
}

func CreateRequestFromVpc(v ccx.VPC) CreateRequest {
	return CreateRequest{
		Name:          v.Name,
		Cloudspace:    v.CloudSpace,
		CloudProvider: v.CloudProvider,
		Region:        v.Region,
		CidrIpv4Block: v.CidrIpv4Block,
	}
}

type CreateResponse struct {
	VPC *struct {
		ID            string `json:"id"`
		Name          string `json:"name"`
		Cloudspace    string `json:"cloudspace"`
		CloudProvider string `json:"cloud"`
		Region        string `json:"region"`
		CidrIpv4Block string `json:"cidr_ipv4_block"`
	} `json:"vpc"`
}

func VpcFromCreateResponse(r CreateResponse) ccx.VPC {
	if r.VPC == nil {
		return ccx.VPC{}
	}

	return ccx.VPC{
		ID:            r.VPC.ID,
		Name:          r.VPC.Name,
		CloudProvider: r.VPC.CloudProvider,
		CloudSpace:    r.VPC.Cloudspace,
		Region:        r.VPC.Region,
		CidrIpv4Block: r.VPC.CidrIpv4Block,
	}
}

func (cli *Client) Create(ctx context.Context, vpc ccx.VPC) (*ccx.VPC, error) {
	cr := CreateRequestFromVpc(vpc)

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(cr); err != nil {
		return nil, errors.Join(ccx.RequestEncodingErr, err)
	}

	url := cli.conn.BaseURL + "/api/vpc/api/v2/vpcs"
	req, err := http.NewRequest(http.MethodPost, url, &b)
	if err != nil {
		return nil, errors.Join(ccx.RequestInitializationErr, err)
	}

	token, err := cli.auth.Auth(ctx)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", token)
	client := &http.Client{Timeout: cli.conn.Timeout}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Join(ccx.RequestSendingErr, err)
	}

	if res.StatusCode == http.StatusBadRequest {
		return nil, chttp.ErrorFromErrorResponse(res.Body)
	}

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%w: status = %d", ccx.ResponseStatusFailedErr, res.StatusCode)
	}

	var rs CreateResponse
	if err := chttp.DecodeJsonInto(res.Body, &rs); err != nil {
		return nil, err
	}

	newVPC := VpcFromCreateResponse(rs)

	return &newVPC, nil
}
