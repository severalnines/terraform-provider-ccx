package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
)

type Datastore struct {
	svc        ccx.DatastoreService
	contentSvc ccx.ContentService
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
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "Replication type of the datastore",
				ForceNew:         true,
				DiffSuppressFunc: caseInsensitiveSuppressor,
			},
			"size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The size of the datastore ( int64 ). 1 or 3 nodes.",
				Default:     1,
			},
			"db_vendor": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Database Vendor",
				ForceNew:         true,
				DiffSuppressFunc: caseInsensitiveSuppressor,
			},
			"db_version": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "Database Version",
				ForceNew:         true,
				DiffSuppressFunc: caseInsensitiveSuppressor,
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "An optional list of tags",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cloud_provider": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Cloud provider name",
				ForceNew:         true,
				DiffSuppressFunc: caseInsensitiveSuppressor,
			},
			"cloud_region": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The region to set up the datastore",
				ForceNew:         true,
				DiffSuppressFunc: caseInsensitiveSuppressor,
			},
			"instance_size": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Instance type/flavor to use",
				DiffSuppressFunc: r.instanceSizeDiffSupressor,
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
				Default:     0,
			},
			"network_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Type of network: public/private",
				DiffSuppressFunc: caseInsensitiveSuppressor,
				ForceNew:         true,
			},
			"network_ha_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "High availability enabled or not",
				Default:     false,
			},
			"network_vpc_uuid": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "VPC to use if network_type is private",
				ForceNew:         true,
				DiffSuppressFunc: caseInsensitiveSuppressor,
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
			"notifications_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable or disable notifications. Default is false",
			},
			"notifications_emails": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "List of email addresses to send notifications to",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"maintenance_day_of_week": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Day of the week to run the maintenance. 1-7, 1 is Monday",
			},
			"maintenance_start_hour": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Hour of the day to start the maintenance. 0-23",
			},
			"maintenance_end_hour": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Hour of the day to end the maintenance. 0-23. Must be start_hour + 2",
			},
		},
		CreateContext: r.Create,
		ReadContext:   r.Read,
		UpdateContext: r.Update,
		DeleteContext: r.Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func (r *Datastore) Create(ctx context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	c, err := schemaToDatastore(d)

	if err != nil {
		return diag.FromErr(err)
	}

	n, err := r.svc.Create(ctx, c)
	if errors.Is(err, ccx.CreateFailedReadErr) && n != nil {
		d.SetId(n.ID)
		return diag.Errorf("creating stores: %s", err)
	} else if err != nil {
		d.SetId("")
		return diag.Errorf("creating stores: %s", err)
	}

	var errs []error

	if c.MaintenanceSettings != nil {
		if err := r.svc.SetMaintenanceSettings(ctx, n.ID, *c.MaintenanceSettings); err != nil {
			errs = append(errs, fmt.Errorf("%w setting: %w", ccx.MaintenanceSettingsErr, err))
		} else {
			n.MaintenanceSettings = c.MaintenanceSettings
		}
	}

	if len(c.DbParams) != 0 {
		if err := r.svc.SetParameters(ctx, n.ID, c.DbParams); err != nil {
			errs = append(errs, fmt.Errorf("%w setting: %w", ccx.ParametersErr, err))
		} else {
			n.DbParams = c.DbParams
		}
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
		return diag.Errorf("creating stores completed only partially: %s", errors.Join(errs...))
	}

	return nil
}

func (r *Datastore) Read(ctx context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	c, err := schemaToDatastore(d)

	if err != nil {
		return diag.FromErr(err)
	}

	n, err := r.svc.Read(ctx, c.ID)
	if errors.Is(err, ccx.ResourceNotFoundErr) {
		d.SetId("")
		return nil
	} else if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(schemaFromDatastore(*n, d))
}

func (r *Datastore) Update(ctx context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	c, err := schemaToDatastore(d)
	if err != nil {
		return diag.FromErr(err)
	}

	old, err := r.svc.Read(ctx, c.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	n := &c
	if d.HasChangesExcept("db_params", "firewall") {
		if n, err = r.svc.Update(ctx, *old, c); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("maintenance_day_of_week", "maintenance_start_hour", "maintenance_end_hour") {
		n.MaintenanceSettings = getMaintenanceSettings(d)
	}

	var errs []error

	if d.HasChange("db_params") {
		if err := r.svc.SetParameters(ctx, n.ID, c.DbParams); err != nil {
			errs = append(errs, fmt.Errorf("%w setting: %w", ccx.ParametersErr, err))
		} else {
			n.DbParams = c.DbParams
		}
	}

	if d.HasChange("firewall") {
		if err := r.svc.SetFirewallRules(ctx, n.ID, c.FirewallRules); err != nil {
			errs = append(errs, fmt.Errorf("%w: setting: %w", ccx.FirewallRulesErr, err))
		} else {
			n.FirewallRules = c.FirewallRules
		}
	}

	n.Notifications = getNotifications(d)

	if err := schemaFromDatastore(*n, d); err != nil {
		errs = append(errs, fmt.Errorf("setting schema: %w", err))
	}

	if len(errs) != 0 {
		return diag.Errorf("updating stores completed only partially: %s", errors.Join(errs...))
	}

	return nil
}

func (r *Datastore) Delete(ctx context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	c, err := schemaToDatastore(d)

	if err != nil {
		return diag.FromErr(err)
	}

	err = r.svc.Delete(ctx, c.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func (r *Datastore) instanceSizeDiffSupressor(_, oldValue, newValue string, d *schema.ResourceData) bool {
	if d.IsNewResource() || r.contentSvc == nil {
		// contentSvc might not have been initialized yet (configured not run by terraform)
		// also no need to check for new resources
		return false
	}

	ctx := context.Background()
	cloudProvider := getString(d, "cloud_provider")

	ok := checkInstanceSizeEquivalence(ctx, r.contentSvc, cloudProvider, oldValue, newValue)

	return ok
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

	if err = d.Set("db_version", c.DBVersion); err != nil {
		return err
	}

	if err = d.Set("type", defaultType(c.DBVendor, c.Type)); err != nil {
		return err
	}

	if err = setTags(d, "tags", c.Tags); err != nil {
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

	if err = d.Set("network_vpc_uuid", c.VpcUUID); err != nil {
		return err
	}

	if err = d.Set("network_ha_enabled", c.HAEnabled); err != nil {
		return err
	}

	if err = setStrings(d, "network_az", c.AvailabilityZones); err != nil {
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
