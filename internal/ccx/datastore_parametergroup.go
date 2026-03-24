package ccx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

func (svc *DatastoresClient) ApplyParameterGroup(ctx context.Context, id, groupID string) error {
	if groupID == "" {
		return errors.New("group ID is required")
	}

	_, err := svc.client.Do(ctx, http.MethodPut, "/api/db-configuration/v1/parameter-groups/apply/"+groupID+"/"+id, nil)
	if err != nil {
		return fmt.Errorf("applying parameter group: %w", err)
	}

	status, err := svc.jobs.Await(ctx, id, ModifyDbConfigJob)
	if err != nil {
		return fmt.Errorf("awaiting modify parameter job: %w", err)
	} else if status != JobStatusFinished {
		return fmt.Errorf("modify parameter job failed: %s", status)
	}

	return nil
}
