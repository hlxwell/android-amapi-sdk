package types

import (
	"google.golang.org/api/androidmanagement/v1"
)

// ProvisioningInfo is an alias for androidmanagement.ProvisioningInfo.
// Use androidmanagement.ProvisioningInfo directly for all provisioning operations.
type ProvisioningInfo = androidmanagement.ProvisioningInfo

// ProvisioningInfoGetRequest represents a request to get provisioning info.
type ProvisioningInfoGetRequest struct {
	// Name is the provisioning info resource name
	Name string `json:"name"`
}
