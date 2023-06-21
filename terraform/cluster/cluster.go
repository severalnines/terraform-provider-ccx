package cluster

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
	ccxprov "github.com/severalnines/terraform-provider-ccx"
	chttp "github.com/severalnines/terraform-provider-ccx/http"
	"github.com/severalnines/terraform-provider-ccx/http/auth"
	clusterclient "github.com/severalnines/terraform-provider-ccx/http/cluster-client"
	"github.com/severalnines/terraform-provider-ccx/terraform"
)

var (
	_ ccxprov.TerraformResource = &Resource{}
)

func ToCluster(d *schema.ResourceData) ccxprov.Cluster {
	c := ccxprov.Cluster{
		ID:                d.Id(),
		ClusterName:       terraform.GetString(d, "cluster_name"),
		ClusterSize:       terraform.GetInt(d, "cluster_size"),
		DBVendor:          terraform.GetString(d, "db_vendor"),
		DBVersion:         terraform.GetString(d, "db_version"),
		ClusterType:       terraform.GetString(d, "cluster_type"),
		Tags:              terraform.GetStrings(d, "tags"),
		CloudSpace:        terraform.GetString(d, "cloud_space"),
		CloudProvider:     terraform.GetString(d, "cloud_provider"),
		CloudRegion:       terraform.GetString(d, "cloud_region"),
		InstanceSize:      terraform.GetString(d, "instance_size"),
		VolumeType:        terraform.GetString(d, "volume_type"),
		VolumeSize:        terraform.GetInt(d, "volume_size"),
		VolumeIOPS:        terraform.GetInt(d, "volume_iops"),
		NetworkType:       terraform.GetString(d, "network_type"),
		HAEnabled:         terraform.GetBool(d, "network_ha_enabled"),
		VpcUUID:           terraform.GetString(d, "network_vpc_uuid"),
		AvailabilityZones: terraform.GetStrings(d, "network_az"),
	}

	return c
}

func ToSchema(d *schema.ResourceData, c ccxprov.Cluster) error {
	d.SetId(c.ID)
	var err error
	if err = d.Set("cluster_name", c.ClusterName); err != nil {
		return err
	}
	if err = d.Set("cluster_size", c.ClusterSize); err != nil {
		return err
	}
	if err = d.Set("db_vendor", c.DBVendor); err != nil {
		return err
	}
	if err = d.Set("db_version", c.DBVersion); err != nil {
		return err
	}
	if err = d.Set("cluster_type", c.ClusterType); err != nil {
		return err
	}
	if err = d.Set("tags", c.Tags); err != nil {
		return err
	}
	if err = d.Set("cloud_space", c.CloudSpace); err != nil {
		return err
	}
	if err = d.Set("cloud_provider", c.CloudProvider); err != nil {
		return err
	}
	if err = d.Set("cloud_region", c.CloudRegion); err != nil {
		return err
	}
	if err = d.Set("instance_size", c.InstanceSize); err != nil {
		return err
	}
	if err = d.Set("volume_type", c.VolumeType); err != nil {
		return err
	}
	if err = d.Set("volume_size", c.VolumeSize); err != nil {
		return err
	}
	if err = d.Set("volume_iops", c.VolumeIOPS); err != nil {
		return err
	}
	if err = d.Set("network_type", c.NetworkType); err != nil {
		return err
	}
	if err = d.Set("network_ha_enabled", c.HAEnabled); err != nil {
		return err
	}
	if err = d.Set("network_vpc_uuid", c.VpcUUID); err != nil {
		return err
	}
	if err = d.Set("network_az", c.AvailabilityZones); err != nil {
		return err
	}
	return nil
}

type Resource struct {
	svc ccxprov.ClusterService
}

func (r *Resource) Name() string {
	return "ccx_cluster"
}

func (r *Resource) Schema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the resource",
				// ValidateFunc: validateName,
			},
			"cluster_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of the resource",
			},
			"cluster_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The size of the cluster ( int64 ). 1 or 3 nodes.",
			},
			"db_vendor": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Database Vendor",
			},
			"db_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Database Version",
			},
			"tags": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				Description:      "An optional list of tags, represented as a key, value pair",
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: terraform.NonNewSuppressor,
			},
			"cloud_provider": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional list of tags, represented as a key, value pair",
			},
			"cloud_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region to set up the cluster",
			},
			"cloud_space": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cloud space information",
			},
			"instance_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional list of tags, represented as a key, value pair",
			},
			"volume_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Volume type",
			},
			"volume_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Volume size",
			},
			"volume_iops": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Volume IOPS",
			},
			"network_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Type of network: public/private",
				DiffSuppressFunc: terraform.NonNewSuppressor,
			},
			"network_ha_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "High availability enabled or not",
			},
			"network_vpc_uuid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "VPC to use if network_type is private",
			},
			"network_az": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Network availability zones",
				Elem:        &schema.Schema{Type: schema.TypeString},
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

func (r *Resource) Configure(ctx context.Context, cfg ccxprov.TerraformConfiguration) error {
	// if p.Config.IsDevMode {
	// 	return r.ConfigureDev(p)
	// }

	authorizer := auth.New(cfg.ClientID, cfg.ClientSecret, chttp.BaseURL(cfg.BaseURL))
	clusterCli, err := clusterclient.New(ctx, authorizer, chttp.BaseURL(cfg.BaseURL))
	if err != nil {
		return errors.Join(err, ccxprov.ResourcesLoadFailedErr)
	}

	r.svc = clusterCli
	return nil
}

// func (r *Resource) ConfigureDev(p *terraform.Provider) error {
// 	var d mockdata
// 	if err := io.LoadData(p.Config.Mockfile, &d); err != nil {
// 		return err
// 	}
//
// 	r.svc = fakecluster.Instance(d.Clusters)
// 	return nil
// }

func (r *Resource) Create(d *schema.ResourceData, _ any) error {
	ctx := context.Background()
	c := ToCluster(d)
	n, err := r.svc.Create(ctx, c)
	if err != nil {
		d.SetId("")
		return err
	}

	return ToSchema(d, *n)
}

func (r *Resource) Read(d *schema.ResourceData, _ any) error {
	ctx := context.Background()
	c := ToCluster(d)
	n, err := r.svc.Read(ctx, c.ID)
	if errors.Is(err, ccxprov.ResourceNotFoundErr) {
		d.SetId("")
		return nil
	} else if err != nil {
		return err
	}

	return ToSchema(d, *n)
}

func (r *Resource) Update(d *schema.ResourceData, _ any) error {
	ctx := context.Background()
	c := ToCluster(d)
	n, err := r.svc.Update(ctx, c)
	if err != nil {
		return err
	}

	return ToSchema(d, *n)
}

func (r *Resource) Delete(d *schema.ResourceData, _ any) error {
	ctx := context.Background()
	c := ToCluster(d)
	err := r.svc.Delete(ctx, c.ID)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
