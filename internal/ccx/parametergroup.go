package ccx

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type ParameterGroupsClient struct {
	client HTTPClient
}

var _ ParameterGroupsService = (*ParameterGroupsClient)(nil)

// NewParameterGroupsClient creates a new ParameterGroupService
func NewParameterGroupsClient(client HTTPClient) *ParameterGroupsClient {
	c := ParameterGroupsClient{
		client: client,
	}

	return &c
}

type createParameterGroupRequest struct {
	Name            string `json:"name"`
	DatabaseVendor  string `json:"database_vendor"`
	DatabaseVersion string `json:"database_version"`
	DatabaseType    string `json:"database_type"`
	Description     string `json:"description"`

	Parameters map[string]string `json:"parameters"`
}

type createParameterGroupResponse struct {
	ID string `json:"uuid"`
}

func (svc *ParameterGroupsClient) Create(ctx context.Context, p ParameterGroup) (*ParameterGroup, error) {
	req := createParameterGroupRequest{
		Name:            p.Name,
		DatabaseVendor:  p.DatabaseVendor,
		DatabaseVersion: p.DatabaseVersion,
		DatabaseType:    p.DatabaseType,
		Parameters:      p.DbParameters,
		Description:     p.Description,
	}

	res, err := svc.client.Do(ctx, http.MethodPost, "/api/db-configuration/v1/parameter-groups", req)

	if err != nil {
		return nil, err
	}

	var rs createParameterGroupResponse
	if err := DecodeJsonInto(res.Body, &rs); err != nil {
		return nil, fmt.Errorf("creating parameter group: %w", err)
	}

	p.ID = rs.ID

	return &p, nil
}

func (svc *ParameterGroupsClient) Read(ctx context.Context, id string) (*ParameterGroup, error) {
	var rs ParameterGroup

	err := svc.client.Get(ctx, "/api/db-configuration/v1/parameter-groups/"+id, &rs)
	if err != nil {
		return nil, err
	}

	return &rs, nil
}

type updateParameterGroupRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Parameters  map[string]string `json:"parameters"`
}

func (svc *ParameterGroupsClient) Update(ctx context.Context, p ParameterGroup) error {
	req := updateParameterGroupRequest{
		Name:        p.Name,
		Description: p.Description,
		Parameters:  p.DbParameters,
	}

	_, err := svc.client.Do(ctx, http.MethodPut, "/api/db-configuration/v1/parameter-groups/"+p.ID+"?sync=true", req)
	if err != nil {
		return err
	}

	return nil
}

func (svc *ParameterGroupsClient) Delete(ctx context.Context, id string) error {
	_, err := svc.client.Do(ctx, http.MethodDelete, "/api/db-configuration/v1/parameter-groups/"+id, nil)
	if errors.Is(err, ErrResourceNotFound) {
		tflog.Warn(ctx, "deleting parameter group: not found", map[string]interface{}{"id": id})
		return nil
	} else if err != nil {
		return fmt.Errorf("deleting parameter group: %w", err)
	}

	return nil
}
