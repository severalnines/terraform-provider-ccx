package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/severalnines/terraform-provider-ccx/ccx"
	"github.com/severalnines/terraform-provider-ccx/internal/lib"
)

type UpdateRequest struct {
	NewName       string `json:"cluster_name"`
	NewVolumeSize uint   `json:"new_volume_size"`
}

func (svc *DatastoreService) Update(ctx context.Context, c ccx.Datastore) (*ccx.Datastore, error) {
	old, err := svc.Read(ctx, c.ID)
	if errors.Is(err, ccx.ResourceNotFoundErr) {
		return nil, ccx.ResourceNotFoundErr
	} else if err != nil {
		return nil, err
	}

	if hasCan, err := hasSupportedChanges(*old, c); err != nil {
		return nil, err
	} else if !hasCan {
		return old, nil
	}

	var ur UpdateRequest

	if old.Name != c.Name {
		ur.NewName = c.Name
	}

	if old.VolumeSize != c.VolumeSize {
		ur.NewVolumeSize = uint(c.VolumeSize)
	}

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(ur); err != nil {
		return nil, errors.Join(ccx.RequestEncodingErr, err)
	}

	url := svc.baseURL + "/api/prov/api/v2/cluster/" + c.ID
	req, err := http.NewRequest(http.MethodPatch, url, &b)
	if err != nil {
		return nil, errors.Join(ccx.RequestInitializationErr, err)
	}

	token, err := svc.auth.Auth(ctx)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", token)
	client := &http.Client{Timeout: ccx.DefaultTimeout}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Join(ccx.RequestSendingErr, err)
	}

	if res.StatusCode == http.StatusBadRequest {
		return nil, lib.ErrorFromErrorResponse(res.Body)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status = %d", ccx.ResponseStatusFailedErr, res.StatusCode)
	}

	var rs DatastoreResponse
	if err := lib.DecodeJsonInto(res.Body, &rs); err != nil {
		return nil, err
	}

	updatedStore := DatastoreFromResponse(rs)

	if err := svc.LoadAll(ctx); err != nil {
		return nil, errors.Join(ccx.ResourcesLoadFailedErr, err)
	}

	return &updatedStore, nil
}

func hasSupportedChanges(old, c ccx.Datastore) (bool, error) {
	var (
		hasCan, hasCant bool
		fields          []string
	)

	if old.Name != c.Name {
		hasCan = true
	}

	if old.Size != c.Size {
		hasCant = true
		fields = append(fields, "cluster_size")
	}

	if old.VolumeSize != c.VolumeSize {
		// hasCan = true
		hasCant = true
		fields = append(fields, "volume_size")
	}

	if old.DBVendor != c.DBVendor {
		hasCant = true
		fields = append(fields, "db_vendor")
	}

	if old.DBVersion != c.DBVersion {
		hasCant = true
		fields = append(fields, "db_version")
	}

	if old.Type != c.Type {
		hasCant = true
		fields = append(fields, "cluster_type")
	}
	if old.CloudProvider != c.CloudProvider {
		hasCant = true
		fields = append(fields, "cloud_provider")
	}

	if old.CloudRegion != c.CloudRegion {
		hasCant = true
		fields = append(fields, "cloud_region")
	}

	if old.InstanceSize != c.InstanceSize {
		hasCant = true
		fields = append(fields, "instance_size")
	}

	if old.VolumeType != c.VolumeType {
		hasCant = true
		fields = append(fields, "volume_type")
	}

	if old.VolumeIOPS != c.VolumeIOPS {
		hasCant = true
		fields = append(fields, "volume_iops")
	}

	if old.HAEnabled != c.HAEnabled {
		hasCant = true
		fields = append(fields, "ha_enabled")
	}

	if old.VpcUUID != c.VpcUUID {
		hasCant = true
		fields = append(fields, "vpc_uuid")
	}

	if !slices.Equal(old.AvailabilityZones, c.AvailabilityZones) {
		hasCant = true
		fields = append(fields, "availability_zones")
	}

	if hasCant {
		return hasCan, fmt.Errorf("%w: %s", ccx.UpdateNotSupportedErr, strings.Join(fields, ", "))
	}

	return hasCan, nil
}
