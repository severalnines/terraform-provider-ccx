package api

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type jobType string

const (
	deployStoreJob    jobType = "JOB_TYPE_DEPLOY_DATASTORE"
	modifyDbConfigJob jobType = "JOB_TYPE_MODIFYDBCONFIG"
	destroyStoreJob   jobType = "JOB_TYPE_DESTROY_DATASTORE"
	addNodeJob        jobType = "JOB_TYPE_ADD_NODE"
	removeNodeJob     jobType = "JOB_TYPE_REMOVE_NODE"
)

type jobStatus string

const (
	jobStatusUnknown  jobStatus = "JOB_STATUS_UNKNOWN"
	jobStatusRunning  jobStatus = "JOB_STATUS_RUNNING"
	jobStatusFinished jobStatus = "JOB_STATUS_FINISHED"
	jobStatusErrored  jobStatus = "JOB_STATUS_ERRORED"
)

func newJobs(httpcli HttpClient, timeout time.Duration) jobs {
	return jobs{
		timeout: timeout,
		httpcli: httpcli,
		tick:    time.Second * 30,
	}
}

type jobs struct {
	httpcli HttpClient
	timeout time.Duration
	tick    time.Duration // time to wait between job status checks
}

type jobsResponse struct {
	Jobs  []jobsResponseJobItem `json:"jobs"`
	Total int
}

type jobsResponseJobItem struct {
	JobID  string    `json:"job_id"`
	Type   jobType   `json:"type"`
	Status jobStatus `json:"status"`
}

type jobService interface {
	Await(ctx context.Context, storeID string, job jobType) (jobStatus, error)
	GetStatus(_ context.Context, storeID string, job jobType) (jobStatus, error)
}

func (svc jobs) Await(ctx context.Context, storeID string, job jobType) (jobStatus, error) {
	timeout := time.Now().Add(svc.timeout)
	ticker := time.NewTicker(svc.tick)

	var (
		status jobStatus
		err    error
	)

	for time.Now().Before(timeout) {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				return jobStatusUnknown, fmt.Errorf("context cancelled with error %w", err)
			}
			return jobStatusUnknown, errors.New("context cancelled")
		default:
			break
		}

		status, err = svc.GetStatus(ctx, storeID, job)

		if err != nil {
			continue
		}

		switch status {
		case jobStatusFinished, jobStatusErrored:
			return status, nil
		}

		<-ticker.C
	}

	if err != nil {
		return jobStatusUnknown, err
	}

	return jobStatusUnknown, fmt.Errorf("job did not finish in %s", svc.timeout)
}

func (svc jobs) GetStatus(ctx context.Context, storeID string, job jobType) (jobStatus, error) {
	var rs jobsResponse
	if err := svc.httpcli.Get(ctx, "/api/deployment/v2/data-stores/"+storeID+"/jobs?limit=10&offset=0", &rs); err != nil {
		return jobStatusUnknown, fmt.Errorf("getting job status: %w", err)
	}

	for i := range rs.Jobs {
		if rs.Jobs[i].Type == job {
			return rs.Jobs[i].Status, nil
		}
	}

	return jobStatusUnknown, nil
}
