package vpc

import (
	"context"

	"github.com/hashicorp/terraform/helper/schema"
	ccxprov "github.com/severalnines/terraform-provider-ccx"
	chttp "github.com/severalnines/terraform-provider-ccx/http"
	"github.com/severalnines/terraform-provider-ccx/http/auth"
	vpcclient "github.com/severalnines/terraform-provider-ccx/http/vpc-client"
	"github.com/severalnines/terraform-provider-ccx/terraform"
)

var (
	_ ccxprov.TerraformResource = &Resource{}
)

func ToVpc(d *schema.ResourceData) ccxprov.VPC {
	return ccxprov.VPC{
		ID:            d.Id(),
		Name:          terraform.GetString(d, "name"),
		CloudSpace:    terraform.GetString(d, "cloud_space"),
		CloudProvider: terraform.GetString(d, "cloud_provider"),
		Region:        terraform.GetString(d, "cloud_region"),
		CidrIpv4Block: terraform.GetString(d, "ipv4_cidr"),
	}
}

func WriteSchema(d *schema.ResourceData, v ccxprov.VPC) {
	d.SetId(v.ID)
	d.Set("name", v.Name)
	d.Set("cloud_space", v.CloudSpace)
	d.Set("cloud_provider", v.CloudProvider)
	d.Set("cloud_region", v.Region)
	d.Set("ipv4_cidr", v.CidrIpv4Block)
}

type Resource struct {
	svc ccxprov.VPCService
}

func (r *Resource) Name() string {
	return "ccx_vpc"
}

func (r *Resource) Configure(_ context.Context, cfg ccxprov.TerraformConfiguration) error {
	// if p.Config.IsDevMode {
	// 	return r.ConfigureDev(p)
	// }

	authorizer := auth.New(cfg.ClientID, cfg.ClientSecret, chttp.BaseURL(cfg.BaseURL))
	vpcCli := vpcclient.New(authorizer, chttp.BaseURL(cfg.BaseURL))

	r.svc = vpcCli
	return nil
}

type mockdata struct {
	VPCs map[string]ccxprov.VPC `json:"vpcs"`
}

// func (r *Resource) ConfigureDev(p *terraform.Provider) error {
// 	var d mockdata
// 	if err := io.LoadData(p.Config.Mockfile, &d); err != nil {
// 		return err
// 	}
//
// 	r.svc = fakevpc.Instance(d.VPCs)
// 	return nil
// }

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
		return err
	}

	WriteSchema(d, *n)
	return nil
}

func (r *Resource) Read(d *schema.ResourceData, _ any) error {
	ctx := context.Background()
	v := ToVpc(d)
	n, err := r.svc.Read(ctx, v.ID)
	if err == ccxprov.ResourceNotFoundErr {
		return err
	} else if err != nil {
		return err
	}

	WriteSchema(d, *n)
	return nil
}

func (r *Resource) Update(d *schema.ResourceData, _ any) error {
	ctx := context.Background()
	v := ToVpc(d)
	n, err := r.svc.Update(ctx, v)
	if err != nil {
		return err
	}

	WriteSchema(d, *n)
	return nil
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
