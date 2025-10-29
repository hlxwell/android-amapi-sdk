package client

import (
	"fmt"
	"net/url"
	"time"

	"google.golang.org/api/androidmanagement/v1"

	"amapi-pkg/pkgs/amapi/types"
)

// EnterpriseService provides enterprise management methods.
type EnterpriseService struct {
	client *Client
}

// Enterprises returns the enterprise service.
func (c *Client) Enterprises() *EnterpriseService {
	return &EnterpriseService{client: c}
}

// GenerateSignupURL generates a signup URL for enterprise creation.
func (es *EnterpriseService) GenerateSignupURL(req *types.SignupURLRequest) (*types.EnterpriseSignupURL, error) {
	if req == nil {
		return nil, types.NewError(types.ErrCodeInvalidInput, "signup URL request is required")
	}

	if req.ProjectID == "" {
		req.ProjectID = es.client.config.ProjectID
	}

	if req.CallbackURL == "" {
		req.CallbackURL = es.client.config.CallbackURL
	}

	var result *androidmanagement.SignupUrl
	var err error

	err = es.client.executeAPICall(func() error {
		call := es.client.service.SignupUrls.Create()
		call.ProjectId(req.ProjectID)

		if req.CallbackURL != "" {
			call.CallbackUrl(req.CallbackURL)
		}

		if req.AdminEmail != "" {
			call.AdminEmail(req.AdminEmail)
		}

		result, err = call.Context(es.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, es.client.wrapAPIError(err, "generate signup URL")
	}

	signupURL := &types.EnterpriseSignupURL{
		URL:         result.Url,
		CallbackURL: req.CallbackURL,
		ProjectID:   req.ProjectID,
		CreatedAt:   time.Now(),
	}

	// Extract completion token from URL if present
	if parsedURL, err := url.Parse(result.Url); err == nil {
		if token := parsedURL.Query().Get("token"); token != "" {
			signupURL.CompletionToken = token
		}
	}

	return signupURL, nil
}

// Create creates a new enterprise.
func (es *EnterpriseService) Create(req *types.EnterpriseCreateRequest) (*types.Enterprise, error) {
	if req == nil {
		return nil, types.NewError(types.ErrCodeInvalidInput, "enterprise create request is required")
	}

	if req.SignupToken == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "signup token is required")
	}

	if req.ProjectID == "" {
		req.ProjectID = es.client.config.ProjectID
	}

	// Create enterprise object
	enterprise := &androidmanagement.Enterprise{}

	if req.ContactInfo != nil {
		enterprise.ContactInfo = req.ContactInfo.ToAMAPIContactInfo()
	}

	var result *androidmanagement.Enterprise
	var err error

	err = es.client.executeAPICall(func() error {
		call := es.client.service.Enterprises.Create(enterprise)
		call.ProjectId(req.ProjectID)
		call.SignupUrlName(req.SignupToken)

		if req.EnterpriseToken != "" {
			call.EnterpriseToken(req.EnterpriseToken)
		}

		result, err = call.Context(es.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, es.client.wrapAPIError(err, "create enterprise")
	}

	return types.FromAMAPIEnterprise(result), nil
}

// Get retrieves an enterprise by its resource name.
func (es *EnterpriseService) Get(enterpriseName string) (*types.Enterprise, error) {
	if enterpriseName == "" {
		return nil, types.ErrInvalidEnterpriseID
	}

	var result *androidmanagement.Enterprise
	var err error

	err = es.client.executeAPICall(func() error {
		result, err = es.client.service.Enterprises.Get(enterpriseName).Context(es.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, es.client.wrapAPIError(err, "get enterprise")
	}

	return types.FromAMAPIEnterprise(result), nil
}

// GetByID retrieves an enterprise by its ID.
func (es *EnterpriseService) GetByID(enterpriseID string) (*types.Enterprise, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	enterpriseName := buildEnterpriseName(enterpriseID)
	return es.Get(enterpriseName)
}

// Update updates an enterprise.
func (es *EnterpriseService) Update(enterpriseName string, req *types.EnterpriseUpdateRequest) (*types.Enterprise, error) {
	if enterpriseName == "" {
		return nil, types.ErrInvalidEnterpriseID
	}

	if req == nil {
		return nil, types.NewError(types.ErrCodeInvalidInput, "enterprise update request is required")
	}

	// Get current enterprise
	current, err := es.Get(enterpriseName)
	if err != nil {
		return nil, err
	}

	// Apply updates
	enterprise := current.ToAMAPIEnterprise()

	// Note: DisplayName is not supported in the androidmanagement.Enterprise struct
	// It may be part of the name field or handled separately

	if req.PrimaryColor != nil {
		enterprise.PrimaryColor = *req.PrimaryColor
	}

	if req.Logo != nil {
		enterprise.Logo = req.Logo
	}

	if req.ContactInfo != nil {
		enterprise.ContactInfo = req.ContactInfo.ToAMAPIContactInfo()
	}

	if req.EnabledNotificationTypes != nil {
		enterprise.EnabledNotificationTypes = req.EnabledNotificationTypes
	}

	if req.AppAutoApprovalEnabled != nil {
		enterprise.AppAutoApprovalEnabled = *req.AppAutoApprovalEnabled
	}

	if req.TermsAndConditions != nil {
		// Convert terms and conditions
		terms := make([]*androidmanagement.TermsAndConditions, len(req.TermsAndConditions))
		for i, tc := range req.TermsAndConditions {
			terms[i] = &androidmanagement.TermsAndConditions{
				Content: tc.Content,
				Header:  tc.Header,
			}
		}
		enterprise.TermsAndConditions = terms
	}

	var result *androidmanagement.Enterprise

	err = es.client.executeAPICall(func() error {
		result, err = es.client.service.Enterprises.Patch(enterpriseName, enterprise).Context(es.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, es.client.wrapAPIError(err, "update enterprise")
	}

	return types.FromAMAPIEnterprise(result), nil
}

// List lists enterprises in the project.
func (es *EnterpriseService) List(req *types.EnterpriseListRequest) (*types.ListResult[types.Enterprise], error) {
	if req == nil {
		req = &types.EnterpriseListRequest{}
	}

	if req.ProjectID == "" {
		req.ProjectID = es.client.config.ProjectID
	}

	var result *androidmanagement.ListEnterprisesResponse
	var err error

	err = es.client.executeAPICall(func() error {
		call := es.client.service.Enterprises.List()
		call.ProjectId(req.ProjectID)

		if req.PageSize > 0 {
			call.PageSize(int64(req.PageSize))
		}

		if req.PageToken != "" {
			call.PageToken(req.PageToken)
		}

		result, err = call.Context(es.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, es.client.wrapAPIError(err, "list enterprises")
	}

	// Convert results
	enterprises := make([]types.Enterprise, len(result.Enterprises))
	for i, enterprise := range result.Enterprises {
		enterprises[i] = *types.FromAMAPIEnterprise(enterprise)
	}

	return &types.ListResult[types.Enterprise]{
		Items:         enterprises,
		NextPageToken: result.NextPageToken,
	}, nil
}

// Delete deletes an enterprise.
func (es *EnterpriseService) Delete(req *types.EnterpriseDeleteRequest) error {
	if req == nil || req.Name == "" {
		return types.ErrInvalidEnterpriseID
	}

	err := es.client.executeAPICall(func() error {
		_, err := es.client.service.Enterprises.Delete(req.Name).Context(es.client.ctx).Do()
		return err
	})

	if err != nil {
		return es.client.wrapAPIError(err, "delete enterprise")
	}

	return nil
}

// DeleteByID deletes an enterprise by its ID.
func (es *EnterpriseService) DeleteByID(enterpriseID string, force bool) error {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return err
	}

	enterpriseName := buildEnterpriseName(enterpriseID)
	return es.Delete(&types.EnterpriseDeleteRequest{
		Name:  enterpriseName,
		Force: force,
	})
}

// CompleteSignup completes the enterprise signup process.
func (es *EnterpriseService) CompleteSignup(signupToken, enterpriseToken string) (*types.Enterprise, error) {
	if signupToken == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "signup token is required")
	}

	req := &types.EnterpriseCreateRequest{
		SignupToken:     signupToken,
		EnterpriseToken: enterpriseToken,
		ProjectID:       es.client.config.ProjectID,
	}

	return es.Create(req)
}

// EnableNotifications enables specific notification types for an enterprise.
func (es *EnterpriseService) EnableNotifications(enterpriseName string, notificationTypes []string) (*types.Enterprise, error) {
	if enterpriseName == "" {
		return nil, types.ErrInvalidEnterpriseID
	}

	if len(notificationTypes) == 0 {
		return nil, types.NewError(types.ErrCodeInvalidInput, "notification types are required")
	}

	// Get current enterprise to merge notification types
	current, err := es.Get(enterpriseName)
	if err != nil {
		return nil, err
	}

	// Merge notification types
	enabledTypes := make(map[string]bool)
	for _, nt := range current.EnabledNotificationTypes {
		enabledTypes[nt] = true
	}
	for _, nt := range notificationTypes {
		enabledTypes[nt] = true
	}

	// Convert back to slice
	var allTypes []string
	for nt := range enabledTypes {
		allTypes = append(allTypes, nt)
	}

	req := &types.EnterpriseUpdateRequest{
		EnabledNotificationTypes: allTypes,
	}

	return es.Update(enterpriseName, req)
}

// DisableNotifications disables specific notification types for an enterprise.
func (es *EnterpriseService) DisableNotifications(enterpriseName string, notificationTypes []string) (*types.Enterprise, error) {
	if enterpriseName == "" {
		return nil, types.ErrInvalidEnterpriseID
	}

	if len(notificationTypes) == 0 {
		return nil, types.NewError(types.ErrCodeInvalidInput, "notification types are required")
	}

	// Get current enterprise
	current, err := es.Get(enterpriseName)
	if err != nil {
		return nil, err
	}

	// Remove specified notification types
	disabledTypes := make(map[string]bool)
	for _, nt := range notificationTypes {
		disabledTypes[nt] = true
	}

	var remainingTypes []string
	for _, nt := range current.EnabledNotificationTypes {
		if !disabledTypes[nt] {
			remainingTypes = append(remainingTypes, nt)
		}
	}

	req := &types.EnterpriseUpdateRequest{
		EnabledNotificationTypes: remainingTypes,
	}

	return es.Update(enterpriseName, req)
}

// SetPubSubTopic sets the Pub/Sub topic for enterprise notifications.
func (es *EnterpriseService) SetPubSubTopic(enterpriseName, topicName string) (*types.Enterprise, error) {
	if enterpriseName == "" {
		return nil, types.ErrInvalidEnterpriseID
	}

	if topicName == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "topic name is required")
	}

	// Get current enterprise
	current, err := es.Get(enterpriseName)
	if err != nil {
		return nil, err
	}

	// Update enterprise with new topic
	enterprise := current.ToAMAPIEnterprise()
	enterprise.PubsubTopic = topicName

	var result *androidmanagement.Enterprise

	err = es.client.executeAPICall(func() error {
		result, err = es.client.service.Enterprises.Patch(enterpriseName, enterprise).Context(es.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, es.client.wrapAPIError(err, "set pub/sub topic")
	}

	return types.FromAMAPIEnterprise(result), nil
}

// GetApplication retrieves a specific application by package name for an enterprise.
func (es *EnterpriseService) GetApplication(enterpriseName string, packageName string) (*androidmanagement.Application, error) {
	if enterpriseName == "" {
		return nil, types.ErrInvalidEnterpriseID
	}

	if packageName == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "package name is required")
	}

	var result *androidmanagement.Application
	var err error

	err = es.client.executeAPICall(func() error {
		// Build the application name: enterprises/{enterprise}/applications/{package}
		appName := fmt.Sprintf("%s/applications/%s", enterpriseName, packageName)
		result, err = es.client.service.Enterprises.Applications.Get(appName).Context(es.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, es.client.wrapAPIError(err, "get application")
	}

	return result, nil
}

// GenerateEnterpriseUpgradeURL generates an upgrade URL for an existing enterprise.
// Note: This method is a placeholder as the actual API method may not be available
func (es *EnterpriseService) GenerateEnterpriseUpgradeURL(req *types.EnterpriseUpgradeURLRequest) (*types.EnterpriseUpgradeURL, error) {
	if req == nil {
		return nil, types.NewError(types.ErrCodeInvalidInput, "enterprise upgrade URL request is required")
	}

	if req.EnterpriseName == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "enterprise name is required")
	}

	if req.ProjectID == "" {
		req.ProjectID = es.client.config.ProjectID
	}

	// For now, return a placeholder URL
	// In a real implementation, this would call the actual API
	upgradeURL := &types.EnterpriseUpgradeURL{
		URL:            "https://play.google.com/console/developers/upgrade?project=" + req.ProjectID,
		EnterpriseName: req.EnterpriseName,
		ProjectID:      req.ProjectID,
		CreatedAt:      time.Now(),
	}

	return upgradeURL, nil
}
