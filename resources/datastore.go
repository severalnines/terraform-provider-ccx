package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx/api"
	"github.com/severalnines/terraform-provider-ccx/internal/lib"
)

var (
	_ TerraformResource = &Datastore{}
)

func schemaToDatastore(d *schema.ResourceData) (ccx.Datastore, error) {
	c := ccx.Datastore{
		ID:                d.Id(),
		Name:              getString(d, "name"),
		Size:              getInt(d, "size"),
		DBVendor:          getString(d, "db_vendor"),
		DBVersion:         getString(d, "db_version"),
		Type:              getString(d, "type"),
		Tags:              getStrings(d, "tags"),
		CloudProvider:     getString(d, "cloud_provider"),
		CloudRegion:       getString(d, "cloud_region"),
		InstanceSize:      getString(d, "instance_size"),
		VolumeType:        getString(d, "volume_type"),
		VolumeSize:        uint64(getInt(d, "volume_size")),
		VolumeIOPS:        uint64(getInt(d, "volume_iops")),
		NetworkType:       getString(d, "network_type"),
		HAEnabled:         getBool(d, "network_ha_enabled"),
		VpcUUID:           getString(d, "network_vpc_uuid"),
		AvailabilityZones: getStrings(d, "network_az"),
	}

	dbparams := getMapString(d, "db_params")
	c.DbParams = dbparams

	firewalls, err := getFirewalls(d)
	if err != nil {
		return c, err
	}

	c.FirewallRules = firewalls

	c.Type = defaultType(c.DBVendor, c.Type)

	c.Notifications = getNotifications(d)
	c.MaintenanceSettings = getMaintenanceSettings(d)

	return c, nil
}

func schemaFromDatastore(c ccx.Datastore, d *schema.ResourceData) error {
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
	if getString(d, "db_version") != "" {
		if err = d.Set("db_version", c.DBVersion); err != nil {
			return err
		}
	}
	if getString(d, "type") != "" || c.Type != defaultType(c.DBVendor, c.Type) {
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

	if err = d.Set("db_params", c.DbParams); err != nil {
		return err
	}

	if err = setFirewalls(d, c.FirewallRules); err != nil {
		return err
	}

	if err = setNotifications(d, c.Notifications); err != nil {
		return err
	}

	if c.MaintenanceSettings != nil {
		if err = setMaintenanceSettings(d, *c.MaintenanceSettings); err != nil {
			return err
		}
	}

	return nil
}

type Datastore struct {
	svc ccx.DatastoreService
}

func (r *Datastore) Name() string {
	return "ccx_datastore"
}

func (r *Datastore) Configure(ctx context.Context, cfg TerraformConfiguration) error {
	httpcli := lib.NewHttpClient(ctx, "datastore", cfg.BaseURL, cfg.ClientID, cfg.ClientSecret, cfg.Logpath)

	svc, err := api.Datastores(ctx, httpcli, cfg.Timeout)
	if err != nil {
		return errors.Join(err, ccx.ResourcesLoadFailedErr)
	}

	r.svc = svc

	return nil
}

func (r *Datastore) Schema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the datastore",
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
				DiffSuppressFunc: nonNewSuppressor,
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
				DiffSuppressFunc: nonNewSuppressor,
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
			"db_params": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Database parameters",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"firewall": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "FirewallRule rules to allow",
				Elem:        (firewall{}).Schema(),
			},
			// "notifications": {
			// 	Type:        schema.TypeMap,
			// 	Optional:    true,
			// 	Description: "Notification settings",
			// 	Elem:        (notifications{}).Schema(),
			// },
			"notifications_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable or disable notifications. Default is false",
			},
			"notifications_emails": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of email addresses to send notifications to",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			// "maintenance_settings": {
			// 	Type:        schema.TypeMap,
			// 	Optional:    true,
			// 	Description: "Maintenance settings",
			// 	Elem:        (maintenance{}).Schema(),
			// },
			"maintenance_day_of_week": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Day of the week to run the maintenance. 1-7, 1 is Monday",
			},
			"maintenance_start_hour": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Hour of the day to start the maintenance. 0-23",
			},
			"maintenance_end_hour": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Hour of the day to end the maintenance. 0-23. Must be start_hour + 2",
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

func (r *Datastore) Create(d *schema.ResourceData, _ any) error {
	ctx := context.Background()
	c, err := schemaToDatastore(d)

	if err != nil {
		return err
	}

	n, err := r.svc.Create(ctx, c)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("creating stores: %w", err)
	}

	var errs []error

	if len(c.DbParams) != 0 {
		if err := r.svc.SetParameters(ctx, n.ID, c.DbParams); err != nil {
			errs = append(errs, fmt.Errorf("%w setting: %w", ccx.ParametersErr, err))
		} else {
			n.DbParams = c.DbParams
		}
	}

	if err := r.svc.SetFirewallRules(ctx, n.ID, c.FirewallRules); err != nil {
		errs = append(errs, fmt.Errorf("%w: setting: %w", ccx.FirewallRulesErr, err))
	} else {
		n.FirewallRules = c.FirewallRules
	}

	if err := schemaFromDatastore(*n, d); err != nil {
		errs = append(errs, fmt.Errorf("setting schema: %w", err))
	}

	if len(errs) != 0 {
		return fmt.Errorf("creating stores completed only partially: %w", errors.Join(errs...))
	}

	return nil
}

func (r *Datastore) Read(d *schema.ResourceData, _ any) error {
	ctx := context.Background()
	c, err := schemaToDatastore(d)

	if err != nil {
		return err
	}

	n, err := r.svc.Read(ctx, c.ID)
	if errors.Is(err, ccx.ResourceNotFoundErr) {
		d.SetId("")
		return nil
	} else if err != nil {
		return err
	}

	return schemaFromDatastore(*n, d)
}

func (r *Datastore) Update(d *schema.ResourceData, _ any) error {
	ctx := context.Background()

	c, err := schemaToDatastore(d)
	if err != nil {
		return err
	}

	old, err := r.svc.Read(ctx, c.ID)
	if err != nil {
		return err
	}

	n, err := r.svc.Update(ctx, *old, c)
	if err != nil {
		return err
	}

	var errs []error

	if err := r.svc.SetParameters(ctx, n.ID, c.DbParams); err != nil {
		errs = append(errs, fmt.Errorf("%w setting: %w", ccx.ParametersErr, err))
	} else {
		n.DbParams = c.DbParams
	}

	if len(c.FirewallRules) != 0 {
		if err := r.svc.SetFirewallRules(ctx, n.ID, c.FirewallRules); err != nil {
			errs = append(errs, fmt.Errorf("%w: setting: %w", ccx.FirewallRulesErr, err))
		} else {
			n.FirewallRules = c.FirewallRules
		}
	}

	if err := schemaFromDatastore(*n, d); err != nil {
		errs = append(errs, fmt.Errorf("setting schema: %w", err))
	}

	if len(errs) != 0 {
		return fmt.Errorf("updating stores completed only partially: %w", errors.Join(errs...))
	}

	return nil
}

func (r *Datastore) Delete(d *schema.ResourceData, _ any) error {
	ctx := context.Background()
	c, err := schemaToDatastore(d)

	if err != nil {
		return err
	}

	err = r.svc.Delete(ctx, c.ID)
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
