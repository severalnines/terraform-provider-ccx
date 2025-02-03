package api

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
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
	InstanceSize string `json:"instance_size"`
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
		return nil, fmt.Errorf("resizing datastore: %w", err)
	}

	if old.ParameterGroupID != next.ParameterGroupID {
		if err := svc.ApplyParameterGroup(ctx, next.ID, next.ParameterGroupID); err != nil {
			return nil, fmt.Errorf("applying parameter group: %w", err)
		}

		updated = true
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

	_, err := svc.client.Do(ctx, http.MethodPatch, "/api/prov/api/v2/cluster/"+next.ID, ur)
	if err != nil {
		return false, fmt.Errorf("updating datastore: %w", err)
	}

	return true, nil
}

func (svc *DatastoreService) resize(ctx context.Context, old, next ccx.Datastore) (bool, error) {
	if old.Size == next.Size {
		return false, nil
	}

	adding := old.Size < next.Size

	if adding {
		have := len(next.AvailabilityZones)
		need := int(next.Size - old.Size)
		missing := need - have

		if next.VpcUUID == "" && missing > 0 { // allocate AZs if public and need is less than have
			allAzs, err := svc.contentSvc.AvailabilityZones(ctx, next.CloudProvider, next.CloudRegion)
			if err != nil {
				return false, fmt.Errorf("%w: %w", ccx.AllocatingAZsErr, err)
			}

			existing := make([]string, 0, len(old.Hosts))

			for i := range old.Hosts {
				existing = append(existing, old.Hosts[i].AZ)
			}

			next.AvailabilityZones = append(next.AvailabilityZones, allocateAzs(allAzs, existing, missing)...)
		}
	}

	ur, err := svc.updateSizeRequest(ctx, old, next)
	if err != nil {
		return false, fmt.Errorf("computing resize: %w", err)
	}

	_, err = svc.client.Do(ctx, http.MethodPatch, "/api/prov/api/v2/cluster/"+next.ID, ur)
	if err != nil {
		return false, err
	}

	var jt ccx.JobType
	if adding {
		jt = ccx.AddNodeJob
	} else {
		jt = ccx.RemoveNodeJob
	}

	status, err := svc.jobs.Await(ctx, old.ID, jt)
	if err != nil {
		return false, fmt.Errorf("awaiting resize job: %w", err)
	} else if status != ccx.JobStatusFinished {
		return false, fmt.Errorf("resize job failed: %s", status)
	}

	return true, nil
}

func (svc *DatastoreService) updateSizeRequest(ctx context.Context, old, next ccx.Datastore) (updateRequest, error) {
	var ur updateRequest

	if old.Size > next.Size { // remove the oldest non-primary nodes
		ids, err := oldestRemovableNodeIds(old.Hosts, int(old.Size-next.Size))
		if err != nil {
			return ur, err
		}

		ur.Remove = &removeHosts{HostIDs: ids}
		return ur, nil
	}

	have := len(next.AvailabilityZones)
	need := int(next.Size - old.Size)
	missing := need - have

	if old.VpcUUID == "" && missing > 0 { // allocate AZs if public and need is less than have
		allAzs, err := svc.contentSvc.AvailabilityZones(ctx, next.CloudProvider, next.CloudRegion)
		if err != nil {
			return ur, fmt.Errorf("%w: %w", ccx.AllocatingAZsErr, err)
		}

		existing := make([]string, 0, len(old.Hosts))

		for i := range old.Hosts {
			existing = append(existing, old.Hosts[i].AZ)
		}

		ls := allocateAzs(allAzs, existing, missing)
		next.AvailabilityZones = append(next.AvailabilityZones, ls...)
	}

	// add new hosts based on newest node spec
	specs, err := newestNodeSpecs(old.Hosts, int(next.Size-old.Size), next.AvailabilityZones)
	if err != nil {
		return ur, err
	}

	ur.Add = &addHosts{Specs: specs}

	return ur, nil
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

func newestNodeSpecs(hosts []ccx.Host, count int, azs []string) ([]hostSpecs, error) {
	if len(hosts) == 0 {
		return nil, fmt.Errorf("no nodes available")
	}

	if len(azs) != count {
		return nil, fmt.Errorf("not enough azs available")
	}

	slices.SortStableFunc(hosts, func(a, b ccx.Host) int {
		return b.CreatedAt.Compare(a.CreatedAt)
	})

	h := hosts[0]
	ls := make([]hostSpecs, 0, count)

	for i := 0; i < count; i++ {
		ls = append(ls, hostSpecs{
			InstanceSize: h.InstanceType,
			AZ:           azs[i],
		})
	}

	return ls, nil
}
