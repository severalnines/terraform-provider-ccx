package api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
)

func Jobs(httpcli HttpClient, timeout time.Duration) JobsService {
	return JobsService{
		timeout: timeout,
		httpcli: httpcli,
		tick:    time.Second * 30,
	}
}

type JobsService struct {
	httpcli HttpClient
	timeout time.Duration
	tick    time.Duration // time to wait between job status checks
}

type jobsResponse struct {
	Jobs  []jobsResponseJobItem `json:"jobs"`
	Total int
}

type jobsResponseJobItem struct {
	JobID  string        `json:"job_id"`
	Type   ccx.JobType   `json:"type"`
	Status ccx.JobStatus `json:"status"`
}

func (svc JobsService) Await(ctx context.Context, storeID string, job ccx.JobType) (ccx.JobStatus, error) {
	timeout := time.Now().Add(svc.timeout)
	ticker := time.NewTicker(svc.tick)

	var (
		status ccx.JobStatus
		err    error
	)

	for time.Now().Before(timeout) {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				return ccx.JobStatusUnknown, fmt.Errorf("context cancelled with error %w", err)
			}
			return ccx.JobStatusUnknown, errors.New("context cancelled")
		default:
		}

		status, err = svc.GetStatus(ctx, storeID, job)

		if err != nil {
			return ccx.JobStatusUnknown, fmt.Errorf("getting job status: %w", err)
		}

		switch status {
		case ccx.JobStatusFinished, ccx.JobStatusErrored:
			return status, nil
		}

		<-ticker.C
	}

	if err != nil {
		return ccx.JobStatusUnknown, err
	}

	return ccx.JobStatusUnknown, fmt.Errorf("job did not finish in %s", svc.timeout)
}

func (svc JobsService) GetStatus(ctx context.Context, storeID string, job ccx.JobType) (ccx.JobStatus, error) {
	var rs jobsResponse
	if err := svc.httpcli.Get(ctx, "/api/deployment/v2/data-stores/"+storeID+"/jobs?limit=10&offset=0", &rs); err != nil {
		return ccx.JobStatusUnknown, fmt.Errorf("getting job status: %w", err)
	}

	for i := range rs.Jobs {
		if rs.Jobs[i].Type == job {
			return rs.Jobs[i].Status, nil
		}
	}

	return ccx.JobStatusUnknown, nil
}
