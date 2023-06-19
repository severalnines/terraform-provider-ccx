package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	ccxprov "github.com/severalnines/terraform-provider-ccx"
	cio "github.com/severalnines/terraform-provider-ccx/io"
)

// DefaultBaseURL to access API services
const DefaultBaseURL = "https://app.mydbservice.net"

const (
	DefaultTimeout = time.Second * 5
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
func BaseURL(url string) ParameterOptionFn {
	return func(p *ConnectionParameters) {
		p.BaseURL = url
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
	Error string `json:"err"`
}

// ErrorFromErrorResponse decodes the body of type ErrorResponse and returns an error
func ErrorFromErrorResponse(body io.ReadCloser) error {
	defer cio.Close(body)
	b, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("%w: reason could not be read", ccxprov.InvalidRequestErr)
	}

	var r ErrorResponse
	err = json.Unmarshal(b, &r)
	if err != nil {
		return fmt.Errorf("%w: reason could not be decoded", ccxprov.InvalidRequestErr)
	}

	return fmt.Errorf("%w: %s", ccxprov.InvalidRequestErr, r.Error)
}

// DecodeJsonInto is a helper to decode JSON body into a target type
func DecodeJsonInto(body io.ReadCloser, target any) error {
	defer cio.Close(body)

	raw, err := io.ReadAll(body)
	if err != nil {
		return errors.Join(ccxprov.ResponseReadFailedErr, err)
	}

	err = json.Unmarshal(raw, target)
	if err != nil {
		return errors.Join(ccxprov.ResponseDecodingErr, err)
	}

	return nil
}

// DumpRequest will dump a request in a http request file. Useful for debugging.
func DumpRequest(method, url, token, body string) {
	var authorization string
	if token != "" {
		authorization = "Authorization " + token + "\n"
	}
	data := fmt.Sprintf("%s %s\n%s\n%s", method, url, authorization, body)
	filename := "req_" + time.Now().Format("2006_01_02-15_04_05_999999999Z07_00") + ".http"
	_ = os.WriteFile(filename, []byte(data), 0644)
}
