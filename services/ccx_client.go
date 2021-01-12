package services

import "net/http"

type CCXLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type Client struct {
	address    string
	userId     string
	httpClient *http.Client
	httpCookie *http.Cookie
}

func NewClient(service_address string, userID string, cookie *http.Cookie) *Client {
	return &Client{
		address:    service_address,
		userId:     userID,
		httpCookie: cookie,
		httpClient: &http.Client{},
	}

}
