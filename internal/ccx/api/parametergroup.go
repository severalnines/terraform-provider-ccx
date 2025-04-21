package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
	"github.com/severalnines/terraform-provider-ccx/internal/lib"
)

type ParameterGroupService struct {
	client HttpClient
}

var _ ccx.ParameterGroupService = (*ParameterGroupService)(nil)

// ParameterGroups creates a new ParameterGroupService
func ParameterGroups(client HttpClient) *ParameterGroupService {
	c := ParameterGroupService{
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

func (svc *ParameterGroupService) Create(ctx context.Context, p ccx.ParameterGroup) (*ccx.ParameterGroup, error) {
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
	if err := lib.DecodeJsonInto(res.Body, &rs); err != nil {
		return nil, fmt.Errorf("creating parameter group: %w", err)
	}

	p.ID = rs.ID

	return &p, nil
}

func (svc *ParameterGroupService) Read(ctx context.Context, id string) (*ccx.ParameterGroup, error) {
	var rs ccx.ParameterGroup

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

func (svc *ParameterGroupService) Update(ctx context.Context, p ccx.ParameterGroup) error {
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

func (svc *ParameterGroupService) Delete(ctx context.Context, id string) error {
	_, err := svc.client.Do(ctx, http.MethodDelete, "/api/db-configuration/v1/parameter-groups/"+id, nil)
	if errors.Is(err, ccx.ErrResourceNotFound) {
		tflog.Warn(ctx, "deleting parameter group: not found", map[string]interface{}{"id": id})
		return nil
	} else if err != nil {
		return fmt.Errorf("deleting parameter group: %w", err)
	}

	return nil
}
