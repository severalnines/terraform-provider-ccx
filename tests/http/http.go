package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type Request struct {
	Host   string
	Path   string
	Query  string
	Method string
	Header http.Header
	Body   string
}

type Response struct {
	Status int
	Header http.Header
	Body   string
	Error  error
}

// Test case with assertions
type Test struct {
	// Request details we will be sending
	Request Request

	// Want stuffs
	Want Response
}

// Assert all as per Want
func (h Test) Assert(t *testing.T, handler http.HandlerFunc) bool {
	var (
		b     io.Reader
		valid = true
	)

	if len(h.Request.Body) != 0 {
		b = strings.NewReader(h.Request.Body)
	}

	r, err := http.NewRequest(h.Request.Method, h.Request.Path, b)
	if err != nil {
		t.Errorf("failed to init request: %s", err.Error())
		return false
	}

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	var (
		result  = w.Result()
		status  = result.StatusCode
		headers = w.Header()

		lg = len(headers)
		lw = len(h.Want.Header)
		wb = h.Want.Body
		gb = w.Body.String()
	)

	if status != h.Want.Status {
		t.Errorf("status not as wanted: want = %d, got = %d", h.Want.Status, status)
		valid = false
	}

	if wb != gb {
		var replacer = strings.NewReplacer(
			"\r", "[R]", "\n", "[N]",
		)

		wb = replacer.Replace(wb)
		gb = replacer.Replace(gb)

		t.Errorf("body not as wanted.\nwant = %s\ngot  = %s", wb, gb)
		valid = false
	}

	if lw != lg {
		t.Errorf("headers count: want = %d, got = %d", lw, lg)
		valid = false
		return false
	}

	for k, wv := range h.Want.Header {
		var gv = headers[k]

		if len(wv) != len(gv) {
			t.Errorf("header [%s] length: want = %d, got = %d", k, len(wv), len(gv))
			valid = false
			continue
		}

		for i := range wv {
			if wv[i] != gv[i] {
				t.Errorf("header[%s][%d]: want = %s, got = %s", k, i, wv[i], gv[i])
				valid = false
			}
		}
	}

	return valid
}
