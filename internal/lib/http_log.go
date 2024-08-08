package lib

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

// LoggingRoundTripper is a custom RoundTripper that logs request and response details
type LoggingRoundTripper struct {
	LogPath string
	Module  string
	Proxied http.RoundTripper
}

// RoundTrip executes a single HTTP transaction and logs the details
func (l *LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Log the request
	requestDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}

	// Measure the time taken for the request
	res, err := l.Proxied.RoundTrip(req)

	if err != nil {
		return nil, err
	}

	// Log the response
	responseDump, err := httputil.DumpResponse(res, true)
	if err != nil {
		return nil, err
	}

	prefix := strings.ReplaceAll(time.Now().Format("2006-01-02-15-04-05.999999999"), ".", "-")

	reqfile := fmt.Sprintf("%s/%s-%s-request.http", l.LogPath, prefix, l.Module)
	if err := os.WriteFile(reqfile, requestDump, 0644); err != nil {
		panic(fmt.Sprintf("failed to write file [%s]: %s", reqfile, err))
	}

	resfile := fmt.Sprintf("%s/%s-%s-response.log", l.LogPath, prefix, l.Module)
	if err := os.WriteFile(resfile, responseDump, 0644); err != nil {
		panic(fmt.Sprintf("failed to write file [%s]: %s", resfile, err))
	}

	return res, nil
}
