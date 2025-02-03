package api

import (
	"context"
	"net/http"
)

// HttpClient is used to make requests to the ccx api
type HttpClient interface {
	// Do sends a request to the API and returns the response
	// errors returned are:
	// - ccx.RequestEncodingErr (if body encoding fails)
	// - ccx.RequestInitializationErr (if request creation fails)
	// - ccx.RequestSendingErr (if request sending fails)
	// - ccx.ResourceNotFoundErr (if API returns 404)
	// - ccx.ApiErr (if API returns 4xx or 5xx)
	Do(ctx context.Context, method, path string, body any) (*http.Response, error)

	// Get sends a GET request to the API and decodes the response into target
	// errors returned are:
	// - ccx.RequestInitializationErr (if request creation fails)
	// - ccx.RequestSendingErr (if request sending fails)
	// - ccx.ResourceNotFoundErr (if API returns 404)
	// - ccx.ApiErr (if API returns 4xx or 5xx)
	Get(ctx context.Context, path string, target any) error
}
