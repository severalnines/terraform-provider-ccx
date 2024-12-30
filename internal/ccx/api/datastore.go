package api

import (
	"time"

	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
)

type DatastoreService struct {
	client     HttpClient
	jobs       jobService
	contentSvc ccx.ContentService
}

var _ ccx.DatastoreService = (*DatastoreService)(nil)

// Datastores creates a new datastores DatastoreService
func Datastores(client HttpClient, timeout time.Duration, contentSvc ccx.ContentService) (*DatastoreService, error) {
	j := newJobs(client, timeout)

	c := DatastoreService{
		client:     client,
		jobs:       j,
		contentSvc: contentSvc,
	}

	return &c, nil
}
