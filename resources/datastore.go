package resources

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

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
				DiffSuppressFunc: vendorSuppressor,
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
				Computed:    true,
				Description: "Volume type",
			},
			"volume_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Volume size",
			},
			"volume_iops": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Volume IOPS",
				Default:     0,
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
				Description:      "VPC to use",
				ForceNew:         true,
				DiffSuppressFunc: caseInsensitiveSuppressor,
			},
			"parameter_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Parameter group ID to use",
			},
			"network_az": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Network availability zones",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"firewall": {
				Type:             schema.TypeList,
				Optional:         true,
				Description:      "FirewallRule rules to allow",
				Elem:             (firewall{}).Schema(),
				DiffSuppressFunc: firewallDiffSupressor,
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
			"primary_url": {
				Type:        schema.TypeString,
				Optional:    false,
				Computed:    true,
				Description: "URL to the primary host(s)",
			},
			"primary_dsn": {
				Type:        schema.TypeString,
				Optional:    false,
				Computed:    true,
				Description: "DSN to the primary host(s)",
			},
			"replica_url": {
				Type:        schema.TypeString,
				Optional:    false,
				Computed:    true,
				Description: "URL to the replica host(s)",
			},
			"replica_dsn": {
				Type:        schema.TypeString,
				Optional:    false,
				Computed:    true,
				Description: "DSN to the replica host(s)",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    false,
				Computed:    true,
				Description: "Username to connect to the datastore",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    false,
				Computed:    true,
				Description: "Password to connect to the datastore",
			},
			"dbname": {
				Type:        schema.TypeString,
				Optional:    false,
				Computed:    true,
				Description: "Database name",
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

func validateCloud(cloudInstances map[string][]ccx.InstanceSize, c ccx.Datastore) error {
	prov, ok := cloudInstances[c.CloudProvider]
	if !ok {
		ls := make([]string, 0, len(cloudInstances))
		for k := range cloudInstances {
			ls = append(ls, k)
		}

		return fmt.Errorf("cloud provider %q not found. available cloud providers: %s", c.CloudProvider, strings.Join(ls, ", "))
	}

	ok = slices.ContainsFunc(prov, func(i ccx.InstanceSize) bool {
		return i.Code == c.InstanceSize || i.Type == c.InstanceSize
	})

	if !ok {
		ls := make([]string, 0, len(prov))
		for _, i := range prov {
			ls = append(ls, i.Code+" / "+i.Type)
		}

		return fmt.Errorf("instance size %q not found for provider %q. available sizes: %s", c.InstanceSize, c.CloudProvider, strings.Join(ls, ", "))
	}

	return nil
}

func validateDBVendor(vendors []ccx.DBVendorInfo, c ccx.Datastore) error {
	var vendor ccx.DBVendorInfo

	if i := slices.IndexFunc(vendors, func(info ccx.DBVendorInfo) bool {
		return info.Code == c.DBVendor
	}); i == -1 {
		ls := make([]string, 0, len(vendors))
		for _, v := range vendors {
			ls = append(ls, fmt.Sprintf("%q (%s)", v.Code, v.Name))
		}

		return fmt.Errorf("database vendor %q not found. available vendors: %s", c.DBVendor, strings.Join(ls, ", "))
	} else {
		vendor = vendors[i]
	}

	if c.DBVersion != "" {
		if i := slices.IndexFunc(vendor.Versions, func(v string) bool {
			return v == c.DBVersion
		}); i == -1 {
			return fmt.Errorf("database version %q not found for vendor %q. available versions: %s", c.DBVersion, c.DBVendor, strings.Join(vendor.Versions, ", "))
		}
	}

	if c.Type != "" {
		ok := slices.ContainsFunc(vendor.Types, func(t ccx.DBVendorInfoType) bool {
			return t.Code == c.Type
		})

		ls := make([]string, 0, len(vendor.Types))
		for _, t := range vendor.Types {
			ls = append(ls, fmt.Sprintf("%q (%s)", t.Code, t.Name))
		}

		if !ok {
			return fmt.Errorf("database type %q not found for vendor %q. available types: %s", c.Type, c.DBVendor, strings.Join(ls, ", "))
		}
	}

	return nil
}

func validateVolume(vendor string, volumeTypes []string, volumeType string, volumeSize uint64) error {
	if volumeType == "" {
		return errors.New("volume type is required")
	}

	if !slices.Contains(volumeTypes, volumeType) {
		return fmt.Errorf("volume type %q not found. available types: %s", volumeType, `"`+strings.Join(volumeTypes, `", "`)+`"`)
	}

	if (vendor == "redis" || vendor == "cache22") && volumeSize != 0 {
		return fmt.Errorf("volume_size is not supported for vendor %q", vendor)
	}

	return nil
}

func validateMaintenanceSettings(m *ccx.MaintenanceSettings) error {
	if m == nil {
		return nil
	}

	if m.DayOfWeek < 1 || m.DayOfWeek > 7 {
		return fmt.Errorf("maintenance_day_of week must be between 1 and 7: %d", m.DayOfWeek)
	}

	if m.StartHour < 0 || m.StartHour > 23 {
		return fmt.Errorf("maintenance_start_hour must be between 0 and 23: %d", m.StartHour)
	}

	if m.EndHour < 0 || m.EndHour > 23 {
		return fmt.Errorf("maintenance_end_hour must be between 0 and 23: %d", m.EndHour)
	}

	if (m.StartHour-m.EndHour) == 2 || (m.EndHour-m.StartHour) == 2 {
		return nil
	}

	s := m.StartHour
	if s == 0 {
		s = 24
	}

	e := m.EndHour
	if e == 0 {
		e = 24
	}

	if (e-s) != 2 && (s-e) != -2 {
		return fmt.Errorf("maintenance_end_hour must be start hour + 2: %d - %d", m.StartHour, m.EndHour)
	}

	return nil
}

func (r *Datastore) Create(ctx context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	c, err := schemaToDatastore(d)

	if err != nil {
		return diag.FromErr(err)
	}

	if err := validateMaintenanceSettings(c.MaintenanceSettings); err != nil {
		return diag.FromErr(fmt.Errorf("validating maintenance settings: %w", err))
	}

	cloudInstances, err := r.contentSvc.InstanceSizes(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("loading instance sizes: %w", err))
	}

	if err := validateCloud(cloudInstances, c); err != nil {
		return diag.FromErr(fmt.Errorf("validating cloud provider: %w", err))
	}

	vendors, err := r.contentSvc.DBVendors(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("loading db vendor information: %w", err))
	}

	if err := validateDBVendor(vendors, c); err != nil {
		return diag.FromErr(fmt.Errorf("validating db vendor: %w", err))
	}

	volumes, err := r.contentSvc.VolumeTypes(ctx, c.CloudProvider)
	if err != nil {
		return diag.FromErr(fmt.Errorf("loading volume types: %w", err))
	}

	if err := validateVolume(c.DBVendor, volumes, c.VolumeType, c.VolumeSize); err != nil {
		return diag.FromErr(fmt.Errorf("validating volume type: %w", err))
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
	if d.HasChangesExcept("firewall") {
		if n, err = r.svc.Update(ctx, *old, c); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("maintenance_day_of_week", "maintenance_start_hour", "maintenance_end_hour") {
		n.MaintenanceSettings = getMaintenanceSettings(d)

		if err := validateMaintenanceSettings(n.MaintenanceSettings); err != nil {
			return diag.FromErr(fmt.Errorf("validating maintenance settings: %w", err))
		}
	}

	var errs []error

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
	if err != nil && !errors.Is(err, ccx.ResourceNotFoundErr) {
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
	case "postgres":
		return "postgres_streaming"
	case "redis":
		return "redis"
	case "microsoft":
		return "mssql_ao_async"
	}

	return ""
}

func vendorFromAlias(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))

	switch s {
	case "mysql":
		return "percona"
	case "psql":
		return "postgres"
	}

	return s
}

func schemaToDatastore(d *schema.ResourceData) (ccx.Datastore, error) {
	c := ccx.Datastore{
		ID:               d.Id(),
		Name:             getString(d, "name"),
		Size:             getInt(d, "size"),
		DBVendor:         getString(d, "db_vendor"),
		DBVersion:        getString(d, "db_version"),
		Type:             getString(d, "type"),
		Tags:             getStrings(d, "tags"),
		CloudProvider:    getString(d, "cloud_provider"),
		CloudRegion:      getString(d, "cloud_region"),
		InstanceSize:     getString(d, "instance_size"),
		VolumeType:       getString(d, "volume_type"),
		VolumeSize:       uint64(getInt(d, "volume_size")),
		VolumeIOPS:       uint64(getInt(d, "volume_iops")),
		HAEnabled:        getBool(d, "network_ha_enabled"),
		ParameterGroupID: getString(d, "parameter_group"),
		VpcUUID:          getString(d, "network_vpc_uuid"),
	}

	if azs, hasAzs := getAzs(d); hasAzs && len(azs) == int(c.Size) {
		c.AvailabilityZones = azs
	} else if hasAzs {
		return c, fmt.Errorf("number of availability zones (%d) must match the size of the cluster (%d)", len(azs), c.Size)
	}

	firewalls, err := getFirewalls(d)
	if err != nil {
		return c, err
	}

	c.FirewallRules = firewalls

	c.DBVendor = vendorFromAlias(c.DBVendor)
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

	if err = d.Set("parameter_group", c.ParameterGroupID); err != nil {
		return err
	}

	if err = d.Set("network_ha_enabled", c.HAEnabled); err != nil {
		return err
	}

	if err = d.Set("primary_url", c.PrimaryUrl); err != nil {
		return err
	}

	if err = d.Set("primary_dsn", c.PrimaryDsn); err != nil {
		return err
	}

	if err = d.Set("replica_url", c.ReplicaUrl); err != nil {
		return err
	}

	if err = d.Set("replica_dsn", c.ReplicaDsn); err != nil {
		return err
	}

	if err = d.Set("username", c.Username); err != nil {
		return err
	}

	if err = d.Set("password", c.Password); err != nil {
		return err
	}

	if err = d.Set("dbname", c.DbName); err != nil {
		return err
	}

	if azs, ok := getAzs(d); ok { // do not set azs from upstream
		if err = setStrings(d, "network_az", azs); err != nil {
			return err
		}
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
