package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/severalnines/terraform-provider-ccx/ccx"
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
	url := svc.baseURL + "/api/firewall/api/v1/firewalls/" + storeID

	var rs getFirewallsResponse

	if err := httpGet(ctx, svc.auth, url, &rs); err != nil {
		return nil, err
	}

	ls := make([]ccx.FirewallRule, 0, len(rs))
	for _, r := range rs {
		ls = append(ls, ccx.FirewallRule{
			Source:      r.Source,
			Description: r.Description,
		})
	}

	return ls, nil
}

func firewallsDiff(have, want []ccx.FirewallRule) (create, del []ccx.FirewallRule) {
	haveM := make(map[string]ccx.FirewallRule, len(have))
	for _, f := range have {
		haveM[f.Source] = f
	}

	for _, f := range want {
		if _, ok := haveM[f.Source]; !ok {
			create = append(create, f)
		}
	}

	wantM := make(map[string]ccx.FirewallRule, len(want))
	for _, f := range want {
		wantM[f.Source] = f
	}

	for _, f := range have {
		if _, ok := wantM[f.Source]; !ok {
			del = append(del, f)
		}
	}

	return create, del
}

func (svc *DatastoreService) CreateFirewallRule(ctx context.Context, storeID string, firewall ccx.FirewallRule) error {
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(firewall); err != nil {
		return errors.Join(ccx.RequestEncodingErr, err)
	}

	url := svc.baseURL + "/api/firewall/api/v1/firewall/" + storeID
	req, err := http.NewRequest(http.MethodPost, url, &b)
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

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("%w: status = %d", lib.ErrorFromErrorResponse(res.Body), res.StatusCode)
	}

	return nil
}

func (svc *DatastoreService) CreateFirewallRules(ctx context.Context, storeID string, firewalls []ccx.FirewallRule) error {
	limiter := make(chan bool, 10)

	var eg errgroup.Group

	for _, f := range firewalls {
		limiter <- true
		f := f
		eg.Go(func() error {
			defer func() { <-limiter }()
			err := svc.CreateFirewallRule(ctx, storeID, f)

			if err != nil {
				return fmt.Errorf("creating rule (source=%s, description=%s): %w", f.Source, f.Description, err)
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func (svc *DatastoreService) DeleteFirewallRule(ctx context.Context, storeID string, firewall ccx.FirewallRule) error {
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(firewall); err != nil {
		return errors.Join(ccx.RequestEncodingErr, err)
	}

	url := svc.baseURL + "/api/firewall/api/v1/firewall/" + storeID
	req, err := http.NewRequest(http.MethodDelete, url, &b)
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
	have, err := svc.GetFirewallRules(ctx, storeID)
	if err != nil {
		return fmt.Errorf("getting firewalls: %w", err)
	}

	create, del := firewallsDiff(have, firewalls)

	if len(create) > 0 {
		err = svc.CreateFirewallRules(ctx, storeID, create)
	}

	if err != nil {
		return fmt.Errorf("creating rules: %w", err)
	}

	if len(del) > 0 {
		err = svc.DeleteFirewallRules(ctx, storeID, del)
	}

	if err != nil {
		return fmt.Errorf("deleting rules: %w", err)
	}

	return nil
}
