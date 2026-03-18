package ccx

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func NewJobsClient(httpcli HTTPClient, timeout time.Duration) JobsService {
	return &JobsClient{
		timeout: timeout,
		httpcli: httpcli,
		tick:    time.Second * 30,
	}
}

type JobsClient struct {
	httpcli HTTPClient
	timeout time.Duration
	tick    time.Duration // time to wait between job status checks
}

type jobsResponse struct {
	Jobs  []jobsResponseJobItem `json:"jobs"`
	Total int
}

type jobsResponseJobItem struct {
	JobID  string    `json:"job_id"`
	Type   JobType   `json:"type"`
	Status JobStatus `json:"status"`
}

func (svc *JobsClient) Await(ctx context.Context, storeID string, job JobType) (JobStatus, error) {
	timeout := time.Now().Add(svc.timeout)
	ticker := time.NewTicker(svc.tick)

	var (
		status JobStatus
		err    error
	)

	for time.Now().Before(timeout) {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				return JobStatusUnknown, fmt.Errorf("context cancelled with error %w", err)
			}
			return JobStatusUnknown, errors.New("context cancelled")
		default:
		}

		status, err = svc.GetStatus(ctx, storeID, job)

		if err != nil {
			return JobStatusUnknown, fmt.Errorf("getting job status: %w", err)
		}

		switch status {
		case JobStatusFinished, JobStatusErrored:
			return status, nil
		}

		<-ticker.C
	}

	if err != nil {
		return JobStatusUnknown, err
	}

	return JobStatusUnknown, fmt.Errorf("job did not finish in %s", svc.timeout)
}

func (svc *JobsClient) GetStatus(ctx context.Context, storeID string, job JobType) (JobStatus, error) {
	var rs jobsResponse
	if err := svc.httpcli.Get(ctx, "/api/deployment/v2/data-stores/"+storeID+"/jobs?limit=10&offset=0", &rs); err != nil {
		return JobStatusUnknown, fmt.Errorf("getting job status: %w", err)
	}

	for i := range rs.Jobs {
		if rs.Jobs[i].Type == job {
			return rs.Jobs[i].Status, nil
		}
	}

	return JobStatusUnknown, nil
}
