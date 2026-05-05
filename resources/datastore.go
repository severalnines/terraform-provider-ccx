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

const datastoreDoc = `
Datastores are a CCX resource, and represents one or more servers working together to host a database system.

For full documenation about CCX see https://docs.severalnines.com/ccx/user/Index/`

type Datastore struct {
	svc        ccx.DatastoresService
	contentSvc ccx.ContentService
	pgSvc      ccx.ParameterGroupsService
}

func (r *Datastore) Schema() *schema.Resource {
	return &schema.Resource{
		Description: datastoreDoc,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the datastore. This is just for your reference, and can be changed later.",
			},
			"type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "Replication type of the datastore. This depends on the db_vendor, e.g. `replication` is the default type for MySQL, MariaDB and PostgreSQL.",
				ForceNew:         true,
				DiffSuppressFunc: caseInsensitiveSuppressor,
			},
			"size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of nodes in the datastore. While a single node is allowed, there will be no redundancy. For multi-master datastores there must be an odd number of nodes.",
				Default:     1,
			},
			"db_vendor": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Database vendor. Allowed values depend on the CCX instance. Commonly available vendors are `mysql` and `postgres`.",
				ForceNew:         true,
				DiffSuppressFunc: vendorSuppressor,
			},
			"db_version": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "Version of the database system. Refer to the CCX instance to find versions available for each vendor.",
				ForceNew:         true,
				DiffSuppressFunc: caseInsensitiveSuppressor,
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "An optional list of tags to identify the datastore. These are are for your own use, and can be any strings.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cloud_provider": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Cloud provider name, e.g. `aws`. Refer to the CCX instance to find which clouds are available.",
				ForceNew:         true,
				DiffSuppressFunc: caseInsensitiveSuppressor,
			},
			"cloud_region": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The region to set up the datastore, within the chosen cloud. E.g. `us-east-1` in AWS.",
				ForceNew:         true,
				DiffSuppressFunc: caseInsensitiveSuppressor,
			},
			"instance_size": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Instance type/flavor to use. Refer to the CCX instance to find which instances types are available in a cloud and region.",
				DiffSuppressFunc: r.instanceSizeDiffSupressor,
			},
			"volume_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Volume type, for that will be used as root and data disks as required.",
			},
			"volume_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Volume size, i.e. how much data storage should be initally allocated. This can be changed later, or autoscaled.",
			},
			"volume_iops": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Volume IOPS defines the performance of the disks used for data storage. This is not always configurable, and allowable values depend on the volume type.",
				Default:     0,
			},
			"network_ha_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "This option does nothing directly, but if HA is set to true then CCX will require that availability zones are specified.",
				Default:     false,
			},
			"network_vpc_uuid": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "ID of a VPC, in which the cluster will be deployed.",
				ForceNew:         true,
				DiffSuppressFunc: caseInsensitiveSuppressor,
			},
			"parameter_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Parameter group ID to use. Parameter groups are another CCX resource, and contain a values for configuratable settings with the database system.",
			},
			"network_az": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Network availability zones. This can be 1) omitted for auto-allocation, 2) a single string, for placing all nodes in the same zone, 3) as many strings as the intended size of the cluster, to place each node separately. The values depend on the chosen cloud and region.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"firewall": {
				Type:             schema.TypeList,
				Optional:         true,
				Description:      "Firewall rules allow access to the database system from the internet. If there are no rules then all access is blocked. Each rule is a human-readable name and a CIDR, allowing access from a block of IP addresses.",
				Elem:             (firewall{}).Schema(),
				DiffSuppressFunc: firewallDiffSupressor,
			},
			"notifications_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable or disable notifications. Default is false.",
			},
			"notifications_emails": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "List of email addresses to send notifications to.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"maintenance_day_of_week": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Day of the week when maintenance tasks can be run. 1-7, 1 is Monday.",
			},
			"maintenance_start_hour": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Hour of the day when maintenance tasks can be run, on the chosen day. 0-23.",
			},
			"maintenance_end_hour": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Hour of the day when it is no longer appropriate to run maintenance tasks. 0-23. This must be approximtely maintenance_start_hour + 2.",
			},
			"primary_url": {
				Type:        schema.TypeString,
				Optional:    false,
				Computed:    true,
				Description: "URL to the primary host(s). This is a DNS name, which will resolve to one or more hosts.",
			},
			"primary_dsn": {
				Type:        schema.TypeString,
				Optional:    false,
				Computed:    true,
				Description: "DSN (data source name) to the primary host(s). This is the information that is needed to connect to the cluster - the format depends on the vendor.",
			},
			"replica_url": {
				Type:        schema.TypeString,
				Optional:    false,
				Computed:    true,
				Description: "URL to the replica host(s). This is a DNS name, which will resolve to zero or more hosts.",
			},
			"replica_dsn": {
				Type:        schema.TypeString,
				Optional:    false,
				Computed:    true,
				Description: "DSN (data source name) to the replica host(s). This is the information that is needed to connect to the cluster - the format depends on the vendor.",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    false,
				Computed:    true,
				Description: "Username to connect to the datastore - this represents the default user which is automatically created, but other users can be created later.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    false,
				Computed:    true,
				Description: "Password to connect to the datastore - this represents the default user which is automatically created, but other users can be created later.",
			},
			"dbname": {
				Type:        schema.TypeString,
				Optional:    false,
				Computed:    true,
				Description: "Name of the default database, which is automatically created when the cluster is created.",
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

func validateInstanceSizes(cloudInstances map[string][]ccx.InstanceSize, c ccx.Datastore) error {
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

func validateDB(vendors []ccx.DBVendorInfo, dbVendor, dbVersion, dbType string) error {
	var vendor ccx.DBVendorInfo

	if i := slices.IndexFunc(vendors, func(info ccx.DBVendorInfo) bool {
		return info.Code == dbVendor
	}); i == -1 {
		ls := make([]string, 0, len(vendors))
		for _, v := range vendors {
			ls = append(ls, fmt.Sprintf("%q (%s)", v.Code, v.Name))
		}

		return fmt.Errorf("database vendor %q not found. available vendors: %s", dbVendor, strings.Join(ls, ", "))
	} else {
		vendor = vendors[i]
	}

	if dbVersion != "" {
		if i := slices.IndexFunc(vendor.Versions, func(v string) bool {
			return v == dbVersion
		}); i == -1 {
			return fmt.Errorf("database version %q not found for vendor %q. available versions: %s", dbVersion, dbVendor, strings.Join(vendor.Versions, ", "))
		}
	}

	if dbType != "" {
		ok := slices.ContainsFunc(vendor.Types, func(t ccx.DBVendorInfoType) bool {
			return t.Code == dbType
		})

		ls := make([]string, 0, len(vendor.Types))
		for _, t := range vendor.Types {
			ls = append(ls, fmt.Sprintf("%q (%s)", t.Code, t.Name))
		}

		if !ok {
			return fmt.Errorf("database type %q not found for vendor %q. available types: %s", dbType, dbVendor, strings.Join(ls, ", "))
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

	if (vendor == "redis" || vendor == "cache22" || vendor == "valkey") && volumeSize != 0 {
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

func validateParameterGroupForStore(ctx context.Context, svc ccx.ParameterGroupsService, c ccx.Datastore, groupId string) error {
	if c.ParameterGroupID == "" {
		return nil
	}

	p, err := svc.Read(ctx, c.ParameterGroupID)
	if err != nil {
		return fmt.Errorf("loading parameter group %q: %w", groupId, err)
	}

	dbType := c.Type
	if dbType == "" {
		dbType = defaultType(c.DBVendor, c.Type)
	}

	v1 := vendorFromAlias(c.DBVendor)
	v2 := vendorFromAlias(p.DatabaseVendor)

	ok := v1 == v2 && c.DBVersion == p.DatabaseVersion && strings.EqualFold(dbType, p.DatabaseType)

	if !ok {
		return fmt.Errorf(
			"parameter_group %q with (db_vendor=%q, db_version=%q, db_type=%q) does not match the datastore %q (db_vendor=%q, db_version=%q, db_type=%q)",
			groupId, p.DatabaseVendor, p.DatabaseVersion, p.DatabaseType, c.ID, c.DBVendor, c.DBVersion, dbType,
		)
	}

	return nil
}

func (r *Datastore) Create(ctx context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	c, err := datastoreFromSchema(d)
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

	if err := validateInstanceSizes(cloudInstances, c); err != nil {
		return diag.FromErr(fmt.Errorf("validating cloud provider: %w", err))
	}

	vendors, err := r.contentSvc.DBVendors(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("loading db vendor information: %w", err))
	}

	if err := validateDB(vendors, c.DBVendor, c.DBVersion, c.Type); err != nil {
		return diag.FromErr(fmt.Errorf("validating db vendor: %w", err))
	}

	volumes, err := r.contentSvc.VolumeTypes(ctx, c.CloudProvider)
	if err != nil {
		return diag.FromErr(fmt.Errorf("loading volume types: %w", err))
	}

	if err := validateVolume(c.DBVendor, volumes, c.VolumeType, c.VolumeSize); err != nil {
		return diag.FromErr(fmt.Errorf("validating volume type: %w", err))
	}

	if c.ParameterGroupID != "" {
		if err := validateParameterGroupForStore(ctx, r.pgSvc, c, c.ParameterGroupID); err != nil {
			return diag.FromErr(fmt.Errorf("validating parameter group: %w", err))
		}
	}

	var errs []error

	n, err := r.svc.Create(ctx, c)
	if errors.Is(err, ccx.ErrCreateFailedRead) && n != nil {
		d.SetId(n.ID)
		return diag.Errorf("creating stores: %s", err)
	} else if err != nil {
		d.SetId("")
		return diag.Errorf("creating stores: %s", err)
	}

	if c.ParameterGroupID != "" {
		err = r.svc.ApplyParameterGroup(ctx, n.ID, c.ParameterGroupID)
		if err != nil {
			errs = append(errs, fmt.Errorf("applying database parameters %q failed: %w", c.ParameterGroupID, err))
		} else {
			n.ParameterGroupID = c.ParameterGroupID
		}
	}

	if c.MaintenanceSettings != nil {
		if err := r.svc.SetMaintenanceSettings(ctx, n.ID, *c.MaintenanceSettings); err != nil {
			errs = append(errs, fmt.Errorf("%w setting: %w", ccx.ErrMaintenanceSettings, err))
		} else {
			n.MaintenanceSettings = c.MaintenanceSettings
		}
	}

	if len(c.FirewallRules) != 0 {
		if err := r.svc.SetFirewallRules(ctx, n.ID, c.FirewallRules); err != nil {
			errs = append(errs, fmt.Errorf("%w: setting: %w", ccx.ErrFirewallRules, err))
		} else {
			n.FirewallRules = c.FirewallRules
		}
	}

	err = fillSchemaFromDatastore(*n, d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("setting schema: %w", err))
	}

	if len(errs) != 0 {
		return diag.Errorf("creating stores completed only partially: %s", errors.Join(errs...))
	}

	return nil
}

func (r *Datastore) Read(ctx context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	c, err := datastoreFromSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}

	n, err := r.svc.Read(ctx, c.ID)
	if errors.Is(err, ccx.ErrResourceNotFound) {
		d.SetId("")
		return nil
	} else if err != nil {
		return diag.FromErr(err)
	}

	err = fillSchemaFromDatastore(*n, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func (r *Datastore) Update(ctx context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	c, err := datastoreFromSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}

	old, err := r.svc.Read(ctx, c.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	n := &c

	if d.HasChanges("maintenance_day_of_week", "maintenance_start_hour", "maintenance_end_hour") {
		n.MaintenanceSettings = getMaintenanceSettings(d)

		if err := validateMaintenanceSettings(n.MaintenanceSettings); err != nil {
			return diag.FromErr(fmt.Errorf("validating maintenance settings: %w", err))
		}
	}

	if old.VolumeSize > c.VolumeSize {
		return diag.Errorf("decreasing volume_size is not supported, from %dGB to %dGB", old.VolumeSize, c.VolumeSize)
	} else if old.VolumeSize != c.VolumeSize && (old.VolumeSize+10) >= c.VolumeSize {
		return diag.Errorf("when increasing volume_size, the new volume_size must be at least old.volume_size+10GB. current volume_size is %dGB, new volume_size is %dGB. new volume_size must be atleast %dGB", old.VolumeSize, c.VolumeSize, old.VolumeSize+10)
	}

	var errs []error

	if d.HasChangesExcept("firewall", "parameter_group") {
		if n, err = r.svc.Update(ctx, *old, c); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("parameter_group") {
		if err := r.svc.ApplyParameterGroup(ctx, n.ID, c.ParameterGroupID); err != nil {
			errs = append(errs, fmt.Errorf("%w: %w", ccx.ErrApplyParameterGroup, err))
		} else {
			n.ParameterGroupID = c.ParameterGroupID
		}
	}

	if d.HasChange("firewall") {
		if err := r.svc.SetFirewallRules(ctx, n.ID, c.FirewallRules); err != nil {
			errs = append(errs, fmt.Errorf("%w: %w", ccx.ErrFirewallRules, err))
		} else {
			n.FirewallRules = c.FirewallRules
		}
	}

	// WHY? didn't we already read the whole schema?
	// n.Notifications = getNotifications(d)

	err = fillSchemaFromDatastore(*n, d)
	if err != nil {
		errs = append(errs, fmt.Errorf("setting schema: %w", err))
	}

	if len(errs) != 0 {
		return diag.Errorf("updating stores completed only partially: %s", errors.Join(errs...))
	}

	return nil
}

func (r *Datastore) Delete(ctx context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	c, err := datastoreFromSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = r.svc.Delete(ctx, c.ID)
	if err != nil && !errors.Is(err, ccx.ErrResourceNotFound) {
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
	case "valkey":
		return "valkey"
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

func datastoreFromSchema(d *schema.ResourceData) (ccx.Datastore, error) {
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

func fillSchemaFromDatastore(c ccx.Datastore, d *schema.ResourceData) error {
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
