package ccx

import (
	"errors"
)

var (

	// RequestEncodingErr occurs when a request body cannot be encoded.
	RequestEncodingErr = errors.New("failed to encode request body")

	// RequestInitializationErr occurs when failing to initialize http request object
	RequestInitializationErr = errors.New("failed to initialize request")

	// RequestSendingErr occurs when failing to send an http request
	RequestSendingErr = errors.New("failed to send request")

	// ResponseStatusFailedErr occurs when the server does not respond with an expected status
	ResponseStatusFailedErr = errors.New("response status failed")

	// ResponseReadFailedErr occurs when failing to read the response from the server
	ResponseReadFailedErr = errors.New("failed to read response")

	// ResponseDecodingErr occurs when failing to decode the response body
	ResponseDecodingErr = errors.New("failed to decode response")

	// ResourceNotFoundErr occurs when trying to get a resource that does not exist
	ResourceNotFoundErr = errors.New("resource not found")

	// ResourcesLoadFailedErr occurs when trying to load resources fails
	ResourcesLoadFailedErr = errors.New("failed to load resources")

	// UpdateNotSupportedErr occurs when trying to update a resource which might be destructive if attempted
	UpdateNotSupportedErr = errors.New("updates for this resource are not supported")

	// AuthenticationFailedErr indicates failure to authenticate with the api server
	AuthenticationFailedErr = errors.New("authentication failed")

	// CreateFailedErr indicates failure to create a resource
	CreateFailedErr = errors.New("failed to create resource")

	// CreateFailedReadErr indicates failure to read a newly created resource. The resource may exist, but we have the id and terraform can possibly read it on next apply.
	CreateFailedReadErr = errors.New("reading newly created resource failed")

	// ParametersErr indicates failure to configure parameters
	ParametersErr = errors.New("failed to configure parameters")

	// FirewallRulesErr indicates failure to configure firewall rules
	FirewallRulesErr = errors.New("failed to configure firewall rules")
)
