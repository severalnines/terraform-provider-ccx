package vpc_client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/severalnines/terraform-provider-ccx/ccx"
	ccxprovio "github.com/severalnines/terraform-provider-ccx/io"
)

type ReadResponse CreateResponse

func VpcFromReadResponse(r ReadResponse) ccx.VPC {
	return VpcFromCreateResponse(CreateResponse(r))
}

func (cli *Client) Read(ctx context.Context, id string) (*ccx.VPC, error) {
	url := cli.conn.BaseURL + "/api/vpc/api/v2/vpcs/" + id

	req, err := http.NewRequest(http.MethodGet, url, nil)
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

	if res.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%w: %s", ccx.ResourceNotFoundErr, id)
	} else if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status = %d", ccx.ResponseStatusFailedErr, res.StatusCode)
	}

	defer ccxprovio.Close(res.Body)

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Join(ccx.ResponseReadFailedErr, err)
	}

	var rs ReadResponse
	if err := json.Unmarshal(b, &rs); err != nil {
		return nil, errors.Join(ccx.ResponseDecodingErr, err)
	}

	vpc := VpcFromReadResponse(rs)

	return &vpc, nil
}
