package api

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/severalnines/terraform-provider-ccx/ccx"
	"github.com/severalnines/terraform-provider-ccx/internal/lib"
)

type authTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

func (r authTokenRequest) URLEncode() string {
	var body = url.Values{}
	body.Set("client_id", r.ClientID)
	body.Set("client_secret", r.ClientSecret)
	body.Set("grant_type", r.GrantType)

	e := body.Encode()

	return e
}

type authTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type authorizer interface {
	Auth(_ context.Context) (string, error)
}

type tokenAuthorizer struct {
	id      string
	secret  string
	baseURL string
}

func (svc tokenAuthorizer) Auth(_ context.Context) (string, error) {
	r := authTokenRequest{
		ClientID:     svc.id,
		ClientSecret: svc.secret,
		GrantType:    "client_credentials",
	}

	b := r.URLEncode()

	req, err := http.NewRequest(http.MethodPost, svc.baseURL+"/api/auth/oauth2/token", strings.NewReader(b))
	if err != nil {
		return "", errors.Join(err, ccx.AuthenticationFailedErr)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Timeout: ccx.DefaultTimeout,
	}

	res, err := client.Do(req)
	if err != nil {
		return "", errors.Join(err, ccx.AuthenticationFailedErr)
	}

	var tokenResponse authTokenResponse
	if err := lib.DecodeJsonInto(res.Body, &tokenResponse); err != nil {
		return "", err
	}

	t := tokenResponse.TokenType + " " + tokenResponse.AccessToken

	return t, nil
}
