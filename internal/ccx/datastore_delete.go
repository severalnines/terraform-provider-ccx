package ccx

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (svc *DatastoresClient) Delete(ctx context.Context, id string) error {
	_, err := svc.client.Do(ctx, http.MethodDelete, "/api/prov/api/v2/cluster"+"/"+id, nil)
	if errors.Is(err, ErrResourceNotFound) {
		tflog.Warn(ctx, "deleting datastore: not found", map[string]interface{}{"id": id})
		return nil
	} else if err != nil {
		return fmt.Errorf("deleting datastore: %w", err)
	}

	status, err := svc.jobs.Await(ctx, id, DestroyStoreJob)
	if err != nil {
		return fmt.Errorf("awaiting destroy job: %w", err)
	} else if status != JobStatusFinished {
		return fmt.Errorf("destroy job failed: %s", status)
	}

	return nil
}
