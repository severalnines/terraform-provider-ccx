package datastore_client

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/severalnines/terraform-provider-ccx/ccx"
	chttp "github.com/severalnines/terraform-provider-ccx/http"
)

var _ ccx.DatastoreService = &Client{}

type Client struct {
	auth   chttp.Authorizer
	conn   *chttp.ConnectionParameters
	stores map[string]ccx.Datastore
	mut    sync.Mutex
}

// New creates a new datastores Client
func New(ctx context.Context, authorizer chttp.Authorizer, opts ...chttp.ParameterOption) (*Client, error) {
	p := chttp.Parameters(opts...)

	c := Client{
		auth: authorizer,
		conn: p,
	}

	err := c.LoadAll(ctx)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// String returns a string representation of the internal stores map
// useful for debugging
func (cli *Client) String() string {
	cli.mut.Lock()
	defer cli.mut.Unlock()

	if len(cli.stores) == 0 {
		return "<empty>"
	}

	var b bytes.Buffer
	for id, store := range cli.stores {
		b.WriteString(fmt.Sprintf("id = %s, name = %s\n", id, store.Name))
	}

	return b.String()
}
