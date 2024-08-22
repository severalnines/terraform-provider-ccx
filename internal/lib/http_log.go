package lib

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// LoggingRoundTripper is a custom RoundTripper that logs request and response details
type LoggingRoundTripper struct {
	Proxied http.RoundTripper
}

// RoundTrip executes a single HTTP transaction and logs the details
func (l *LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Log the request
	requestDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}

	res, err := l.Proxied.RoundTrip(req)

	if err != nil {
		return nil, err
	}

	// Log the response
	responseDump, err := httputil.DumpResponse(res, true)
	if err != nil {
		return nil, err
	}

	tflog.Debug(req.Context(), fmt.Sprintf("ccx api request %s", req.URL.Path), map[string]any{
		"request":  string(requestDump),
		"response": string(responseDump),
	})

	return res, nil
}
