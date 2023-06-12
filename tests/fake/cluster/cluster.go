package cluster

import (
	"bytes"
	"context"
	"encoding/json"
	"sync"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	ccxprov "github.com/severalnines/terraform-provider-ccx"
)

var _ ccxprov.ClusterService = &Client{}

var (
	instance     *Client
	instanceLock sync.Mutex
)

func Instance(m map[string]ccxprov.Cluster) *Client {
	instanceLock.Lock()
	defer instanceLock.Unlock()
	if instance == nil {
		instance = New(m)
	}
	return instance
}

func New(m map[string]ccxprov.Cluster) *Client {
	var c Client

	if m == nil {
		c.Clusters = make(map[string]ccxprov.Cluster)
	} else {
		c.Clusters = m
	}

	return &c
}

func debugClusters(m map[string]ccxprov.Cluster) []byte {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)

	for k := range m {
		_ = enc.Encode(m[k])
		b.WriteString("\n")
	}

	return b.Bytes()
}

type Client struct {
	Clusters map[string]ccxprov.Cluster
}

func (m Client) Create(ctx context.Context, c ccxprov.Cluster) (*ccxprov.Cluster, error) {
	c.ID = uuid.NewString()
	m.Clusters[c.ID] = c
	return &c, nil
}

func (m Client) Read(ctx context.Context, id string) (*ccxprov.Cluster, error) {
	c, ok := m.Clusters[id]
	tflog.Info(ctx, "reading")
	if !ok {
		return nil, ccxprov.ResourceNotFoundErr
	}
	return &c, nil
}

func (m Client) Update(ctx context.Context, c ccxprov.Cluster) (*ccxprov.Cluster, error) {
	m.Clusters[c.ID] = c
	return &c, nil
}

func (m Client) Delete(ctx context.Context, id string) error {
	delete(m.Clusters, id)
	return nil
}
