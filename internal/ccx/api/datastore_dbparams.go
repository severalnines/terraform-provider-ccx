package api

import (
	"bytes"
	"context"
	"encoding/json"
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

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(cr); err != nil {
		return errors.Join(ccx.RequestEncodingErr, err)
	}

	url := svc.baseURL + "/api/db-configuration/v1/" + storeID

	req, err := http.NewRequest(http.MethodPut, url, &b)
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

	if res.StatusCode == http.StatusBadRequest {
		return lib.ErrorFromErrorResponse(res.Body)
	}

	if res.StatusCode != http.StatusAccepted {
		return fmt.Errorf("%w: status = %d", lib.ErrorFromErrorResponse(res.Body), res.StatusCode)
	}

	status, err := svc.jobs.Await(ctx, storeID, modifyDbConfigJob, svc.timeout)
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
	url := svc.baseURL + "/api/db-configuration/v1/" + storeID

	var rs getParamsResponse

	if err := httpGet(ctx, svc.auth, url, &rs); err != nil {
		return nil, err
	}

	parameters := make(map[string]string, len(rs.Parameters))
	for k, v := range rs.Parameters {
		parameters[k] = v.Value
	}

	return parameters, nil
}
