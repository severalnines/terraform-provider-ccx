package services

import "net/http"

type Client struct {
	address    string
	userId     string
	httpClient *http.Client
}

func NewClient(service_address string, userID string, client *http.Client) *Client {
	return &Client{
		address:    service_address,
		userId:     userID,
		httpClient: client,
	}
}
