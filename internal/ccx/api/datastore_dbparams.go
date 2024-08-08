package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
	"github.com/severalnines/terraform-provider-ccx/internal/lib"
)

type setParamsRequest struct {
	Parameters map[string]setParamsParameter `json:"parameters"`
}

type setParamsParameter struct {
	Value string `json:"value"`
}

func (svc *DatastoreService) SetParameters(ctx context.Context, storeID string, parameters map[string]string) error {
	var cr setParamsRequest

	cr.Parameters = make(map[string]setParamsParameter, len(parameters))

	for k, v := range parameters {
		cr.Parameters[k] = setParamsParameter{Value: v}
	}

	res, err := svc.httpcli.Do(ctx, http.MethodPut, "/api/db-configuration/v1/"+storeID, cr)
	if err != nil {
		return errors.Join(ccx.RequestSendingErr, err)
	}

	if res.StatusCode == http.StatusBadRequest {
		return lib.ErrorFromErrorResponse(res.Body)
	}

	if res.StatusCode != http.StatusAccepted {
		return fmt.Errorf("%w: status = %d", lib.ErrorFromErrorResponse(res.Body), res.StatusCode)
	}

	status, err := svc.jobs.Await(ctx, storeID, modifyDbConfigJob)
	if err != nil {
		return fmt.Errorf("awaiting parameters job: %w", err)
	} else if status != jobStatusFinished {
		return fmt.Errorf("parameters job failed: %s", status)
	}

	return nil
}

type getParamsResponse struct {
	Parameters map[string]struct {
		Value             string `json:"value"`
		Description       string `json:"description"`
		Type              string `json:"type"`
		ValidationOptions string `json:"validation_options"`
		DefaultValue      string `json:"default_value"`
	} `json:"parameters"`

	Status    string `json:"status"`
	Error     string `json:"error"`
	UpdatedAt string `json:"updated_at"`
}

func (svc *DatastoreService) GetParameters(ctx context.Context, storeID string) (map[string]string, error) {
	var rs getParamsResponse

	if err := svc.httpcli.Get(ctx, "/api/db-configuration/v1/"+storeID, &rs); err != nil {
		return nil, err
	}

	parameters := make(map[string]string, len(rs.Parameters))
	for k, v := range rs.Parameters {
		parameters[k] = v.Value
	}

	return parameters, nil
}
