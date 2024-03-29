package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/severalnines/terraform-provider-ccx/ccx"
	cio "github.com/severalnines/terraform-provider-ccx/io"
)

// DefaultBaseURL to access API services
const DefaultBaseURL = "https://app.mydbservice.net"

const (
	DefaultTimeout = time.Second * 30
)

type ConnectionParameters struct {
	BaseURL   string
	Timeout   time.Duration
	Transport http.RoundTripper
}

// Parameters is a factory for ConnectionParameters with sane defaults and optional overrides
func Parameters(opts ...ParameterOption) *ConnectionParameters {
	p := &ConnectionParameters{
		BaseURL:   DefaultBaseURL,
		Timeout:   DefaultTimeout,
		Transport: http.DefaultTransport,
	}

	for i := range opts {
		opts[i].Set(p)
	}

	return p
}

type ParameterOption interface {
	Set(params *ConnectionParameters)
}

type ParameterOptionFn func(p *ConnectionParameters)

func (f ParameterOptionFn) Set(p *ConnectionParameters) {
	f(p)
}

// BaseURL to specify a different url for the provisioning services
// If an empty url is provided, it will default to DefaultBaseURL
func BaseURL(url string) ParameterOptionFn {
	if url == "" {
		url = DefaultBaseURL
	}

	return func(p *ConnectionParameters) {
		p.BaseURL = strings.TrimSuffix(url, "/")
	}
}

// Timeout duration for requests
func Timeout(duration time.Duration) ParameterOptionFn {
	return func(p *ConnectionParameters) {
		p.Timeout = duration
	}
}

// RoundTripper sets a custom requester for a
func RoundTripper(t http.RoundTripper) ParameterOptionFn {
	return func(p *ConnectionParameters) {
		p.Transport = t
	}
}

// Authorizer retrieves authorization tokens
type Authorizer interface {
	Auth(ctx context.Context) (string, error)
}

// ErrorResponse represents generic error responses from ccx api
type ErrorResponse struct {
	Code    json.Number `json:"code"`
	Err     string      `json:"err"`
	ErrLong string      `json:"error"`
}

func (r ErrorResponse) Error() string {
	prefix := r.Code.String()
	if prefix != "" {
		prefix += ": "
	}
	if r.Err != "" {
		return prefix + r.Err
	}
	return prefix + r.ErrLong
}

// ErrorFromErrorResponse decodes the body of type ErrorResponse and returns an error
func ErrorFromErrorResponse(body io.ReadCloser) error {
	defer cio.Close(body)
	b, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("%w: reason could not be read", err)
	}

	var r ErrorResponse
	err = json.Unmarshal(b, &r)
	if err != nil {
		return fmt.Errorf("%w: reason could not be decoded", err)
	}

	return r
}

// DecodeJsonInto is a helper to decode JSON body into a target type
func DecodeJsonInto(body io.ReadCloser, target any) error {
	defer cio.Close(body)

	raw, err := io.ReadAll(body)
	if err != nil {
		return errors.Join(ccx.ResponseReadFailedErr, err)
	}

	err = json.Unmarshal(raw, target)
	if err != nil {
		return errors.Join(ccx.ResponseDecodingErr, err)
	}

	return nil
}
