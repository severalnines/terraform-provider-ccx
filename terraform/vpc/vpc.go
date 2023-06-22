package vpc

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/severalnines/terraform-provider-ccx/ccx"
	chttp "github.com/severalnines/terraform-provider-ccx/http"
	"github.com/severalnines/terraform-provider-ccx/http/auth"
	vpcclient "github.com/severalnines/terraform-provider-ccx/http/vpc-client"
	"github.com/severalnines/terraform-provider-ccx/terraform"
)

var (
	_ ccx.TerraformResource = &Resource{}
)

func ToVpc(d *schema.ResourceData) ccx.VPC {
	return ccx.VPC{
		ID:            d.Id(),
		Name:          terraform.GetString(d, "name"),
		CloudSpace:    terraform.GetString(d, "cloud_space"),
		CloudProvider: terraform.GetString(d, "cloud_provider"),
		Region:        terraform.GetString(d, "cloud_region"),
		CidrIpv4Block: terraform.GetString(d, "ipv4_cidr"),
	}
}

func ToSchema(d *schema.ResourceData, v ccx.VPC) error {
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

type Resource struct {
	svc ccx.VPCService
}

func (r *Resource) Name() string {
	return "ccx_vpc"
}

func (r *Resource) Configure(_ context.Context, cfg ccx.TerraformConfiguration) error {
	authorizer := auth.New(cfg.ClientID, cfg.ClientSecret, chttp.BaseURL(cfg.BaseURL))
	vpcCli := vpcclient.New(authorizer, chttp.BaseURL(cfg.BaseURL))

	r.svc = vpcCli
	return nil
}

func (r *Resource) Schema() *schema.Resource {
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

func (r *Resource) Create(d *schema.ResourceData, _ any) error {
	ctx := context.Background()
	v := ToVpc(d)
	n, err := r.svc.Create(ctx, v)
	if err != nil {
		d.SetId("")
		return err
	}

	return ToSchema(d, *n)
}

func (r *Resource) Read(d *schema.ResourceData, _ any) error {
	ctx := context.Background()
	v := ToVpc(d)
	n, err := r.svc.Read(ctx, v.ID)
	if errors.Is(err, ccx.ResourceNotFoundErr) {
		d.SetId("")
		return nil
	} else if err != nil {
		return err
	}

	return ToSchema(d, *n)
}

func (r *Resource) Update(d *schema.ResourceData, _ any) error {
	ctx := context.Background()
	v := ToVpc(d)
	n, err := r.svc.Update(ctx, v)
	if err != nil {
		return err
	}

	return ToSchema(d, *n)
}

func (r *Resource) Delete(d *schema.ResourceData, _ any) error {
	ctx := context.Background()
	v := ToVpc(d)
	err := r.svc.Delete(ctx, v.ID)
	if err != nil {
		return err
	}
	return nil
}
