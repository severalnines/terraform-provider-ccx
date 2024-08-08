package api

import (
	"context"
	"net/http"
)

// HttpClient is used to make requests to the ccx api
type HttpClient interface {
	Do(ctx context.Context, method, path string, body any) (*http.Response, error)
	Get(ctx context.Context, path string, target any) error
}
