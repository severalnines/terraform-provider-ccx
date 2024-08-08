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
)

type jobStatus string

const (
	jobStatusUnknown  jobStatus = "JOB_STATUS_UNKNOWN"
	jobStatusRunning  jobStatus = "JOB_STATUS_RUNNING"
	jobStatusFinished jobStatus = "JOB_STATUS_FINISHED"
	jobStatusErrored  jobStatus = "JOB_STATUS_ERRORED"
)

type jobs struct {
	baseURL             string
	auth                authorizer
	awaitTickerDuration time.Duration
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

func (svc jobs) Await(ctx context.Context, storeID string, job jobType, wait time.Duration) (jobStatus, error) {
	timeout := time.Now().Add(wait)
	ticker := time.NewTicker(svc.awaitTickerDuration)

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

	return jobStatusUnknown, fmt.Errorf("job did not finish in %s", wait)
}

func (svc jobs) GetStatus(ctx context.Context, storeID string, job jobType) (jobStatus, error) {
	url := svc.baseURL + "/api/deployment/v2/data-stores/" + storeID + "/jobs?limit=10&offset=0"

	var rs jobsResponse
	if err := httpGet(ctx, svc.auth, url, &rs); err != nil {
		return jobStatusUnknown, fmt.Errorf("getting job status: %w", err)
	}

	for i := range rs.Jobs {
		if rs.Jobs[i].Type == job {
			return rs.Jobs[i].Status, nil
		}
	}

	return jobStatusUnknown, nil
}
