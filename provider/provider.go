package provider

import (
	"log"
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/severalnines/terraform-provider-ccx/services"
)

type CCXLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CCX_USERNAME", ""),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CCX_PASSWORD", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ccx_cluster": clusterResource(),
			"ccx_vpc":     vpcResource(),
		},
		ConfigureFunc: configureProvider,
	}
}
func configureProvider(d *schema.ResourceData) (interface{}, error) {
	var BaseURLV1 string
	switch os.Getenv("ENVIRONMENT") {
	case "dev":
		BaseURLV1 = services.AuthServiceUrlDev
	case "test":
		BaseURLV1 = services.AuthServiceUrlTest
	case "prod":
		BaseURLV1 = services.AuthServiceUrlProd
	default:
		BaseURLV1 = services.AuthServiceUrlProd
	}
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	userID, httpCookie, err := services.GetUserId(BaseURLV1, username, password)
	if err != nil {
		return nil, err
	}
	log.Println(userID, httpCookie)
	return services.NewClient(BaseURLV1, userID, httpCookie), nil
}
