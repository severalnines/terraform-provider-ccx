package cluster

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ccxprov "github.com/severalnines/terraform-provider-ccx"
	chttp "github.com/severalnines/terraform-provider-ccx/http"
	"github.com/severalnines/terraform-provider-ccx/http/auth"
	clusterclient "github.com/severalnines/terraform-provider-ccx/http/cluster-client"
	"github.com/severalnines/terraform-provider-ccx/io"
	"github.com/severalnines/terraform-provider-ccx/terraform"
	fakecluster "github.com/severalnines/terraform-provider-ccx/tests/fake/cluster"
)

var (
	_ resource.Resource = &Resource{}
)

type Model struct {
	// Cluster is a database cluster
	ID                types.String `tfsdk:"id"`
	ClusterName       types.String `tfsdk:"cluster_name"`
	ClusterSize       types.Int64  `tfsdk:"cluster_size"`
	DBVendor          types.String `tfsdk:"db_vendor"`
	DBVersion         types.String `tfsdk:"db_version"`
	ClusterType       types.String `tfsdk:"cluster_type"`
	Tags              types.List   `tfsdk:"tags"`
	CloudSpace        types.String `tfsdk:"cloud_space"`
	CloudProvider     types.String `tfsdk:"cloud_provider"`
	CloudRegion       types.String `tfsdk:"cloud_region"`
	InstanceSize      types.String `tfsdk:"instance_size"` // "Tiny" ... "2X-Large"
	VolumeType        types.String `tfsdk:"volume_type"`
	VolumeSize        types.Int64  `tfsdk:"volume_size"`
	VolumeIOPS        types.Int64  `tfsdk:"volume_iops"`
	NetworkType       types.String `tfsdk:"network_type"` // public/private
	HAEnabled         types.Bool   `tfsdk:"network_ha_enabled"`
	VpcUUID           types.String `tfsdk:"network_vpc_uuid"`
	AvailabilityZones types.List   `tfsdk:"network_az"`
}

func ToModel(c ccxprov.Cluster) Model {
	return Model{
		ID:                types.StringValue(c.ID),
		ClusterName:       types.StringValue(c.ClusterName),
		ClusterSize:       types.Int64Value(c.ClusterSize),
		DBVendor:          types.StringValue(c.DBVendor),
		DBVersion:         types.StringValue(c.DBVersion),
		ClusterType:       types.StringValue(c.ClusterType),
		Tags:              terraform.StringsToList(c.Tags),
		CloudSpace:        types.StringValue(c.CloudSpace),
		CloudProvider:     types.StringValue(c.CloudProvider),
		CloudRegion:       types.StringValue(c.CloudRegion),
		InstanceSize:      types.StringValue(c.InstanceSize),
		VolumeType:        types.StringValue(c.VolumeType),
		VolumeSize:        types.Int64Value(c.VolumeSize),
		VolumeIOPS:        types.Int64Value(c.VolumeIOPS),
		NetworkType:       types.StringValue(c.NetworkType),
		HAEnabled:         types.BoolValue(c.HAEnabled),
		VpcUUID:           types.StringValue(c.VpcUUID),
		AvailabilityZones: terraform.StringsToList(c.AvailabilityZones),
	}
}

func FromModel(m Model) ccxprov.Cluster {
	return ccxprov.Cluster{
		ID:                m.ID.ValueString(),
		ClusterName:       m.ClusterName.ValueString(),
		ClusterSize:       m.ClusterSize.ValueInt64(),
		DBVendor:          m.DBVendor.ValueString(),
		DBVersion:         m.DBVersion.ValueString(),
		ClusterType:       m.ClusterType.ValueString(),
		Tags:              terraform.ListToStrings(m.Tags),
		CloudSpace:        m.CloudSpace.ValueString(),
		CloudProvider:     m.CloudProvider.ValueString(),
		CloudRegion:       m.CloudRegion.ValueString(),
		InstanceSize:      m.InstanceSize.ValueString(),
		VolumeType:        m.VolumeType.ValueString(),
		VolumeSize:        m.VolumeSize.ValueInt64(),
		VolumeIOPS:        m.VolumeIOPS.ValueInt64(),
		NetworkType:       m.NetworkType.ValueString(),
		HAEnabled:         m.HAEnabled.ValueBool(),
		VpcUUID:           m.VpcUUID.ValueString(),
		AvailabilityZones: terraform.ListToStrings(m.AvailabilityZones),
	}
}

type Resource struct {
	svc ccxprov.ClusterService
}

type mockdata struct {
	Clusters map[string]ccxprov.Cluster `json:"clusters"`
}

func (r *Resource) String() string {
	return "Cluster Resource"
}

func (r *Resource) Configure(ctx context.Context, p *terraform.Provider, _ diag.Diagnostics) error {
	if p.Config.IsDevMode {
		return r.ConfigureDev(p)
	}

	authorizer := auth.New(p.Config.ClientID, p.Config.ClientSecret, chttp.BaseURL(p.Config.BaseURL))
	clusterCli, err := clusterclient.New(ctx, authorizer, chttp.BaseURL(p.Config.BaseURL))
	if err != nil {
		return errors.Join(err, ccxprov.ResourcesLoadFailedErr)
	}

	r.svc = clusterCli
	return nil
}

func (r *Resource) ConfigureDev(p *terraform.Provider) error {
	var d mockdata
	if err := io.LoadData(p.Config.Mockfile, &d); err != nil {
		return err
	}

	r.svc = fakecluster.Instance(d.Clusters)
	return nil
}

func (r *Resource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "ccx_cluster"
}

func (r *Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the resource",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cluster_name": schema.StringAttribute{
				Description: "The name of the cluster",
				Required:    true,
			},
			"cluster_type": schema.StringAttribute{
				Description: "The name of the cluster",
				Optional:    true,
				Computed:    true,
			},
			"cluster_size": schema.Int64Attribute{
				Description: "The size of the cluster ( int64 ). 1 or 3 nodes.",
				Required:    true,
			},
			"db_vendor": schema.StringAttribute{
				Description: "The database vendor",
				Required:    true,
			},
			"db_version": schema.StringAttribute{
				Description: "Optional Database version",
				Required:    true,
			},
			"tags": schema.ListAttribute{
				Description: "An optional list of tags, represented as a key, value pair",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"cloud_provider": schema.StringAttribute{
				Description: "Cloud Provider",
				Required:    true,
			},
			"cloud_region": schema.StringAttribute{
				Description: "Cloud Region",
				Required:    true,
			},
			"cloud_space": schema.StringAttribute{
				Description: "Cloud Space",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"instance_size": schema.StringAttribute{
				Description: "The instance size",
				Required:    true,
				// Computed:    true,
			},
			"volume_type": schema.StringAttribute{
				Description: "",
				// Required:    true,
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString("gp2"),
			},
			"volume_size": schema.Int64Attribute{
				Description: "Volume Size",
				// Required:    true,
				Computed: true,
				Optional: true,
				Default:  int64default.StaticInt64(80),
			},
			"volume_iops": schema.Int64Attribute{
				Description: "",
				// Required:    true,
				Computed: true,
				Optional: true,
				Default:  int64default.StaticInt64(80),
			},
			"network_type": schema.StringAttribute{
				Description: "Network type: public or private",
				// Required:    true,
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString("public"),
			},
			"network_ha_enabled": schema.BoolAttribute{
				Description: "",
				// Required:    true,
				Computed: true,
				Optional: true,
				Default:  booldefault.StaticBool(false),
			},
			"network_vpc_uuid": schema.StringAttribute{
				Description: "",
				// Required:    true,
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
			},
			"network_az": schema.ListAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				// PlanModifiers: []planmodifier.List{
				// 	terraform.ListPlanModifier(),
				// },
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

	c := FromModel(m)
	newCluster, err := r.svc.Create(ctx, c)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create clusters", err.Error())
		return
	}

	resp.Diagnostics.AddWarning("new cluster info", newCluster.String())

	newModel := ToModel(*newCluster)
	resp.Diagnostics.Append(resp.State.Set(ctx, newModel)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var m Model

	resp.Diagnostics.Append(req.State.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}

	c := FromModel(m)
	readCluster, err := r.svc.Read(ctx, c.ID)
	if err == ccxprov.ResourceNotFoundErr {
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Unable to get clusters", err.Error())
		return
	}

	readModel := ToModel(*readCluster)
	resp.Diagnostics.Append(resp.State.Set(ctx, &readModel)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var m Model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}

	c := FromModel(m)
	newCluster, err := r.svc.Update(ctx, c)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update clusters", err.Error())
		return
	}

	newModel := ToModel(*newCluster)
	resp.Diagnostics.Append(resp.State.Set(ctx, newModel)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var m Model

	resp.Diagnostics.Append(req.State.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}

	c := FromModel(m)
	err := r.svc.Delete(ctx, c.ID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete clusters", err.Error())
	}
}
