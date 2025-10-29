package client

import (
	"strings"

	"google.golang.org/api/androidmanagement/v1"

	"amapi-pkg/pkgs/amapi/types"
)

// PolicyService provides policy management methods.
type PolicyService struct {
	client *Client
}

// Policies returns the policy service.
func (c *Client) Policies() *PolicyService {
	return &PolicyService{client: c}
}

// Create creates a new policy.
func (ps *PolicyService) Create(req *types.PolicyCreateRequest) (*types.Policy, error) {
	if req == nil {
		return nil, types.NewError(types.ErrCodeInvalidInput, "policy create request is required")
	}

	if req.EnterpriseName == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "enterprise name is required")
	}

	if req.PolicyID == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "policy ID is required")
	}

	if req.Policy == nil {
		return nil, types.NewError(types.ErrCodeInvalidInput, "policy configuration is required")
	}

	// Validate policy
	if err := req.Policy.Validate(); err != nil {
		return nil, err
	}

	// Convert to AMAPI policy
	policy := types.ToAMAPIPolicy(req.Policy)

	var result *androidmanagement.Policy
	var err error

	err = ps.client.executeAPICall(func() error {
		result, err = ps.client.service.Enterprises.Policies.Patch(
			buildPolicyName(req.EnterpriseName, req.PolicyID),
			policy,
		).Context(ps.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, ps.client.wrapAPIError(err, "create policy")
	}

	return types.FromAMAPIPolicy(result), nil
}

// CreateByEnterpriseID creates a new policy using enterprise ID.
func (ps *PolicyService) CreateByEnterpriseID(enterpriseID, policyID string, policy *types.Policy) (*types.Policy, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	if err := validatePolicyID(policyID); err != nil {
		return nil, err
	}

	enterpriseName := buildEnterpriseName(enterpriseID)
	req := &types.PolicyCreateRequest{
		EnterpriseName: enterpriseName,
		PolicyID:       policyID,
		Policy:         policy,
	}

	return ps.Create(req)
}

// Get retrieves a policy by its resource name.
func (ps *PolicyService) Get(policyName string) (*types.Policy, error) {
	if policyName == "" {
		return nil, types.ErrInvalidPolicyID
	}

	var result *androidmanagement.Policy
	var err error

	err = ps.client.executeAPICall(func() error {
		result, err = ps.client.service.Enterprises.Policies.Get(policyName).Context(ps.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, ps.client.wrapAPIError(err, "get policy")
	}

	return types.FromAMAPIPolicy(result), nil
}

// GetByID retrieves a policy by enterprise ID and policy ID.
func (ps *PolicyService) GetByID(enterpriseID, policyID string) (*types.Policy, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	if err := validatePolicyID(policyID); err != nil {
		return nil, err
	}

	policyName := buildPolicyName(enterpriseID, policyID)
	return ps.Get(policyName)
}

// Update updates an existing policy.
func (ps *PolicyService) Update(req *types.PolicyUpdateRequest) (*types.Policy, error) {
	if req == nil {
		return nil, types.NewError(types.ErrCodeInvalidInput, "policy update request is required")
	}

	if req.Name == "" {
		return nil, types.ErrInvalidPolicyID
	}

	if req.Policy == nil {
		return nil, types.NewError(types.ErrCodeInvalidInput, "policy configuration is required")
	}

	// Validate policy
	if err := req.Policy.Validate(); err != nil {
		return nil, err
	}

	// Convert to AMAPI policy
	policy := types.ToAMAPIPolicy(req.Policy)

	var result *androidmanagement.Policy
	var err error

	err = ps.client.executeAPICall(func() error {
		call := ps.client.service.Enterprises.Policies.Patch(req.Name, policy)

		if len(req.UpdateMask) > 0 {
			// Set update mask if provided - use comma-separated string
			maskString := strings.Join(req.UpdateMask, ",")
			call.UpdateMask(maskString)
		}

		result, err = call.Context(ps.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, ps.client.wrapAPIError(err, "update policy")
	}

	return types.FromAMAPIPolicy(result), nil
}

// UpdateByID updates a policy by enterprise ID and policy ID.
func (ps *PolicyService) UpdateByID(enterpriseID, policyID string, policy *types.Policy) (*types.Policy, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	if err := validatePolicyID(policyID); err != nil {
		return nil, err
	}

	policyName := buildPolicyName(enterpriseID, policyID)
	req := &types.PolicyUpdateRequest{
		Name:   policyName,
		Policy: policy,
	}

	return ps.Update(req)
}

// List lists policies for an enterprise.
func (ps *PolicyService) List(req *types.PolicyListRequest) (*types.ListResult[types.Policy], error) {
	if req == nil || req.EnterpriseName == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "enterprise name is required")
	}

	var result *androidmanagement.ListPoliciesResponse
	var err error

	err = ps.client.executeAPICall(func() error {
		call := ps.client.service.Enterprises.Policies.List(req.EnterpriseName)

		if req.PageSize > 0 {
			call.PageSize(int64(req.PageSize))
		}

		if req.PageToken != "" {
			call.PageToken(req.PageToken)
		}

		result, err = call.Context(ps.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, ps.client.wrapAPIError(err, "list policies")
	}

	// Convert results
	policies := make([]types.Policy, len(result.Policies))
	for i, policy := range result.Policies {
		policies[i] = *types.FromAMAPIPolicy(policy)
	}

	return &types.ListResult[types.Policy]{
		Items:         policies,
		NextPageToken: result.NextPageToken,
	}, nil
}

// ListByEnterpriseID lists policies for an enterprise by enterprise ID.
func (ps *PolicyService) ListByEnterpriseID(enterpriseID string, options *types.ListOptions) (*types.ListResult[types.Policy], error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	enterpriseName := buildEnterpriseName(enterpriseID)
	req := &types.PolicyListRequest{
		EnterpriseName: enterpriseName,
	}

	if options != nil {
		req.ListOptions = *options
	}

	return ps.List(req)
}

// Delete deletes a policy.
func (ps *PolicyService) Delete(req *types.PolicyDeleteRequest) error {
	if req == nil || req.Name == "" {
		return types.ErrInvalidPolicyID
	}

	err := ps.client.executeAPICall(func() error {
		_, err := ps.client.service.Enterprises.Policies.Delete(req.Name).Context(ps.client.ctx).Do()
		return err
	})

	if err != nil {
		return ps.client.wrapAPIError(err, "delete policy")
	}

	return nil
}

// DeleteByID deletes a policy by enterprise ID and policy ID.
func (ps *PolicyService) DeleteByID(enterpriseID, policyID string) error {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return err
	}

	if err := validatePolicyID(policyID); err != nil {
		return err
	}

	policyName := buildPolicyName(enterpriseID, policyID)
	req := &types.PolicyDeleteRequest{
		Name: policyName,
	}

	return ps.Delete(req)
}

// Clone creates a copy of an existing policy with a new ID.
func (ps *PolicyService) Clone(sourcePolicyName, targetEnterpriseID, targetPolicyID string) (*types.Policy, error) {
	// Get the source policy
	sourcePolicy, err := ps.Get(sourcePolicyName)
	if err != nil {
		return nil, err
	}

	// Clone the policy
	clonedPolicy := sourcePolicy.Clone()

	// Clear the name and version for the new policy
	clonedPolicy.Name = ""
	clonedPolicy.Version = 0

	// Create the new policy
	req := &types.PolicyCreateRequest{
		EnterpriseName: buildEnterpriseName(targetEnterpriseID),
		PolicyID:       targetPolicyID,
		Policy:         clonedPolicy,
	}

	return ps.Create(req)
}

// AddApplication adds an application to a policy.
func (ps *PolicyService) AddApplication(policyName string, app *androidmanagement.ApplicationPolicy) (*types.Policy, error) {
	// Get current policy
	policy, err := ps.Get(policyName)
	if err != nil {
		return nil, err
	}

	// Add application
	policy.AddApplication(app)

	// Update policy
	req := &types.PolicyUpdateRequest{
		Name:   policyName,
		Policy: policy,
	}

	return ps.Update(req)
}

// RemoveApplication removes an application from a policy.
func (ps *PolicyService) RemoveApplication(policyName, packageName string) (*types.Policy, error) {
	// Get current policy
	policy, err := ps.Get(policyName)
	if err != nil {
		return nil, err
	}

	// Remove application
	policy.RemoveApplication(packageName)

	// Update policy
	req := &types.PolicyUpdateRequest{
		Name:   policyName,
		Policy: policy,
	}

	return ps.Update(req)
}

// SetApplicationInstallType sets the install type for an application in a policy.
func (ps *PolicyService) SetApplicationInstallType(policyName, packageName string, installType types.ApplicationInstallType) (*types.Policy, error) {
	// Get current policy
	policy, err := ps.Get(policyName)
	if err != nil {
		return nil, err
	}

	// Find or create application policy
	var app *androidmanagement.ApplicationPolicy
	for _, existing := range policy.Applications {
		if existing.PackageName == packageName {
			app = existing
			break
		}
	}

	if app == nil {
		// Create new application policy
		app = &androidmanagement.ApplicationPolicy{
			PackageName: packageName,
		}
		policy.Applications = append(policy.Applications, app)
	}

	// Set install type
	app.InstallType = string(installType)

	// Update policy
	req := &types.PolicyUpdateRequest{
		Name:   policyName,
		Policy: policy,
	}

	return ps.Update(req)
}

// EnableSystemApp enables a system application in a policy.
func (ps *PolicyService) EnableSystemApp(policyName, packageName string) (*types.Policy, error) {
	return ps.SetApplicationInstallType(policyName, packageName, types.InstallTypePreinstalled)
}

// BlockApplication blocks an application in a policy.
func (ps *PolicyService) BlockApplication(policyName, packageName string) (*types.Policy, error) {
	return ps.SetApplicationInstallType(policyName, packageName, types.InstallTypeBlocked)
}

// RequireApplication requires an application in a policy.
func (ps *PolicyService) RequireApplication(policyName, packageName string) (*types.Policy, error) {
	return ps.SetApplicationInstallType(policyName, packageName, types.InstallTypeRequired)
}

// SetKioskMode configures a policy for kiosk mode with a single application.
func (ps *PolicyService) SetKioskMode(policyName, kioskAppPackage string) (*types.Policy, error) {
	// Get current policy
	policy, err := ps.Get(policyName)
	if err != nil {
		return nil, err
	}

	// Configure for kiosk mode
	policy.StatusBarDisabled = true
	policy.KeyguardDisabled = true
	policy.AddUserDisabled = true
	policy.UninstallAppsDisabled = true

	// Set kiosk application
	kioskApp := types.NewKioskApp(kioskAppPackage)
	policy.AddApplication(kioskApp)

	// Update policy
	req := &types.PolicyUpdateRequest{
		Name:   policyName,
		Policy: policy,
	}

	return ps.Update(req)
}

// SetFullyManagedMode configures a policy for fully managed device mode.
func (ps *PolicyService) SetFullyManagedMode(policyName string) (*types.Policy, error) {
	// Get current policy
	policy, err := ps.Get(policyName)
	if err != nil {
		return nil, err
	}

	// Configure for fully managed mode
	policy.AddUserDisabled = true
	policy.UninstallAppsDisabled = true
	policy.StatusBarDisabled = false
	policy.KeyguardDisabled = false

	// Update policy
	req := &types.PolicyUpdateRequest{
		Name:   policyName,
		Policy: policy,
	}

	return ps.Update(req)
}

// SetWorkProfileMode configures a policy for work profile mode.
func (ps *PolicyService) SetWorkProfileMode(policyName string) (*types.Policy, error) {
	// Get current policy
	policy, err := ps.Get(policyName)
	if err != nil {
		return nil, err
	}

	// Configure for work profile mode (less restrictive)
	policy.AddUserDisabled = false
	policy.UninstallAppsDisabled = false
	policy.StatusBarDisabled = false
	policy.KeyguardDisabled = false

	// Update policy
	req := &types.PolicyUpdateRequest{
		Name:   policyName,
		Policy: policy,
	}

	return ps.Update(req)
}

// GetDevicesUsingPolicy returns devices that are using a specific policy.
func (ps *PolicyService) GetDevicesUsingPolicy(policyName string) (*types.ListResult[types.Device], error) {
	// Extract enterprise ID from policy name
	enterpriseID, _, err := parsePolicyName(policyName)
	if err != nil {
		return nil, err
	}

	// Get all devices for the enterprise
	deviceService := ps.client.Devices()
	allDevices, err := deviceService.ListByEnterpriseID(enterpriseID, nil)
	if err != nil {
		return nil, err
	}

	// Filter devices using this policy
	var devicesUsingPolicy []types.Device
	for _, device := range allDevices.Items {
		if device.AppliedPolicyName == policyName {
			devicesUsingPolicy = append(devicesUsingPolicy, device)
		}
	}

	return &types.ListResult[types.Device]{
		Items: devicesUsingPolicy,
	}, nil
}

// ValidatePolicy validates a policy configuration without saving it.
func (ps *PolicyService) ValidatePolicy(policy *types.Policy) error {
	if policy == nil {
		return types.NewError(types.ErrCodeInvalidInput, "policy is required")
	}

	return policy.Validate()
}
