package datastore

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/severalnines/terraform-provider-ccx/ccx"
	chttp "github.com/severalnines/terraform-provider-ccx/http"
	"github.com/severalnines/terraform-provider-ccx/http/auth"
	datastoreclient "github.com/severalnines/terraform-provider-ccx/http/datastore-client"
	"github.com/severalnines/terraform-provider-ccx/terraform"
)

var (
	_ ccx.TerraformResource = &Resource{}
)

func ToDatastore(d *schema.ResourceData) ccx.Datastore {
	c := ccx.Datastore{
		ID:                d.Id(),
		Name:              terraform.GetString(d, "name"),
		Size:              terraform.GetInt(d, "size"),
		DBVendor:          terraform.GetString(d, "db_vendor"),
		DBVersion:         terraform.GetString(d, "db_version"),
		Type:              terraform.GetString(d, "type"),
		Tags:              terraform.GetStrings(d, "tags"),
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
	c.Type = defaultType(c.DBVendor, c.Type)

	return c
}

func ToSchema(d *schema.ResourceData, c ccx.Datastore) error {
	d.SetId(c.ID)
	var err error
	if err = d.Set("name", c.Name); err != nil {
		return err
	}
	if err = d.Set("size", c.Size); err != nil {
		return err
	}
	if err = d.Set("db_vendor", c.DBVendor); err != nil {
		return err
	}
	if terraform.GetString(d, "db_version") != "" {
		if err = d.Set("db_version", c.DBVersion); err != nil {
			return err
		}
	}
	if terraform.GetString(d, "type") != "" || c.Type != defaultType(c.DBVendor, c.Type) {
		if err = d.Set("type", defaultType(c.DBVendor, c.Type)); err != nil {
			return err
		}
	}
	if err = d.Set("tags", c.Tags); err != nil {
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
	svc ccx.DatastoreService
}

func (r *Resource) Name() string {
	return "ccx_datastore"
}

func (r *Resource) Schema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the datastore",
				// ValidateFunc: validateName,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Replication type of the datastore",
			},
			"size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The size of the datastore ( int64 ). 1 or 3 nodes.",
				Default:     1,
			},
			"db_vendor": {
				Type:        schema.TypeString,
				Required:    true,
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
				Required:    true,
				Description: "Cloud provider name",
			},
			"cloud_region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The region to set up the datastore",
			},
			"instance_size": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance type/flavor to use",
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

func (r *Resource) Configure(ctx context.Context, cfg ccx.TerraformConfiguration) error {
	authorizer := auth.New(cfg.ClientID, cfg.ClientSecret, chttp.BaseURL(cfg.BaseURL))
	datastoreCli, err := datastoreclient.New(ctx, authorizer, chttp.BaseURL(cfg.BaseURL))
	if err != nil {
		return errors.Join(err, ccx.ResourcesLoadFailedErr)
	}

	r.svc = datastoreCli
	return nil
}

func (r *Resource) Create(d *schema.ResourceData, _ any) error {
	ctx := context.Background()
	c := ToDatastore(d)
	n, err := r.svc.Create(ctx, c)
	if err != nil {
		d.SetId("")
		return err
	}

	return ToSchema(d, *n)
}

func (r *Resource) Read(d *schema.ResourceData, _ any) error {
	ctx := context.Background()
	c := ToDatastore(d)
	n, err := r.svc.Read(ctx, c.ID)
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
	c := ToDatastore(d)
	n, err := r.svc.Update(ctx, c)
	if err != nil {
		return err
	}

	return ToSchema(d, *n)
}

func (r *Resource) Delete(d *schema.ResourceData, _ any) error {
	ctx := context.Background()
	c := ToDatastore(d)
	err := r.svc.Delete(ctx, c.ID)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func defaultType(vendor, dbType string) string {
	if dbType != "" {
		return dbType
	}
	switch vendor {
	case "mariadb", "percona":
		return "replication"
	case "psql", "postgres":
		return "postgres_streaming"
	case "redis":
		return "redis"
	}
	return ""
}
