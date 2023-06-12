package http

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	cio "github.com/severalnines/terraform-provider-ccx/io"
	testhttp "github.com/severalnines/terraform-provider-ccx/tests/http"
)

type Result struct {
	Response *http.Response
	Error    error
}

// RoundTripper helps to mock server responses while recording http.Request calls
type RoundTripper struct {
	Results []Result
	Calls   []testhttp.Request
	mut     sync.Mutex
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	r.mut.Lock()
	defer r.mut.Unlock()

	nc := len(r.Calls)
	nr := len(r.Results)
	if nr-1 < nc {
		return nil, fmt.Errorf("ran out of responses for call [%d], have only [%d] responses", nc+1, nr)
	}

	var body string
	if req.Body != nil {
		defer cio.Close(req.Body)
		if b, err := io.ReadAll(req.Body); err == nil {
			body = string(b)
		} else {
			return nil, fmt.Errorf("failed to read body: %w", err)
		}
	}

	call := testhttp.Request{
		Host:   req.URL.Host,
		Path:   req.URL.Path,
		Query:  req.URL.RawQuery,
		Method: req.Method,
		Header: req.Header,
		Body:   body,
	}

	r.Calls = append(r.Calls, call)

	rs := r.Results[nc]
	return rs.Response, rs.Error
}

func Status(code int) ResponseOptionFn {
	return func(r *http.Response) {
		r.StatusCode = code
		r.Status = fmt.Sprintf("%d %s", code, http.StatusText(code))
	}
}

func Body(s string) ResponseOptionFn {
	return func(r *http.Response) {
		r.Body = io.NopCloser(strings.NewReader(s))
	}
}

func Response(opts ...ResponseOption) *http.Response {
	r := &http.Response{
		StatusCode:    http.StatusOK,
		Status:        "200 OK",
		ContentLength: 0,
	}

	for i := range opts {
		opts[i].Apply(r)
	}

	return r
}

type ResponseOption interface {
	Apply(r *http.Response)
}

type ResponseOptionFn func(r *http.Response)

func (o ResponseOptionFn) Apply(r *http.Response) {
	o(r)
}
