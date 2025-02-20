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

	// AllocatingAZsErr occurs when failing to allocate availability zones
	AllocatingAZsErr = errors.New("allocating availability zones")

	// ApiErr occurs when an error is returned from the api, typically when status code is 4xx or 5xx
	ApiErr = errors.New("api error")

	// UpdateNotSupportedErr occurs when trying to update a resource which might be destructive if attempted
	UpdateNotSupportedErr = errors.New("updates for this resource are not supported")

	// AuthenticationFailedErr indicates failure to authenticate with the api server
	AuthenticationFailedErr = errors.New("authentication failed")

	// CreateFailedReadErr indicates failure to read a newly created resource. The resource may exist, but we have the id and terraform can possibly read it on next apply.
	CreateFailedReadErr = errors.New("reading newly created resource failed")

	// ApplyDbParametersFailedErr indicates failure to apply database parameter group
	ApplyDbParametersFailedErr = errors.New("failed to apply database parameter group")

	// FirewallRulesErr indicates failure to configure firewall rules
	FirewallRulesErr = errors.New("failed to configure firewall rules")

	// MaintenanceSettingsErr indicates failure to configure maintenance settings
	MaintenanceSettingsErr = errors.New("failed to configure maintenance settings")
)
