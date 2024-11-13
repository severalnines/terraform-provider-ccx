package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
)

func (svc *DatastoreService) ApplyParameterGroup(ctx context.Context, id, group string) error {
	if group == "" {
		return errors.New("group name is required")
	}

	rs, err := svc.client.Do(ctx, http.MethodPut, "/api/db-configuration/v1/parameter-groups/apply/"+group+"/"+id, nil)
	if err != nil {
		return errors.Join(ccx.RequestSendingErr, err)
	}

	if rs.StatusCode != http.StatusAccepted {
		return errors.New("failed to apply parameter group")
	}

	status, err := svc.jobs.Await(ctx, id, modifyDbConfigJob)
	if err != nil {
		return fmt.Errorf("%w: awaiting modify parameter job: %w", ccx.CreateFailedErr, err)
	} else if status != jobStatusFinished {
		return fmt.Errorf("%w: modify parameter job failed: %s", ccx.CreateFailedErr, status)
	}

	return nil
}
