package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
)

type ParameterGroup struct {
	svc        ccx.ParameterGroupService
	contentSvc ccx.ContentService
}

func (r *ParameterGroup) Schema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this parameter group",
			},
			"database_vendor": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Database vendor for which this parameter group is applicable",
				DiffSuppressFunc: vendorSuppressor,
			},
			"database_version": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Database version for which this parameter group is applicable",
				DiffSuppressFunc: caseInsensitiveSuppressor,
			},
			"database_type": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Database type for which this parameter group is applicable",
				DiffSuppressFunc: caseInsensitiveSuppressor,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of this parameter group",
			},
			"parameters": {
				Type:        schema.TypeMap,
				Required:    true,
				Description: "Parameters for this parameter group",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

func (r *ParameterGroup) Create(ctx context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	p, err := schemaToParameterGroup(d)
	if err != nil {
		return diag.FromErr(err)
	}

	vendors, err := r.contentSvc.DBVendors(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("loading db vendor information: %w", err))
	}

	if err := validateDb(vendors, p.DatabaseVendor, p.DatabaseVersion, p.DatabaseType); err != nil {
		return diag.FromErr(fmt.Errorf("validating db vendor: %w", err))
	}

	n, err := r.svc.Create(ctx, p)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := schemaFromParameterGroup(*n, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func (r *ParameterGroup) Read(ctx context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	id := d.Id()

	p, err := r.svc.Read(ctx, id)
	if errors.Is(err, ccx.ErrResourceNotFound) {
		d.SetId("")
		return nil
	} else if err != nil {
		return diag.FromErr(err)
	}

	if err := schemaFromParameterGroup(*p, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func (r *ParameterGroup) Update(ctx context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	id := d.Id()

	c, err := schemaToParameterGroup(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges("database_vendor", "database_version", "database_type") {
		return diag.FromErr(errors.New("database_vendor, database_version, database_type update is not supported"))
	}

	if err := r.svc.Update(ctx, c); err != nil {
		return diag.FromErr(err)
	}

	p, err := r.svc.Read(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := schemaFromParameterGroup(*p, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func (r *ParameterGroup) Delete(ctx context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	id := d.Id()

	if err := r.svc.Delete(ctx, id); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func schemaFromParameterGroup(p ccx.ParameterGroup, d *schema.ResourceData) error {
	d.SetId(p.ID)

	if err := d.Set("name", p.Name); err != nil {
		return fmt.Errorf("setting name: %w", err)
	}

	if err := d.Set("database_vendor", p.DatabaseVendor); err != nil {
		return fmt.Errorf("setting database_vendor: %w", err)
	}

	if err := d.Set("database_version", p.DatabaseVersion); err != nil {
		return fmt.Errorf("setting database_version: %w", err)
	}

	if err := d.Set("database_type", p.DatabaseType); err != nil {
		return fmt.Errorf("setting database_type: %w", err)
	}

	if err := setDbParameters(d, p.DbParameters); err != nil {
		return fmt.Errorf("setting parameters: %w", err)
	}

	if err := d.Set("description", p.Description); err != nil {
		return fmt.Errorf("setting description: %w", err)
	}

	return nil
}

func schemaToParameterGroup(d *schema.ResourceData) (ccx.ParameterGroup, error) {
	g := ccx.ParameterGroup{
		ID:              d.Id(),
		Name:            getString(d, "name"),
		DatabaseVendor:  getString(d, "database_vendor"),
		DatabaseVersion: getString(d, "database_version"),
		DatabaseType:    getString(d, "database_type"),
		Description:     getString(d, "description"),
	}

	g.DatabaseVendor = vendorFromAlias(g.DatabaseVendor)

	if v, err := getDbParameters(d); err == nil {
		g.DbParameters = v
	} else {
		return ccx.ParameterGroup{}, err
	}

	return g, nil
}

func getDbParameters(d *schema.ResourceData) (map[string]string, error) {
	parameters := make(map[string]string)

	if raw, ok := d.GetOk("parameters"); ok {
		for key, value := range raw.(map[string]any) {
			if strValue, ok := value.(string); ok {
				parameters[key] = strValue
			} else {
				return nil, fmt.Errorf("expected string for key %q in parameters, got %T", key, value)
			}
		}
	}

	return parameters, nil
}

func setDbParameters(d *schema.ResourceData, ls map[string]string) error {
	parameters := make(map[string]any, len(ls))
	for key, value := range ls {
		parameters[key] = value
	}

	if err := d.Set("parameters", parameters); err != nil {
		return fmt.Errorf("error setting parameters: %w", err)
	}

	return nil
}
