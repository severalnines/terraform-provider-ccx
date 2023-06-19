package vpc_client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	ccxprov "github.com/severalnines/terraform-provider-ccx"
	ccxprovio "github.com/severalnines/terraform-provider-ccx/io"
)

type ReadResponse CreateResponse

func VpcFromReadResponse(r ReadResponse) ccxprov.VPC {
	return VpcFromCreateResponse(CreateResponse(r))
}

func (cli *Client) Read(ctx context.Context, id string) (*ccxprov.VPC, error) {
	url := cli.conn.BaseURL + "/api/vpc/api/v2/vpcs/" + id

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Join(ccxprov.RequestInitializationErr, err)
	}

	token, err := cli.auth.Auth(ctx)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", token)
	client := &http.Client{Timeout: cli.conn.Timeout}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Join(ccxprov.RequestSendingErr, err)
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%w: resource not found: %s", ccxprov.ResourceNotFoundErr, id)
	} else if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status = %d", ccxprov.ResponseStatusFailedErr, res.StatusCode)
	}

	defer ccxprovio.Close(res.Body)

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Join(ccxprov.ResponseReadFailedErr, err)
	}

	var rs ReadResponse
	if err := json.Unmarshal(b, &rs); err != nil {
		return nil, errors.Join(ccxprov.ResponseDecodingErr, err)
	}

	vpc := VpcFromReadResponse(rs)

	return &vpc, nil
}
