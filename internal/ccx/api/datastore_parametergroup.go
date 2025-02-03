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

	_, err := svc.client.Do(ctx, http.MethodPut, "/api/db-configuration/v1/parameter-groups/apply/"+group+"/"+id, nil)
	if err != nil {
		return fmt.Errorf("applying parameter group: %w", err)
	}

	status, err := svc.jobs.Await(ctx, id, ccx.ModifyDbConfigJob)
	if err != nil {
		return fmt.Errorf("awaiting modify parameter job: %w", err)
	} else if status != ccx.JobStatusFinished {
		return fmt.Errorf("modify parameter job failed: %s", status)
	}

	return nil
}
