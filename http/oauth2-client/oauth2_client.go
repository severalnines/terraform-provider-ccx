package oauth2_client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	ccxprov "github.com/severalnines/terraform-provider-ccx"
	chttp "github.com/severalnines/terraform-provider-ccx/http"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type GetUserIDResponse struct {
	UserId string `json:"user_id"`
}

// Client returns a http.Client with oauth2 authentication
func Client(ctx context.Context, clientID, clientSecret string, opts ...chttp.ParameterOption) (*http.Client, error) {
	p := chttp.Parameters(opts...)

	conf := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     p.BaseURL + "/api/auth/oauth2/token",
	}

	client := &http.Client{Transport: p.Transport}

	// obtaining the HTTP client managing OAuth2 tokens (including refreshing)
	ctx = context.WithValue(ctx, oauth2.HTTPClient, client)

	token := conf.TokenSource(ctx)
	client = oauth2.NewClient(ctx, token)
	client.Timeout = p.Timeout

	// get the user ID (also checking the provided credentials work)
	req, err := http.NewRequest(http.MethodGet, p.BaseURL+"/api/auth/oauth2userid", nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		// if err, ok := err.(*url.Error); ok {
		// 	if err, ok := err.Unwrap().(*oauth2.RetrieveError); ok {
		// 		return nil, errors.Join(err, ccxprov.AuthenticationFailedErr)
		// 	}
		// }
		return nil, errors.Join(err, ccxprov.AuthenticationFailedErr)
	}

	defer func() { _ = resp.Body.Close() }()
	var userIDResp GetUserIDResponse
	if err := json.NewDecoder(resp.Body).Decode(&userIDResp); err != nil {
		return nil, errors.Join(err, ccxprov.AuthenticationFailedErr)
	}

	return client, nil
}
