package provider

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/severalnines/terraform-provider-ccx/services"
)

func resourceItem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the resource, also acts as it's unique ID",
				ForceNew:     true,
				ValidateFunc: validateName,
			},
			"cluster_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A description of an item",
			},
			"cloud_provider": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional list of tags, represented as a key, value pair",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional list of tags, represented as a key, value pair",
			},
			"db_vendor": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "An optional list of tags, represented as a key, value pair",
			},
			"instance_size": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "An optional list of tags, represented as a key, value pair",
			},
			"instance_iops": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "An optional list of tags, represented as a key, value pair",
			},
			"db_username": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "An optional list of tags, represented as a key, value pair",
			},
			"db_password": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "An optional list of tags, represented as a key, value pair",
			},
			"db_host": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "An optional list of tags, represented as a key, value pair",
			},
		},
		Create: resourceCreateItem,
		Read:   resourceReadItem,
		Update: resourceCreateItem,
		Delete: resourceCreateItem,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
func validateName(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("Expected name to be string"))
		return warns, errs
	}
	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(value)) {
		errs = append(errs, fmt.Errorf("name cannot contain whitespace. Got %s", value))
		return warns, errs
	}
	return warns, errs
}
func resourceCreateItem(d *schema.ResourceData, m interface{}) error {

	address := d.Get("address").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	clusterName := d.Get("cluster_name").(string)
	clusterType := d.Get("cluster_type").(string)
	clusterProvider := d.Get("cluster_provider").(string)
	region := d.Get("region").(string)
	dbVendor := d.Get("db_vendor").(string)
	instanceSize := d.Get("instance_size").(string)
	instanceIops := d.Get("instance_iops").(int)
	dbUsername := d.Get("db_username").(string)
	dbPassword := d.Get("db_password").(string)
	dbHost := d.Get("db_host").(string)
	userID, cookie := services.GetUserId(address, username, password)
	err := services.CreateCluster(userID, clusterName, clusterType,
		clusterProvider, region, dbVendor, instanceSize, instanceIops, dbUsername, dbPassword,
		dbHost, cookie)
	pendingClusterID := services.GetClusters(userID, cookie)
	for i := range pendingClusterID {
		if pendingClusterID[i].ClusterStatus == "CREATING_CLUSTER" {
			d.SetId(pendingClusterID[i].UUID)
		}

	}

	if err != nil {
		return err
	}
	return nil
}
func resourceReadItem(d *schema.ResourceData, m interface{}) error {
	address := d.Get("address").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	clusterID := d.Id()
	_, cookie := services.GetUserId(address, username, password)
	clusterInfo := services.GetClusterByID(clusterID, cookie)
	d.SetId(clusterInfo.UUID)
	return nil
}
