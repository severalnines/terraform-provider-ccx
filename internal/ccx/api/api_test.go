package api

import (
	"io"
	"net/http"
	"strings"
)

func fakeHttpResponse(code int, body string) *http.Response {
	var b io.ReadCloser

	if body != "" {
		b = io.NopCloser(strings.NewReader(body))
	}

	return &http.Response{
		StatusCode: code,
		Body:       b,
	}
}
