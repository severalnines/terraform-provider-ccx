package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/severalnines/terraform-provider-ccx/ccx"
)

func (svc *DatastoreService) Delete(ctx context.Context, id string) error {
	url := svc.baseURL + "/api/prov/api/v2/cluster" + "/" + id
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return errors.Join(ccx.RequestInitializationErr, err)
	}

	token, err := svc.auth.Auth(ctx)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", token)
	client := &http.Client{Timeout: ccx.DefaultTimeout}

	res, err := client.Do(req)
	if err != nil {
		return errors.Join(ccx.RequestSendingErr, err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: status = %d", ccx.ResponseStatusFailedErr, res.StatusCode)
	}

	status, err := svc.jobs.Await(ctx, id, deployStoreJob, svc.timeout)
	if err != nil {
		return fmt.Errorf("awaiting destroy job: %w", err)
	} else if status != jobStatusFinished {
		return fmt.Errorf("destroy job failed: %s", status)
	}

	if err := svc.LoadAll(ctx); err != nil {
		return errors.Join(ccx.ResourcesLoadFailedErr, err)
	}

	return nil
}
