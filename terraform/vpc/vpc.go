package vpc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ccxprov "github.com/severalnines/terraform-provider-ccx"
	chttp "github.com/severalnines/terraform-provider-ccx/http"
	"github.com/severalnines/terraform-provider-ccx/http/auth"
	vpcclient "github.com/severalnines/terraform-provider-ccx/http/vpc-client"
	"github.com/severalnines/terraform-provider-ccx/io"
	"github.com/severalnines/terraform-provider-ccx/terraform"
	fakevpc "github.com/severalnines/terraform-provider-ccx/tests/fake/vpc"
)

var (
	_ resource.Resource = &Resource{}
)

type Model struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	CloudSpace    types.String `tfsdk:"cloud_space"`
	CloudProvider types.String `tfsdk:"cloud_provider"`
	Region        types.String `tfsdk:"cloud_region"`
	CidrIpv4Block types.String `tfsdk:"cidr_ipv4_block"`
}

func ToModel(v ccxprov.VPC) Model {
	return Model{
		ID:            types.StringValue(v.ID),
		Name:          types.StringValue(v.Name),
		CloudSpace:    types.StringValue(v.CloudSpace),
		CloudProvider: types.StringValue(v.CloudProvider),
		Region:        types.StringValue(v.Region),
		CidrIpv4Block: types.StringValue(v.CidrIpv4Block),
	}
}

func FromModel(m Model) ccxprov.VPC {
	return ccxprov.VPC{
		ID:            m.ID.ValueString(),
		Name:          m.Name.ValueString(),
		CloudSpace:    m.CloudSpace.ValueString(),
		CloudProvider: m.CloudProvider.ValueString(),
		Region:        m.Region.ValueString(),
		CidrIpv4Block: m.CidrIpv4Block.ValueString(),
	}
}

type Resource struct {
	svc ccxprov.VPCService
}

func (r *Resource) String() string {
	return "VPC Resource"
}

func (r *Resource) Configure(_ context.Context, p *terraform.Provider, _ diag.Diagnostics) error {
	if p.Config.IsDevMode {
		return r.ConfigureDev(p)
	}

	authorizer := auth.New(p.Config.ClientID, p.Config.ClientSecret, chttp.BaseURL(p.Config.BaseURL))
	vpcCli := vpcclient.New(authorizer, chttp.BaseURL(p.Config.BaseURL))

	r.svc = vpcCli
	return nil
}

type mockdata struct {
	VPCs map[string]ccxprov.VPC `json:"vpcs"`
}

func (r *Resource) ConfigureDev(p *terraform.Provider) error {
	var d mockdata
	if err := io.LoadData(p.Config.Mockfile, &d); err != nil {
		return err
	}

	r.svc = fakevpc.Instance(d.VPCs)
	return nil
}

func (r *Resource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "ccx_vpc"
}

func (r *Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the VPC",
				Required:    true,
			},
			"cloud_provider": schema.StringAttribute{
				Description: "",
				Required:    true,
			},
			"cloud_region": schema.StringAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"cloud_space": schema.StringAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"cidr_ipv4_block": schema.StringAttribute{
				Description: "",
				Required:    true,
			},
		},
	}
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var m Model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}

	v := FromModel(m)
	newVPC, err := r.svc.Create(ctx, v)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create VPC", err.Error())
		return
	}

	newModel := ToModel(*newVPC)
	resp.Diagnostics.Append(resp.State.Set(ctx, newModel)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var m Model

	resp.Diagnostics.Append(req.State.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}

	c := FromModel(m)
	readVPC, err := r.svc.Read(ctx, c.ID)
	if err == ccxprov.ResourceNotFoundErr {
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Unable to get VPC", err.Error())
		return
	}

	readModel := ToModel(*readVPC)
	resp.Diagnostics.Append(resp.State.Set(ctx, &readModel)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var m Model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}

	v := FromModel(m)
	newVPC, err := r.svc.Update(ctx, v)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update VPC", err.Error())
		return
	}

	newModel := ToModel(*newVPC)
	resp.Diagnostics.Append(resp.State.Set(ctx, newModel)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var m Model

	resp.Diagnostics.Append(req.State.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}

	v := FromModel(m)
	err := r.svc.Delete(ctx, v.ID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete VPC", err.Error())
	}
}
