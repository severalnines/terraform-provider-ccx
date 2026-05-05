package ccx

import (
	"time"
)

type DatastoresClient struct {
	client     HTTPClient
	jobs       JobsService
	contentSvc ContentService
}

var _ DatastoresService = (*DatastoresClient)(nil)

// NewDatastoresClient creates a new datastores DatastoreService
func NewDatastoresClient(client HTTPClient, timeout time.Duration, contentSvc ContentService) (DatastoresService, error) {
	j := NewJobsClient(client, timeout)

	c := DatastoresClient{
		client:     client,
		jobs:       j,
		contentSvc: contentSvc,
	}

	return &c, nil
}
