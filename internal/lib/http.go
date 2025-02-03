package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
	"golang.org/x/oauth2/clientcredentials"
)

// ErrorResponse represents generic error responses from ccx api
type ErrorResponse struct {
	Code    json.Number `json:"code"`
	Err     string      `json:"err"`
	ErrLong string      `json:"error"`
	status  int
}

func (r ErrorResponse) Error() string {
	s := r.Err

	if s == "" && r.ErrLong != "" {
		s = r.ErrLong
	}

	if s == "" {
		s = "an error occurred"
	}

	s += " ("

	if r.Code != "" {
		s += "code: " + r.Code.String() + ", "
	}

	s += "response: " + strconv.Itoa(r.status)

	if t := http.StatusText(r.status); t != "" {
		s += " - " + t
	}

	s += ")"

	return s
}

// ErrorFromResponse decodes the body of type ErrorResponse and returns an error
func ErrorFromResponse(rs *http.Response) error {
	var e ErrorResponse

	e.status = rs.StatusCode

	if rs.Body == nil {
		return e
	}

	defer Closed(rs.Body)
	b, err := io.ReadAll(rs.Body)
	if err != nil {
		e.ErrLong = fmt.Sprintf("could not read reason: %s", err.Error())
		return e
	}

	err = json.Unmarshal(b, &e)
	if err != nil {
		e.ErrLong = fmt.Sprintf("could not decode reason: %s", err.Error())
		return e
	}

	return e
}

// DecodeJsonInto is a helper to decode JSON body into a target type
func DecodeJsonInto(body io.ReadCloser, target any) error {
	defer Closed(body)

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

// Closed will close an io.Closer object silently
// useful for deferred closing
func Closed(c io.Closer) {
	if c != nil {
		_ = c.Close()
	}
}

func NewHttpClient(baseURL, clientID, clientSecret string) *HttpClient {
	creds := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     baseURL + "/api/auth/oauth2/token",
	}

	// TF context is canceled to soon on import
	cli := creds.Client(context.Background())

	cli.Timeout = ccx.DefaultTimeout
	cli.Transport = &LoggingRoundTripper{
		Proxied: cli.Transport,
	}

	return &HttpClient{
		baseURL: baseURL,
		cli:     cli,
	}
}

func NewTestHttpClient(baseURL string) *HttpClient {
	cli := http.DefaultClient

	return &HttpClient{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		cli:     cli,
	}
}

type HttpClient struct {
	baseURL string
	cli     *http.Client
}

// Do sends a request to the ccx api
// errors returned are:
// - ccx.RequestEncodingErr (if body encoding fails)
// - ccx.RequestInitializationErr (if request creation fails)
// - ccx.RequestSendingErr (if request sending fails)
// - ccx.ResourceNotFoundErr (if API returns 404)
// - ccx.ApiErr (if API returns 4xx or 5xx)
func (h *HttpClient) Do(ctx context.Context, method, path string, body any) (*http.Response, error) {
	var b bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&b).Encode(body); err != nil {
			return nil, errors.Join(ccx.RequestEncodingErr, err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, h.baseURL+path, &b)
	if err != nil {
		return nil, errors.Join(ccx.RequestInitializationErr, err)
	}

	rs, err := h.cli.Do(req)

	if err != nil {
		return nil, errors.Join(ccx.RequestSendingErr, err)
	} else if rs.StatusCode == http.StatusNotFound {
		return nil, ccx.ResourceNotFoundErr
	} else if rs.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("%w: %w", ccx.ApiErr, ErrorFromResponse(rs))
	}

	return rs, nil
}

// Get sends a GET request to the ccx api
// errors returned are:
// - ccx.RequestInitializationErr (if request creation fails)
// - ccx.RequestSendingErr (if request sending fails)
// - ccx.ResourceNotFoundErr (if API returns 404)
// - ccx.ApiErr (if API returns 4xx or 5xx)
func (h *HttpClient) Get(ctx context.Context, path string, target any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.baseURL+path, nil)
	if err != nil {
		return errors.Join(ccx.RequestInitializationErr, err)
	}

	rs, err := h.cli.Do(req)

	if err != nil {
		return errors.Join(ccx.RequestSendingErr, err)
	} else if rs.StatusCode == http.StatusNotFound {
		return ccx.ResourceNotFoundErr
	} else if rs.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("%w: %w", ccx.ApiErr, ErrorFromResponse(rs))
	}

	defer func() {
		if rs.Body != nil {
			_ = rs.Body.Close()
		}
	}()

	b, err := io.ReadAll(rs.Body)
	if err != nil {
		return errors.Join(ccx.ResponseReadFailedErr, err)
	}

	if err := json.Unmarshal(b, target); err != nil {
		return errors.Join(ccx.ResponseDecodingErr, err)
	}

	return nil
}
