package types

import "google.golang.org/api/androidmanagement/v1"

// Type aliases for compatibility with callers that expect types from this package.
type (
	Enterprise       = androidmanagement.Enterprise
	Policy           = androidmanagement.Policy
	Device           = androidmanagement.Device
	EnrollmentToken  = androidmanagement.EnrollmentToken
	MigrationToken   = androidmanagement.MigrationToken
	WebApp           = androidmanagement.WebApp
	WebToken         = androidmanagement.WebToken
	ProvisioningInfo = androidmanagement.ProvisioningInfo
)

// ListOptions captures common pagination parameters used by higher-level callers.
type ListOptions struct {
	PageSize  int
	PageToken string
	Filter    string
}

