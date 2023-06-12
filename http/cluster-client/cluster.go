package cluster_client

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	ccxprov "github.com/severalnines/terraform-provider-ccx"
	chttp "github.com/severalnines/terraform-provider-ccx/http"
)

var _ ccxprov.ClusterService = &Client{}

type Client struct {
	auth     chttp.Authorizer
	conn     *chttp.ConnectionParameters
	clusters map[string]ccxprov.Cluster
	mut      sync.Mutex
}

// New creates a new clusters Client
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

// String returns a string representation of the internal clusters map
// useful for debugging
func (cli *Client) String() string {
	cli.mut.Lock()
	defer cli.mut.Unlock()

	if len(cli.clusters) == 0 {
		return "<empty>"
	}

	var b bytes.Buffer
	for id, cluster := range cli.clusters {
		b.WriteString(fmt.Sprintf("id = %s, name = %s\n", id, cluster.ClusterName))
	}

	return b.String()
}
