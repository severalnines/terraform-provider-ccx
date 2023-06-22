package datastore_client

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/severalnines/terraform-provider-ccx/ccx"
)

func (cli *Client) Delete(ctx context.Context, id string) error {
	url := cli.conn.BaseURL + "/api/prov/api/v2/cluster" + "/" + id
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return errors.Join(ccx.RequestInitializationErr, err)
	}

	token, err := cli.auth.Auth(ctx)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", token)
	client := &http.Client{Timeout: cli.conn.Timeout}

	res, err := client.Do(req)
	if err != nil {
		return errors.Join(ccx.RequestSendingErr, err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: status = %d", ccx.ResponseStatusFailedErr, res.StatusCode)
	}

	if err := cli.LoadAll(ctx); err != nil {
		return errors.Join(ccx.ResourcesLoadFailedErr, err)
	}

	return nil
}
