package resources

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx/api"
)

var (
	_ TerraformResource = &VPC{}
)

func schemaToVPC(d *schema.ResourceData) ccx.VPC {
	return ccx.VPC{
		ID:            d.Id(),
		Name:          getString(d, "name"),
		CloudSpace:    getString(d, "cloud_space"),
		CloudProvider: getString(d, "cloud_provider"),
		Region:        getString(d, "cloud_region"),
		CidrIpv4Block: getString(d, "ipv4_cidr"),
	}
}

func vpcToSchema(v ccx.VPC, d *schema.ResourceData) error {
	d.SetId(v.ID)

	if err := d.Set("name", v.Name); err != nil {
		return err
	}

	if err := d.Set("cloud_space", v.CloudSpace); err != nil {
		return err
	}

	if err := d.Set("cloud_provider", v.CloudProvider); err != nil {
		return err
	}

	if err := d.Set("cloud_region", v.Region); err != nil {
		return err
	}

	if err := d.Set("ipv4_cidr", v.CidrIpv4Block); err != nil {
		return err
	}

	return nil
}

type VPC struct {
	svc ccx.VPCService
}

func (r *VPC) Name() string {
	return "ccx_vpc"
}

func (r *VPC) Configure(_ context.Context, cfg TerraformConfiguration) error {
	vpcCli := api.Vpcs(cfg.BaseURL, cfg.ClientID, cfg.ClientSecret)

	r.svc = vpcCli
	return nil
}

func (r *VPC) Schema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cloud_provider": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cloud_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipv4_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Create: r.Create,
		Read:   r.Read,
		Update: r.Update,
		Delete: r.Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func (r *VPC) Create(d *schema.ResourceData, _ any) error {
	ctx := context.Background()
	v := schemaToVPC(d)
	n, err := r.svc.Create(ctx, v)
	if err != nil {
		d.SetId("")
		return err
	}

	return vpcToSchema(*n, d)
}

func (r *VPC) Read(d *schema.ResourceData, _ any) error {
	ctx := context.Background()
	v := schemaToVPC(d)
	n, err := r.svc.Read(ctx, v.ID)
	if errors.Is(err, ccx.ResourceNotFoundErr) {
		d.SetId("")
		return nil
	} else if err != nil {
		return err
	}

	return vpcToSchema(*n, d)
}

func (r *VPC) Update(d *schema.ResourceData, _ any) error {
	ctx := context.Background()
	v := schemaToVPC(d)
	n, err := r.svc.Update(ctx, v)
	if err != nil {
		return err
	}

	return vpcToSchema(*n, d)
}

func (r *VPC) Delete(d *schema.ResourceData, _ any) error {
	ctx := context.Background()
	v := schemaToVPC(d)
	err := r.svc.Delete(ctx, v.ID)
	if err != nil {
		return err
	}
	return nil
}
