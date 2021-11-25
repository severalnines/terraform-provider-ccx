package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/severalnines/terraform-provider-ccx/services"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	clientTimeout = 5 * time.Second
)

type userIdResponse struct {
	UserId string `json:"user_id"`
}

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CCX_CLIENT_ID", ""),
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CCX_CLIENT_SECRET", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ccx_cluster": clusterResource(),
			"ccx_vpc":     vpcResource(),
		},
		ConfigureFunc: configureProvider,
	}
}

func getOAuth2Client(clientId, clientSecret, baseUrl string) (*http.Client, string, error) {
	conf := &clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		TokenURL:     baseUrl + "/oauth2/token",
	}
	client := &http.Client{Timeout: clientTimeout}
	// obtaining the HTTP client managing OAuth2 tokens (including refreshing)
	client = conf.Client(context.WithValue(context.Background(), oauth2.HTTPClient, client))
	// get the user ID (also checking the provided credentials work)
	req, err := http.NewRequest(http.MethodGet, baseUrl+"/oauth2userid", nil)
	if err != nil {
		return nil, "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		if err, ok := err.(*url.Error); ok {
			if err, ok := err.Unwrap().(*oauth2.RetrieveError); ok {
				return nil, "", fmt.Errorf("retrieve error: %w", err)
			}
		}
		return nil, "", err
	}
	defer func() { _ = resp.Body.Close() }()
	out := new(userIdResponse)
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return nil, "", err
	}
	return client, out.UserId, nil
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
	clientId := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	httpClient, userID, err := getOAuth2Client(clientId, clientSecret, BaseURLV1)
	if err != nil {
		return nil, err
	}
	log.Println("user ID:", userID)
	return services.NewClient(BaseURLV1, userID, httpClient), nil
}
