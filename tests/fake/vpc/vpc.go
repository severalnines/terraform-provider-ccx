package vpc

import (
	"context"
	"sync"

	"github.com/google/uuid"
	ccxprov "github.com/severalnines/terraform-provider-ccx"
)

var _ ccxprov.VPCService = &Client{}

var (
	instance     *Client
	instanceLock sync.Mutex
)

func Instance(m map[string]ccxprov.VPC) *Client {
	instanceLock.Lock()
	defer instanceLock.Unlock()
	if instance == nil {
		instance = New(m)
	}
	return instance
}

func New(m map[string]ccxprov.VPC) *Client {
	var c Client

	if m == nil {
		c.VPCs = make(map[string]ccxprov.VPC)
	} else {
		c.VPCs = m
	}

	return &c
}

type Client struct {
	VPCs map[string]ccxprov.VPC
}

func (m Client) Create(ctx context.Context, v ccxprov.VPC) (*ccxprov.VPC, error) {
	v.ID = uuid.NewString()
	m.VPCs[v.ID] = v
	return &v, nil
}

func (m Client) Read(ctx context.Context, id string) (*ccxprov.VPC, error) {
	c, ok := m.VPCs[id]
	if !ok {
		return nil, ccxprov.ResourceNotFoundErr
	}
	return &c, nil
}

func (m Client) Update(ctx context.Context, c ccxprov.VPC) (*ccxprov.VPC, error) {
	m.VPCs[c.ID] = c
	return &c, nil
}

func (m Client) Delete(ctx context.Context, id string) error {
	delete(m.VPCs, id)
	return nil
}
