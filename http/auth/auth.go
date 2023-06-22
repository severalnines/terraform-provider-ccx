package auth

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/severalnines/terraform-provider-ccx/ccx"
	chttp "github.com/severalnines/terraform-provider-ccx/http"
)

type GetTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

func (r GetTokenRequest) URLEncode() string {
	var body = url.Values{}
	body.Set("client_id", r.ClientID)
	body.Set("client_secret", r.ClientSecret)
	body.Set("grant_type", r.GrantType)

	e := body.Encode()

	return e
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type Authorizer struct {
	id     string
	secret string
	conn   *chttp.ConnectionParameters
}

func New(clientID, clientSecret string, opts ...chttp.ParameterOption) *Authorizer {
	p := chttp.Parameters(opts...)

	return &Authorizer{
		id:     clientID,
		secret: clientSecret,
		conn:   p,
	}
}

func (a *Authorizer) Auth(_ context.Context) (string, error) {
	r := GetTokenRequest{
		ClientID:     a.id,
		ClientSecret: a.secret,
		GrantType:    "client_credentials",
	}

	b := r.URLEncode()

	req, err := http.NewRequest(http.MethodPost, a.conn.BaseURL+"/api/auth/oauth2/token", strings.NewReader(b))
	if err != nil {
		return "", errors.Join(err, ccx.AuthenticationFailedErr)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Timeout: a.conn.Timeout,
	}

	res, err := client.Do(req)
	if err != nil {
		return "", errors.Join(err, ccx.AuthenticationFailedErr)
	}

	var tokenResponse TokenResponse
	if err := chttp.DecodeJsonInto(res.Body, &tokenResponse); err != nil {
		return "", err
	}

	t := tokenResponse.TokenType + " " + tokenResponse.AccessToken

	return t, nil
}
