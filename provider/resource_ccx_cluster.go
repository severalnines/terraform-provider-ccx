package provider

import (
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/severalnines/terraform-provider-ccx/services"
)

func clusterResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the resource, also acts as it's unique ID",
				ForceNew:     true,
				ValidateFunc: validateName,
			},
			"cluster_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The size of the cluster ( int64 ). 1 or 3 nodes.",
			},
			"db_vendor": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional list of tags, represented as a key, value pair",
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "An optional list of tags, represented as a key, value pair",
				Elem:        &schema.Schema{Type: schema.TypeString},
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

			"instance_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional list of tags, represented as a key, value pair",
			},
			"volume_iops": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional list of tags, represented as a key, value pair",
			},
			"volume_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "An optional list of tags, represented as a key, value pair",
			},
			"volume_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional list of tags, represented as a key, value pair",
			},
			"network_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional list of tags, represented as a key, value pair",
			},
			"network_ha_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "An optional list of tags, represented as a key, value pair",
			},
			"network_vpc_uuid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional list of tags, represented as a key, value pair",
			},
			"network_az": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "An optional list of tags, represented as a key, value pair",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Create: resourceCreateItem,
		Read:   resourceReadItem,
		Update: resourceCreateItem,
		Delete: resourceDeleteItem,
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
		errs = append(errs, fmt.Errorf("expected name to be string"))
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
	//General Settings
	clusterName := d.Get("cluster_name").(string)
	clusterSize := d.Get("cluster_size").(int)
	dbVendor := d.Get("db_vendor").(string)
	objectTags := []string{}
	for _, tag := range d.Get("tags").([]interface{}) {
		objectTags = append(objectTags, tag.(string))
	}
	//Cloud Settings
	cloudProvider := d.Get("cloud_provider").(string)
	cloudRegion := d.Get("region").(string)
	//Instance Settings
	instanceSize := d.Get("instance_size").(string)
	volumeType := d.Get("volume_type").(string)
	volumeIops := d.Get("volume_iops").(string)
	volumeSize := d.Get("volume_size").(int)
	//Network
	networkType := d.Get("network_type").(string)
	vpcUUID := d.Get("network_vpc_uuid").(string)
	//
	client := m.(*services.Client)
	serviceResponse, err := client.CreateCluster(
		clusterName, clusterSize, dbVendor, objectTags, cloudRegion, cloudProvider, instanceSize,
		volumeType, volumeSize, volumeIops, networkType, vpcUUID)

	if err != nil {
		log.Println(err)
		return err
	}
	d.SetId(serviceResponse.ClusterUUID)
	return nil
}
func resourceReadItem(d *schema.ResourceData, m interface{}) error {
	client := m.(*services.Client)
	client.GetClusterByID(d.Id())
	d.SetId(d.Id())
	return nil
}
func resourceDeleteItem(d *schema.ResourceData, m interface{}) error {
	client := m.(*services.Client)
	client.DeleteCluster(d.Id())
	d.SetId(d.Id())
	return nil
}
