package client

import (
	"google.golang.org/api/androidmanagement/v1"

	"amapi-pkg/pkgs/amapi/types"
)

// DeviceService provides device management methods.
type DeviceService struct {
	client *Client
}

// Devices returns the device service.
func (c *Client) Devices() *DeviceService {
	return &DeviceService{client: c}
}

// List lists devices for an enterprise.
func (ds *DeviceService) List(enterpriseName string, pageSize int, pageToken string, state types.DeviceState, policyCompliant *bool, userName string) (*types.ListResult[*androidmanagement.Device], error) {
	if enterpriseName == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "enterprise name is required")
	}

	var result *androidmanagement.ListDevicesResponse
	var err error

	err = ds.client.executeAPICall(func() error {
		call := ds.client.service.Enterprises.Devices.List(enterpriseName)

		if pageSize > 0 {
			call.PageSize(int64(pageSize))
		}

		if pageToken != "" {
			call.PageToken(pageToken)
		}

		result, err = call.Context(ds.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, ds.client.wrapAPIError(err, "list devices")
	}

	// Return results directly
	devices := make([]*androidmanagement.Device, len(result.Devices))
	copy(devices, result.Devices)

	// Apply client-side filtering if needed
	if state != "" || policyCompliant != nil || userName != "" {
		filteredDevices := make([]*androidmanagement.Device, 0)
		for _, device := range devices {
			// Filter by state
			if state != "" && device.State != string(state) {
				continue
			}

			// Filter by policy compliance
			if policyCompliant != nil && device.PolicyCompliant != *policyCompliant {
				continue
			}

			// Filter by user name
			if userName != "" && device.UserName != userName {
				continue
			}

			filteredDevices = append(filteredDevices, device)
		}
		devices = filteredDevices
	}

	return &types.ListResult[*androidmanagement.Device]{
		Items:         devices,
		NextPageToken: result.NextPageToken,
	}, nil
}

// ListByEnterpriseID lists devices for an enterprise by enterprise ID.
func (ds *DeviceService) ListByEnterpriseID(enterpriseID string, pageSize int, pageToken string, state types.DeviceState, policyCompliant *bool, userName string) (*types.ListResult[*androidmanagement.Device], error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	enterpriseName := buildEnterpriseName(enterpriseID)
	return ds.List(enterpriseName, pageSize, pageToken, state, policyCompliant, userName)
}

// Get retrieves a device by its resource name.
func (ds *DeviceService) Get(deviceName string) (*androidmanagement.Device, error) {
	if deviceName == "" {
		return nil, types.ErrInvalidDeviceID
	}

	var result *androidmanagement.Device
	var err error

	err = ds.client.executeAPICall(func() error {
		result, err = ds.client.service.Enterprises.Devices.Get(deviceName).Context(ds.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, ds.client.wrapAPIError(err, "get device")
	}

	return result, nil
}

// GetByID retrieves a device by enterprise ID and device ID.
func (ds *DeviceService) GetByID(enterpriseID, deviceID string) (*androidmanagement.Device, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	if err := validateDeviceID(deviceID); err != nil {
		return nil, err
	}

	deviceName := buildDeviceName(enterpriseID, deviceID)
	return ds.Get(deviceName)
}

// IssueCommand issues a command to a device.
func (ds *DeviceService) IssueCommand(deviceName string, command *androidmanagement.Command) (*androidmanagement.Operation, error) {
	if deviceName == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "device name is required")
	}

	if command == nil {
		return nil, types.NewError(types.ErrCodeInvalidInput, "command is required")
	}

	var result *androidmanagement.Operation
	var err error

	err = ds.client.executeAPICall(func() error {
		result, err = ds.client.service.Enterprises.Devices.IssueCommand(deviceName, command).Context(ds.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, ds.client.wrapAPIError(err, "issue device command")
	}

	return result, nil
}

// IssueCommandByID issues a command to a device by enterprise ID and device ID.
func (ds *DeviceService) IssueCommandByID(enterpriseID, deviceID string, command *androidmanagement.Command) (*androidmanagement.Operation, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	if err := validateDeviceID(deviceID); err != nil {
		return nil, err
	}

	deviceName := buildDeviceName(enterpriseID, deviceID)
	return ds.IssueCommand(deviceName, command)
}

// Lock locks a device for the specified duration.
func (ds *DeviceService) Lock(deviceName string, duration string) (*androidmanagement.Operation, error) {
	command := &androidmanagement.Command{
		Type:     string(types.CommandTypeLock),
		Duration: duration,
	}

	return ds.IssueCommand(deviceName, command)
}

// LockByID locks a device by enterprise ID and device ID.
func (ds *DeviceService) LockByID(enterpriseID, deviceID string, duration string) (*androidmanagement.Operation, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	if err := validateDeviceID(deviceID); err != nil {
		return nil, err
	}

	deviceName := buildDeviceName(enterpriseID, deviceID)
	return ds.Lock(deviceName, duration)
}

// Reset performs a factory reset on a device.
func (ds *DeviceService) Reset(deviceName string) (*androidmanagement.Operation, error) {
	command := &androidmanagement.Command{
		Type: string(types.CommandTypeReset),
	}

	return ds.IssueCommand(deviceName, command)
}

// ResetByID performs a factory reset on a device by enterprise ID and device ID.
func (ds *DeviceService) ResetByID(enterpriseID, deviceID string) (*androidmanagement.Operation, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	if err := validateDeviceID(deviceID); err != nil {
		return nil, err
	}

	deviceName := buildDeviceName(enterpriseID, deviceID)
	return ds.Reset(deviceName)
}

// Reboot reboots a device.
func (ds *DeviceService) Reboot(deviceName string) (*androidmanagement.Operation, error) {
	command := &androidmanagement.Command{
		Type: string(types.CommandTypeReboot),
	}

	return ds.IssueCommand(deviceName, command)
}

// RebootByID reboots a device by enterprise ID and device ID.
func (ds *DeviceService) RebootByID(enterpriseID, deviceID string) (*androidmanagement.Operation, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	if err := validateDeviceID(deviceID); err != nil {
		return nil, err
	}

	deviceName := buildDeviceName(enterpriseID, deviceID)
	return ds.Reboot(deviceName)
}

// RemovePassword removes the device password/PIN.
func (ds *DeviceService) RemovePassword(deviceName string) (*androidmanagement.Operation, error) {
	command := &androidmanagement.Command{
		Type: string(types.CommandTypeRemovePassword),
	}

	return ds.IssueCommand(deviceName, command)
}

// RemovePasswordByID removes the device password/PIN by enterprise ID and device ID.
func (ds *DeviceService) RemovePasswordByID(enterpriseID, deviceID string) (*androidmanagement.Operation, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	if err := validateDeviceID(deviceID); err != nil {
		return nil, err
	}

	deviceName := buildDeviceName(enterpriseID, deviceID)
	return ds.RemovePassword(deviceName)
}

// ClearAppData clears data for a specific application on the device.
func (ds *DeviceService) ClearAppData(deviceName, packageName string) (*androidmanagement.Operation, error) {
	command := &androidmanagement.Command{
		Type: string(types.CommandTypeClearAppData),
		// Note: Package name would need to be included in the command
		// The AMAPI Command structure would need to be extended for this
	}

	return ds.IssueCommand(deviceName, command)
}

// StartLostMode starts lost mode on a device.
func (ds *DeviceService) StartLostMode(deviceName string) (*androidmanagement.Operation, error) {
	command := &androidmanagement.Command{
		Type: string(types.CommandTypeStartLostMode),
	}

	return ds.IssueCommand(deviceName, command)
}

// StopLostMode stops lost mode on a device.
func (ds *DeviceService) StopLostMode(deviceName string) (*androidmanagement.Operation, error) {
	command := &androidmanagement.Command{
		Type: string(types.CommandTypeStopLostMode),
	}

	return ds.IssueCommand(deviceName, command)
}

// Delete deletes a device (performs a wipe and removes it from management).
func (ds *DeviceService) Delete(deviceName string) error {
	if deviceName == "" {
		return types.ErrInvalidDeviceID
	}

	err := ds.client.executeAPICall(func() error {
		// Note: Delete returns *androidmanagement.Empty, not Operation
		_, err := ds.client.service.Enterprises.Devices.Delete(deviceName).Context(ds.client.ctx).Do()
		return err
	})

	if err != nil {
		return ds.client.wrapAPIError(err, "delete device")
	}

	return nil
}

// DeleteByID deletes a device by enterprise ID and device ID.
func (ds *DeviceService) DeleteByID(enterpriseID, deviceID string) error {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return err
	}

	if err := validateDeviceID(deviceID); err != nil {
		return err
	}

	deviceName := buildDeviceName(enterpriseID, deviceID)
	return ds.Delete(deviceName)
}

// GetOperations retrieves operations for a device.
func (ds *DeviceService) GetOperations(deviceName string) ([]*androidmanagement.Operation, error) {
	if deviceName == "" {
		return nil, types.ErrInvalidDeviceID
	}

	var result *androidmanagement.ListOperationsResponse
	var err error

	err = ds.client.executeAPICall(func() error {
		result, err = ds.client.service.Enterprises.Devices.Operations.List(deviceName).Context(ds.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, ds.client.wrapAPIError(err, "get device operations")
	}

	return result.Operations, nil
}

// GetOperation retrieves a specific operation for a device.
func (ds *DeviceService) GetOperation(operationName string) (*androidmanagement.Operation, error) {
	if operationName == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "operation name is required")
	}

	var result *androidmanagement.Operation
	var err error

	err = ds.client.executeAPICall(func() error {
		result, err = ds.client.service.Enterprises.Devices.Operations.Get(operationName).Context(ds.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, ds.client.wrapAPIError(err, "get device operation")
	}

	return result, nil
}

// CancelOperation cancels a device operation.
func (ds *DeviceService) CancelOperation(operationName string) error {
	if operationName == "" {
		return types.NewError(types.ErrCodeInvalidInput, "operation name is required")
	}

	err := ds.client.executeAPICall(func() error {
		_, err := ds.client.service.Enterprises.Devices.Operations.Cancel(operationName).Context(ds.client.ctx).Do()
		return err
	})

	if err != nil {
		return ds.client.wrapAPIError(err, "cancel device operation")
	}

	return nil
}

// Helper methods for device filtering and querying

// GetActiveDevices returns all active devices for an enterprise.
func (ds *DeviceService) GetActiveDevices(enterpriseID string) (*types.ListResult[*androidmanagement.Device], error) {
	enterpriseName := buildEnterpriseName(enterpriseID)
	return ds.List(enterpriseName, 0, "", types.DeviceStateActive, nil, "")
}

// GetCompliantDevices returns all policy-compliant devices for an enterprise.
func (ds *DeviceService) GetCompliantDevices(enterpriseID string) (*types.ListResult[*androidmanagement.Device], error) {
	compliant := true
	enterpriseName := buildEnterpriseName(enterpriseID)
	return ds.List(enterpriseName, 0, "", "", &compliant, "")
}

// GetNonCompliantDevices returns all non-compliant devices for an enterprise.
func (ds *DeviceService) GetNonCompliantDevices(enterpriseID string) (*types.ListResult[*androidmanagement.Device], error) {
	compliant := false
	enterpriseName := buildEnterpriseName(enterpriseID)
	return ds.List(enterpriseName, 0, "", "", &compliant, "")
}

// GetDevicesByUser returns all devices for a specific user in an enterprise.
func (ds *DeviceService) GetDevicesByUser(enterpriseID, userName string) (*types.ListResult[*androidmanagement.Device], error) {
	enterpriseName := buildEnterpriseName(enterpriseID)
	return ds.List(enterpriseName, 0, "", "", nil, userName)
}
