package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
	"github.com/severalnines/terraform-provider-ccx/internal/lib"
	"golang.org/x/sync/errgroup"
)

type getFirewallsResponse []struct {
	// Source is CIDR of remote
	Source string `json:"source"`
	// Description is human readable description of the source
	Description string `json:"description"`
	// Ports is all the ports available to the source
	Ports []struct {
		// Port is "service" or something
		Port string `json:"port"`
		// PortNo is informational - it is the port being used the chosen purpose
		PortNo int `json:"port_no"`
	} `json:"ports"`
}

func (svc *DatastoreService) GetFirewallRules(ctx context.Context, storeID string) ([]ccx.FirewallRule, error) {
	var rs getFirewallsResponse

	if err := svc.client.Get(ctx, "/api/firewall/api/v1/firewalls/"+storeID, &rs); err != nil {
		return nil, err
	}

	ls := make([]ccx.FirewallRule, 0, len(rs))
	for _, r := range rs {
		ls = append(ls, ccx.FirewallRule{
			Source:      r.Source,
			Description: r.Description,
		})
	}

	slices.SortStableFunc(ls, func(a, b ccx.FirewallRule) int {
		return strings.Compare(a.Source, b.Source)
	})

	return ls, nil
}

func firewallsDiff(have, want []ccx.FirewallRule) (create, del []ccx.FirewallRule) {
	haveM := make(map[ccx.FirewallRule]struct{}, len(have))
	for _, f := range have {
		haveM[f] = struct{}{}
	}

	for _, f := range want {
		if _, ok := haveM[f]; !ok {
			create = append(create, f)
		}
	}

	wantM := make(map[ccx.FirewallRule]struct{}, len(want))
	for _, f := range want {
		wantM[f] = struct{}{}
	}

	for _, f := range have {
		if _, ok := wantM[f]; !ok {
			del = append(del, f)
		}
	}

	return create, del
}

func (svc *DatastoreService) CreateFirewallRule(ctx context.Context, storeID string, firewall ccx.FirewallRule) error {
	res, err := svc.client.Do(ctx, http.MethodPost, "/api/firewall/api/v1/firewall/"+storeID, firewall)
	if err != nil {
		return errors.Join(ccx.RequestSendingErr, err)
	}

	if res.StatusCode == http.StatusBadRequest {
		return lib.ErrorFromErrorResponse(res.Body)
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("%w: status = %d", lib.ErrorFromErrorResponse(res.Body), res.StatusCode)
	}

	return nil
}

func (svc *DatastoreService) CreateFirewallRules(ctx context.Context, storeID string, firewalls []ccx.FirewallRule) error {
	var errs []error

	for _, f := range firewalls {
		err := svc.CreateFirewallRule(ctx, storeID, f)
		if err != nil {
			errs = append(errs, fmt.Errorf("creating rule (source=%s, description=%s): %w", f.Source, f.Description, err))
			break
		}
	}

	if len(errs) != 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (svc *DatastoreService) DeleteFirewallRule(ctx context.Context, storeID string, firewall ccx.FirewallRule) error {
	res, err := svc.client.Do(ctx, http.MethodDelete, "/api/firewall/api/v1/firewall/"+storeID, firewall)
	if err != nil {
		return errors.Join(ccx.RequestSendingErr, err)
	}

	if res.StatusCode == http.StatusBadRequest {
		return lib.ErrorFromErrorResponse(res.Body)
	}

	if res.StatusCode != http.StatusAccepted {
		return fmt.Errorf("%w: status = %d", lib.ErrorFromErrorResponse(res.Body), res.StatusCode)
	}

	return nil
}

func (svc *DatastoreService) DeleteFirewallRules(ctx context.Context, storeID string, firewalls []ccx.FirewallRule) error {
	limiter := make(chan bool, 10)

	var eg errgroup.Group

	for _, f := range firewalls {
		limiter <- true
		f := f
		eg.Go(func() error {
			defer func() { <-limiter }()
			err := svc.DeleteFirewallRule(ctx, storeID, f)

			if err != nil {
				return fmt.Errorf("deleting rule (source=%s, description=%s): %w", f.Source, f.Description, err)
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func (svc *DatastoreService) SetFirewallRules(ctx context.Context, storeID string, firewalls []ccx.FirewallRule) error {
	slices.SortStableFunc(firewalls, func(a, b ccx.FirewallRule) int {
		return strings.Compare(a.Source, b.Source)
	})

	have, err := svc.GetFirewallRules(ctx, storeID)
	if err != nil {
		return fmt.Errorf("getting firewalls: %w", err)
	}

	create, del := firewallsDiff(have, firewalls)

	if len(del) > 0 {
		if err = svc.DeleteFirewallRules(ctx, storeID, del); err != nil {
			return fmt.Errorf("deleting rules: %w", err)
		}
	}

	if len(create) > 0 {
		if err = svc.CreateFirewallRules(ctx, storeID, create); err != nil {
			return fmt.Errorf("creating rules: %w", err)
		}
	}

	return nil
}
