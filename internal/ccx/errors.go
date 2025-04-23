package ccx

import (
	"errors"
)

var (

	// ErrRequestEncoding occurs when a request body cannot be encoded.
	ErrRequestEncoding = errors.New("failed to encode request body")

	// ErrRequestInitialization occurs when failing to initialize http request object
	ErrRequestInitialization = errors.New("failed to initialize request")

	// ErrRequestSending occurs when failing to send an http request
	ErrRequestSending = errors.New("failed to send request")

	// ErrResponseReadFailed occurs when failing to read the response from the server
	ErrResponseReadFailed = errors.New("failed to read response")

	// ErrResponseDecoding occurs when failing to decode the response body
	ErrResponseDecoding = errors.New("failed to decode response")

	// ErrResourceNotFound occurs when trying to get a resource that does not exist
	ErrResourceNotFound = errors.New("resource not found")

	// ErrAllocatingAZs occurs when failing to allocate availability zones
	ErrAllocatingAZs = errors.New("allocating availability zones")

	// ErrApi occurs when an error is returned from the api, typically when status code is 4xx or 5xx
	ErrApi = errors.New("api error")

	// ErrUpdateNotSupported occurs when trying to update a resource which might be destructive if attempted
	ErrUpdateNotSupported = errors.New("updates for this resource are not supported")

	// ErrCreateFailedRead indicates failure to read a newly created resource. The resource may exist, but we have the id and terraform can possibly read it on next apply.
	ErrCreateFailedRead = errors.New("reading newly created resource failed")

	// ErrApplyDbParametersFailed indicates failure to apply database parameter group
	ErrApplyDbParametersFailed = errors.New("failed to apply database parameter group")

	// ErrFirewallRules indicates failure to configure firewall rules
	ErrFirewallRules = errors.New("failed to configure firewall rules")

	// ErrMaintenanceSettings indicates failure to configure maintenance settings
	ErrMaintenanceSettings = errors.New("failed to configure maintenance settings")
)
