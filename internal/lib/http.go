package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/severalnines/terraform-provider-ccx/internal/ccx"
	"golang.org/x/oauth2/clientcredentials"
)

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
	defer Closed(body)
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

func NewHttpClient(module, baseURL, clientID, clientSecret, logpath string) *HttpClient {
	if logpath != "" {
		if err := os.MkdirAll(logpath, 0755); err != nil {
			panic(fmt.Sprintf("failed to create log directory [%s]: %s", logpath, err))
		}
	}

	creds := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     baseURL + "/api/auth/oauth2/token",
	}

	return &HttpClient{
		module:  module,
		baseURL: baseURL,
		logpath: logpath,
		creds:   creds,
	}
}

type HttpClient struct {
	module  string
	baseURL string
	logpath string

	creds *clientcredentials.Config
}

func (h *HttpClient) Do(ctx context.Context, method, path string, body any) (*http.Response, error) {
	var b bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&b).Encode(body); err != nil {
			return nil, errors.Join(ccx.RequestEncodingErr, err)
		}
	}

	req, err := http.NewRequest(method, h.baseURL+path, &b)
	if err != nil {
		return nil, errors.Join(ccx.RequestInitializationErr, err)
	}

	cli := h.cli(ctx)

	return cli.Do(req)
}

func (h *HttpClient) Get(ctx context.Context, path string, target any) error {
	req, err := http.NewRequest(http.MethodGet, h.baseURL+path, nil)
	if err != nil {
		return errors.Join(ccx.RequestInitializationErr, err)
	}

	client := h.cli(ctx)

	res, err := client.Do(req)
	if err != nil {
		return errors.Join(ccx.RequestSendingErr, err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: status = %d", ErrorFromErrorResponse(res.Body), res.StatusCode)
	}

	defer func() {
		if res.Body != nil {
			_ = res.Body.Close()
		}
	}()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.Join(ccx.ResponseReadFailedErr, err)
	}

	if err := json.Unmarshal(b, target); err != nil {
		return errors.Join(ccx.ResponseDecodingErr, err)
	}

	return nil
}

func (h *HttpClient) cli(ctx context.Context) *http.Client {
	cli := h.creds.Client(ctx)
	cli.Timeout = ccx.DefaultTimeout

	if h.logpath != "" {
		cli.Transport = &LoggingRoundTripper{
			LogPath: h.logpath,
			Module:  h.module,
			Proxied: cli.Transport,
		}
	}

	return cli
}
