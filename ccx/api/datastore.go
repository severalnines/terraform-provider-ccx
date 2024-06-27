package api

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/severalnines/terraform-provider-ccx/ccx"
)

type DatastoreService struct {
	baseURL string
	auth    authorizer
	jobs    jobs

	stores  map[string]ccx.Datastore
	mut     sync.Mutex
	timeout time.Duration
}

// Datastores creates a new datastores DatastoreService
func Datastores(ctx context.Context, baseURL, clientID, clientSecret string, timeout time.Duration) (*DatastoreService, error) {
	a := tokenAuthorizer{
		id:      clientID,
		secret:  clientSecret,
		baseURL: baseURL,
	}

	j := jobs{
		baseURL: baseURL,
		auth:    a,
	}

	c := DatastoreService{
		baseURL: baseURL,
		auth:    a,
		jobs:    j,
		timeout: timeout,
	}

	err := c.LoadAll(ctx)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// String returns a string representation of the internal stores map
// useful for debugging
func (svc *DatastoreService) String() string {
	svc.mut.Lock()
	defer svc.mut.Unlock()

	if len(svc.stores) == 0 {
		return "<empty>"
	}

	var b bytes.Buffer
	for id, store := range svc.stores {
		b.WriteString(fmt.Sprintf("id = %s, name = %s\n", id, store.Name))
	}

	return b.String()
}
