package provider

import (
	"fmt"
	"log"

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
			"auth_service_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CCX_AUTH_SERVICE_URL", ""),
			},
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
			"ccx_cluster": resourceItem(),
		},
		ConfigureFunc: configureProvider,
	}
}
func configureProvider(d *schema.ResourceData) (interface{}, error) {
	log.Println("Entered configure provider sequence")
	address := d.Get("auth_service_url").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	userID, httpCookie, _ := services.GetUserId(address, username, password)
	fmt.Println(userID, httpCookie)
	return services.NewClient(address, userID, httpCookie), nil
}
