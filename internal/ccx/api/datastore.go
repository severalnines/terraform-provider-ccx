package api

import (
	"context"
	"time"

	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
)

type DatastoreService struct {
	httpcli HttpClient
	jobs    jobService
}

var _ ccx.DatastoreService = (*DatastoreService)(nil)

// Datastores creates a new datastores DatastoreService
func Datastores(_ context.Context, httpcli HttpClient, timeout time.Duration) (*DatastoreService, error) {
	j := newJobs(httpcli, timeout)

	c := DatastoreService{
		httpcli: httpcli,
		jobs:    j,
	}

	return &c, nil
}
