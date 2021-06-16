package provider

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/severalnines/terraform-provider-ccx/services"
)

func vpcResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vpc_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_cloud_provider": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_cloud_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_ipv4_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Create: createVpc,
		Read:   readVpc,
		Update: createVpc,
		Delete: deleteVpc,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
func createVpc(d *schema.ResourceData, m interface{}) error {
	//General Settings
	log.Printf(d.Get("vpc_name").(string))
	vpcName := d.Get("vpc_name").(string)
	vpcCloudProviderName := d.Get("vpc_cloud_provider").(string)
	vpcCloudRegion := d.Get("vpc_cloud_region").(string)
	vpcCidr := d.Get("vpc_ipv4_cidr").(string)
	client := m.(*services.Client)
	serviceResponse, err := client.CreateVpc(vpcName, vpcCloudProviderName, vpcCloudRegion, vpcCidr)
	if err != nil {
		return err
	}
	d.SetId(serviceResponse.VpcUUID)
	return nil
}
func readVpc(d *schema.ResourceData, m interface{}) error {
	client := m.(*services.Client)
	client.GetVPCbyUUID(d.Id())
	d.SetId(d.Id())
	return nil
}
func deleteVpc(d *schema.ResourceData, m interface{}) error {
	client := m.(*services.Client)
	client.DeleteVPCbyUUID(d.Id())
	d.SetId(d.Id())
	return nil
}
