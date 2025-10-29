package types

import (
	"time"

	"google.golang.org/api/androidmanagement/v1"
)

// ProvisioningInfo represents device provisioning information.
type ProvisioningInfo struct {
	// Name is the resource name of the provisioning info
	Name string `json:"name"`

	// EnterpriseID is the enterprise this provisioning info belongs to
	EnterpriseID string `json:"enterprise_id"`

	// DeviceID is the device ID this provisioning info is for
	DeviceID string `json:"device_id"`

	// PolicyName is the policy to apply to the device
	PolicyName string `json:"policy_name"`

	// EnrollmentTokenName is the enrollment token used for provisioning
	EnrollmentTokenName string `json:"enrollment_token_name"`

	// ProvisioningState is the current state of provisioning
	ProvisioningState string `json:"provisioning_state"`

	// CreatedAt is when the provisioning info was created
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is when the provisioning info was last updated
	UpdatedAt time.Time `json:"updated_at"`

	// IsActive indicates if the provisioning info is active
	IsActive bool `json:"is_active"`

	// DeviceInfo contains information about the device being provisioned
	DeviceInfo *DeviceProvisioningInfo `json:"device_info,omitempty"`

	// ErrorMessage contains any error message if provisioning failed
	ErrorMessage string `json:"error_message,omitempty"`
}

// DeviceProvisioningInfo contains device-specific provisioning information.
type DeviceProvisioningInfo struct {
	// DeviceModel is the model of the device
	DeviceModel string `json:"device_model,omitempty"`

	// AndroidVersion is the Android version of the device
	AndroidVersion string `json:"android_version,omitempty"`

	// DeviceSerialNumber is the serial number of the device
	DeviceSerialNumber string `json:"device_serial_number,omitempty"`

	// IMEI is the IMEI of the device
	IMEI string `json:"imei,omitempty"`

	// WiFiMACAddress is the WiFi MAC address of the device
	WiFiMACAddress string `json:"wifi_mac_address,omitempty"`

	// BluetoothMACAddress is the Bluetooth MAC address of the device
	BluetoothMACAddress string `json:"bluetooth_mac_address,omitempty"`
}

// ProvisioningInfoGetRequest represents a request to get provisioning info.
type ProvisioningInfoGetRequest struct {
	// Name is the provisioning info resource name
	Name string `json:"name"`
}

// ProvisioningInfo helper methods

// GetID extracts the provisioning info ID from the resource name.
func (pi *ProvisioningInfo) GetID() string {
	if pi.Name == "" {
		return ""
	}

	// Extract ID from name format: provisioningInfo/{provisioningInfoId}
	for i := len(pi.Name) - 1; i >= 0; i-- {
		if pi.Name[i] == '/' {
			return pi.Name[i+1:]
		}
	}

	return pi.Name
}

// GetEnterpriseID extracts the enterprise ID from the provisioning info.
func (pi *ProvisioningInfo) GetEnterpriseID() string {
	return pi.EnterpriseID
}

// GetDeviceID extracts the device ID from the provisioning info.
func (pi *ProvisioningInfo) GetDeviceID() string {
	return pi.DeviceID
}

// IsProvisioningComplete checks if the provisioning is complete.
func (pi *ProvisioningInfo) IsProvisioningComplete() bool {
	return pi.ProvisioningState == "PROVISIONING_COMPLETE" || pi.ProvisioningState == "ACTIVE"
}

// HasError checks if there was an error during provisioning.
func (pi *ProvisioningInfo) HasError() bool {
	return pi.ErrorMessage != "" || pi.ProvisioningState == "PROVISIONING_FAILED"
}

// FromAMAPIProvisioningInfo converts an Android Management API provisioning info to our type.
func FromAMAPIProvisioningInfo(info *androidmanagement.ProvisioningInfo) *ProvisioningInfo {
	if info == nil {
		return nil
	}

	provisioningInfo := &ProvisioningInfo{
		Name:      info.Name,
		IsActive:  true, // Assume active if not specified
		CreatedAt: time.Now(), // AMAPI doesn't provide creation time
		UpdatedAt: time.Now(),
	}

	// Extract enterprise ID and device ID from name
	// This would need to be implemented based on the actual AMAPI response structure

	return provisioningInfo
}

// ToAMAPIProvisioningInfo converts our provisioning info to Android Management API format.
func (pi *ProvisioningInfo) ToAMAPIProvisioningInfo() *androidmanagement.ProvisioningInfo {
	if pi == nil {
		return nil
	}

	return &androidmanagement.ProvisioningInfo{
		Name: pi.Name,
	}
}
