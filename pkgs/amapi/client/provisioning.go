package client

import (
	"google.golang.org/api/androidmanagement/v1"

	"amapi-pkg/pkgs/amapi/types"
)

// ProvisioningService provides provisioning info management methods.
type ProvisioningService struct {
	client *Client
}

// ProvisioningInfo returns the provisioning service.
func (c *Client) ProvisioningInfo() *ProvisioningService {
	return &ProvisioningService{client: c}
}

// Get retrieves provisioning information by its resource name.
func (ps *ProvisioningService) Get(provisioningInfoName string) (*androidmanagement.ProvisioningInfo, error) {
	if provisioningInfoName == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "provisioning info name is required")
	}

	var result *androidmanagement.ProvisioningInfo
	var err error

	err = ps.client.executeAPICall(func() error {
		result, err = ps.client.service.ProvisioningInfo.Get(provisioningInfoName).Context(ps.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, ps.client.wrapAPIError(err, "get provisioning info")
	}

	return result, nil
}

// GetByID retrieves provisioning information by ID.
func (ps *ProvisioningService) GetByID(provisioningInfoID string) (*androidmanagement.ProvisioningInfo, error) {
	if provisioningInfoID == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "provisioning info ID is required")
	}

	provisioningInfoName := buildProvisioningInfoName(provisioningInfoID)
	return ps.Get(provisioningInfoName)
}

// GetByDeviceID retrieves provisioning information by device ID.
func (ps *ProvisioningService) GetByDeviceID(deviceID string) (*androidmanagement.ProvisioningInfo, error) {
	if deviceID == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "device ID is required")
	}

	// Build provisioning info name from device ID
	// The exact format depends on the API specification
	provisioningInfoName := "provisioningInfo/" + deviceID
	return ps.Get(provisioningInfoName)
}

// GetByEnterpriseID retrieves provisioning information by enterprise ID.
func (ps *ProvisioningService) GetByEnterpriseID(enterpriseID string) (*androidmanagement.ProvisioningInfo, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	// Build provisioning info name from enterprise ID
	// The exact format depends on the API specification
	provisioningInfoName := "provisioningInfo/" + enterpriseID
	return ps.Get(provisioningInfoName)
}

// Helper function to build provisioning info name
func buildProvisioningInfoName(provisioningInfoID string) string {
	return "provisioningInfo/" + provisioningInfoID
}
