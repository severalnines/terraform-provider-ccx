package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
)

func (svc *DatastoreService) Delete(ctx context.Context, id string) error {
	res, err := svc.client.Do(ctx, http.MethodDelete, "/api/prov/api/v2/cluster"+"/"+id, nil)
	if err != nil {
		return errors.Join(ccx.RequestSendingErr, err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: status = %d", ccx.ResponseStatusFailedErr, res.StatusCode)
	}

	status, err := svc.jobs.Await(ctx, id, deployStoreJob)
	if err != nil {
		return fmt.Errorf("awaiting destroy job: %w", err)
	} else if status != jobStatusFinished {
		return fmt.Errorf("destroy job failed: %s", status)
	}

	return nil
}
