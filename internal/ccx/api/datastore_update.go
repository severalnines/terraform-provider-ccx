package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
	"github.com/severalnines/terraform-provider-ccx/internal/lib"
)

type updateRequest struct {
	NewName       string         `json:"cluster_name"`
	NewVolumeSize uint           `json:"new_volume_size"`
	Remove        *removeHosts   `json:"remove_nodes"`
	Add           *addHosts      `json:"add_nodes"`
	Notifications *notifications `json:"notifications"`
	Maintenance   *maintenance   `json:"maintenance_settings"`
}

type removeHosts struct {
	HostIDs []string `json:"node_uuids"`
}

type addHosts struct {
	Specs []hostSpecs `json:"specs"`
}

type notifications struct {
	Enabled bool     `json:"enabled"`
	Emails  []string `json:"emails"`
}

type maintenance struct {
	DayOfWeek uint32 `json:"day_of_week"`
	StartHour uint64 `json:"start_hour"`
	EndHour   uint64 `json:"end_hour"`
}

type hostSpecs struct {
	InstanceType string `json:"instance_size"`
	AZ           string `json:"availability_zone"`
}

func (svc *DatastoreService) Update(ctx context.Context, old, next ccx.Datastore) (*ccx.Datastore, error) {
	out := &old

	updated, err := svc.update(ctx, old, next)
	if err != nil {
		return nil, err
	}

	resized, err := svc.resize(ctx, old, next)
	if err != nil {
		return nil, err
	}

	if updated || resized {
		if out, err = svc.Read(ctx, old.ID); err != nil {
			return nil, err
		}
	}

	return out, nil
}

func (svc *DatastoreService) update(ctx context.Context, old, next ccx.Datastore) (bool, error) {
	ur, ok := svc.updateRequest(old, next)
	if !ok {
		return false, nil
	}

	res, err := svc.client.Do(ctx, http.MethodPatch, "/api/prov/api/v2/cluster/"+next.ID, ur)
	if err != nil {
		return false, errors.Join(ccx.RequestSendingErr, err)
	}

	if res.StatusCode == http.StatusBadRequest {
		return false, lib.ErrorFromErrorResponse(res.Body)
	}

	if res.StatusCode != http.StatusOK {
		return false, fmt.Errorf("%w: status = %d", ccx.ResponseStatusFailedErr, res.StatusCode)
	}

	return true, nil
}

func (svc *DatastoreService) resize(ctx context.Context, old, next ccx.Datastore) (bool, error) {
	ur, ok := svc.updateSizeRequest(old, next)
	if !ok {
		return false, nil
	}

	res, err := svc.client.Do(ctx, http.MethodPatch, "/api/prov/api/v2/cluster/"+next.ID, ur)
	if err != nil {
		return false, errors.Join(ccx.RequestSendingErr, err)
	}

	if res.StatusCode == http.StatusBadRequest {
		return false, lib.ErrorFromErrorResponse(res.Body)
	}

	if res.StatusCode != http.StatusOK {
		return false, fmt.Errorf("%w: status = %d", ccx.ResponseStatusFailedErr, res.StatusCode)
	}

	var jt jobType

	if old.Size > next.Size {
		jt = removeNodeJob
	} else if old.Size < next.Size {
		jt = addNodeJob
	} else {
		// should not be here
		return false, nil
	}

	status, err := svc.jobs.Await(ctx, old.ID, jt)
	if err != nil {
		return false, fmt.Errorf("%w: awaiting resize job: %w", ccx.CreateFailedErr, err)
	} else if status != jobStatusFinished {
		return false, fmt.Errorf("%w: resize job failed: %s", ccx.CreateFailedErr, status)
	}

	return true, nil
}

func (svc *DatastoreService) updateSizeRequest(old, next ccx.Datastore) (updateRequest, bool) {
	var (
		ur updateRequest
		ok bool
	)

	if old.Size > next.Size { // remove the oldest non-primary nodes
		ids, err := oldestRemovableNodeIds(old.Hosts, int(old.Size-next.Size))
		if err != nil {
			return ur, false
		}

		ur.Remove = &removeHosts{HostIDs: ids}
		ok = true
	} else if old.Size < next.Size { // add new hosts based on newest node spec
		specs, err := newestNodeSpecs(old.Hosts, int(next.Size-old.Size))
		if err != nil {
			return ur, false
		}

		ur.Add = &addHosts{Specs: specs}
		ok = true
	}

	return ur, ok
}

func (svc *DatastoreService) updateRequest(old, next ccx.Datastore) (updateRequest, bool) {
	var (
		ur updateRequest
		ok bool
	)

	if old.Name != next.Name {
		ur.NewName = next.Name
		ok = true
	}

	if old.VolumeSize != next.VolumeSize {
		ur.NewVolumeSize = uint(next.VolumeSize)
		ok = true
	}

	if old.Notifications.Enabled != next.Notifications.Enabled || !slices.Equal(old.Notifications.Emails, next.Notifications.Emails) {
		ur.Notifications = &notifications{
			Enabled: next.Notifications.Enabled,
			Emails:  next.Notifications.Emails,
		}

		ok = true
	}

	if next.MaintenanceSettings != nil {
		ur.Maintenance = &maintenance{
			DayOfWeek: uint32(next.MaintenanceSettings.DayOfWeek),
			StartHour: uint64(next.MaintenanceSettings.StartHour),
			EndHour:   uint64(next.MaintenanceSettings.EndHour),
		}

		ok = true
	}

	return ur, ok
}

func oldestRemovableNodeIds(hosts []ccx.Host, count int) ([]string, error) {
	hosts = slices.DeleteFunc(hosts, func(h ccx.Host) bool {
		return h.IsPrimary()
	})

	if len(hosts) < count {
		return nil, fmt.Errorf("cannot remove %d nodes, only %d non-primary available", count, len(hosts))
	}

	slices.SortStableFunc(hosts, func(a, b ccx.Host) int {
		return a.CreatedAt.Compare(b.CreatedAt)
	})

	ls := make([]string, 0, count)
	for i := 0; i < count; i++ {
		ls = append(ls, hosts[i].ID)
	}

	return ls, nil
}

func newestNodeSpecs(hosts []ccx.Host, count int) ([]hostSpecs, error) {
	if len(hosts) == 0 {
		return nil, fmt.Errorf("no nodes available")
	}

	slices.SortStableFunc(hosts, func(a, b ccx.Host) int {
		return b.CreatedAt.Compare(a.CreatedAt)
	})

	h := hosts[0]
	ls := make([]hostSpecs, 0, count)

	for i := 0; i < count; i++ {
		ls = append(ls, hostSpecs{
			InstanceType: h.InstanceType,
			AZ:           h.AZ,
		})
	}

	return ls, nil
}
