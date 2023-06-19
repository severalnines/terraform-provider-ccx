package cluster_client

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	ccxprov "github.com/severalnines/terraform-provider-ccx"
)

func (cli *Client) Delete(ctx context.Context, id string) error {
	url := cli.conn.BaseURL + "/api/prov/api/v2/cluster" + "/" + id
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return errors.Join(ccxprov.RequestInitializationErr, err)
	}

	token, err := cli.auth.Auth(ctx)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", token)
	client := &http.Client{Timeout: cli.conn.Timeout}

	res, err := client.Do(req)
	if err != nil {
		return errors.Join(ccxprov.RequestSendingErr, err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: status = %d", ccxprov.ResponseStatusFailedErr, res.StatusCode)
	}

	if err := cli.LoadAll(ctx); err != nil {
		return errors.Join(ccxprov.ResourcesLoadFailedErr, err)
	}

	return nil
}
