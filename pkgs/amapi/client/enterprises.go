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
func (es *EnterpriseService) GenerateSignupURL(projectID, callbackURL, adminEmail, enterpriseDisplayName, locale string) (*types.EnterpriseSignupURL, error) {
	if projectID == "" {
		projectID = es.client.config.ProjectID
	}

	if callbackURL == "" {
		callbackURL = es.client.config.CallbackURL
	}

	var result *androidmanagement.SignupUrl
	var err error

	err = es.client.executeAPICall(func() error {
		call := es.client.service.SignupUrls.Create()
		call.ProjectId(projectID)

		if callbackURL != "" {
			call.CallbackUrl(callbackURL)
		}

		if adminEmail != "" {
			call.AdminEmail(adminEmail)
		}

		// Note: EnterpriseDisplayName and Locale are not available in the API
		// They are accepted as parameters but not used

		result, err = call.Context(es.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, es.client.wrapAPIError(err, "generate signup URL")
	}

	signupURL := &types.EnterpriseSignupURL{
		URL:         result.Url,
		CallbackURL: callbackURL,
		ProjectID:   projectID,
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
func (es *EnterpriseService) Create(signupToken, projectID, enterpriseToken string, contactInfo *androidmanagement.ContactInfo) (*androidmanagement.Enterprise, error) {
	if signupToken == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "signup token is required")
	}

	if projectID == "" {
		projectID = es.client.config.ProjectID
	}

	// Create enterprise object
	enterprise := &androidmanagement.Enterprise{}
	if contactInfo != nil {
		enterprise.ContactInfo = contactInfo
	}

	var result *androidmanagement.Enterprise
	var err error

	err = es.client.executeAPICall(func() error {
		call := es.client.service.Enterprises.Create(enterprise)
		call.ProjectId(projectID)
		call.SignupUrlName(signupToken)

		if enterpriseToken != "" {
			call.EnterpriseToken(enterpriseToken)
		}

		result, err = call.Context(es.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, es.client.wrapAPIError(err, "create enterprise")
	}

	return result, nil
}

// Get retrieves an enterprise by its resource name.
func (es *EnterpriseService) Get(enterpriseName string) (*androidmanagement.Enterprise, error) {
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

	return result, nil
}

// GetByID retrieves an enterprise by its ID.
func (es *EnterpriseService) GetByID(enterpriseID string) (*androidmanagement.Enterprise, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	enterpriseName := buildEnterpriseName(enterpriseID)
	return es.Get(enterpriseName)
}

// Update updates an enterprise.
func (es *EnterpriseService) Update(enterpriseName string, primaryColor *int64, logo *androidmanagement.ExternalData, contactInfo *androidmanagement.ContactInfo, enabledNotificationTypes []string, appAutoApprovalEnabled *bool, termsAndConditions []*androidmanagement.TermsAndConditions) (*androidmanagement.Enterprise, error) {
	if enterpriseName == "" {
		return nil, types.ErrInvalidEnterpriseID
	}

	// Get current enterprise
	current, err := es.Get(enterpriseName)
	if err != nil {
		return nil, err
	}

	// Apply updates if provided
	if primaryColor != nil {
		current.PrimaryColor = *primaryColor
	}

	if logo != nil {
		current.Logo = logo
	}

	if contactInfo != nil {
		current.ContactInfo = contactInfo
	}

	if enabledNotificationTypes != nil {
		current.EnabledNotificationTypes = enabledNotificationTypes
	}

	if appAutoApprovalEnabled != nil {
		current.AppAutoApprovalEnabled = *appAutoApprovalEnabled
	}

	if termsAndConditions != nil {
		current.TermsAndConditions = termsAndConditions
	}

	var result *androidmanagement.Enterprise

	err = es.client.executeAPICall(func() error {
		result, err = es.client.service.Enterprises.Patch(enterpriseName, current).Context(es.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, es.client.wrapAPIError(err, "update enterprise")
	}

	return result, nil
}

// List lists enterprises in the project.
func (es *EnterpriseService) List(projectID string, pageSize int, pageToken string) (*types.ListResult[*androidmanagement.Enterprise], error) {
	if projectID == "" {
		projectID = es.client.config.ProjectID
	}

	var result *androidmanagement.ListEnterprisesResponse
	var err error

	err = es.client.executeAPICall(func() error {
		call := es.client.service.Enterprises.List()
		call.ProjectId(projectID)

		if pageSize > 0 {
			call.PageSize(int64(pageSize))
		}

		if pageToken != "" {
			call.PageToken(pageToken)
		}

		result, err = call.Context(es.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, es.client.wrapAPIError(err, "list enterprises")
	}

	// Return results directly
	enterprises := make([]*androidmanagement.Enterprise, len(result.Enterprises))
	copy(enterprises, result.Enterprises)

	return &types.ListResult[*androidmanagement.Enterprise]{
		Items:         enterprises,
		NextPageToken: result.NextPageToken,
	}, nil
}

// Delete deletes an enterprise.
func (es *EnterpriseService) Delete(enterpriseName string) error {
	if enterpriseName == "" {
		return types.ErrInvalidEnterpriseID
	}

	err := es.client.executeAPICall(func() error {
		_, err := es.client.service.Enterprises.Delete(enterpriseName).Context(es.client.ctx).Do()
		return err
	})

	if err != nil {
		return es.client.wrapAPIError(err, "delete enterprise")
	}

	return nil
}

// DeleteByID deletes an enterprise by its ID.
func (es *EnterpriseService) DeleteByID(enterpriseID string) error {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return err
	}

	enterpriseName := buildEnterpriseName(enterpriseID)
	return es.Delete(enterpriseName)
}

// CompleteSignup completes the enterprise signup process.
func (es *EnterpriseService) CompleteSignup(signupToken, enterpriseToken string) (*androidmanagement.Enterprise, error) {
	if signupToken == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "signup token is required")
	}

	return es.Create(signupToken, es.client.config.ProjectID, enterpriseToken, nil)
}

// EnableNotifications enables specific notification types for an enterprise.
func (es *EnterpriseService) EnableNotifications(enterpriseName string, notificationTypes []string) (*androidmanagement.Enterprise, error) {
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

	return es.Update(enterpriseName, nil, nil, nil, allTypes, nil, nil)
}

// DisableNotifications disables specific notification types for an enterprise.
func (es *EnterpriseService) DisableNotifications(enterpriseName string, notificationTypes []string) (*androidmanagement.Enterprise, error) {
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

	return es.Update(enterpriseName, nil, nil, nil, remainingTypes, nil, nil)
}

// SetPubSubTopic sets the Pub/Sub topic for enterprise notifications.
func (es *EnterpriseService) SetPubSubTopic(enterpriseName, topicName string) (*androidmanagement.Enterprise, error) {
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
	current.PubsubTopic = topicName

	var result *androidmanagement.Enterprise

	err = es.client.executeAPICall(func() error {
		result, err = es.client.service.Enterprises.Patch(enterpriseName, current).Context(es.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, es.client.wrapAPIError(err, "set pub/sub topic")
	}

	return result, nil
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
func (es *EnterpriseService) GenerateEnterpriseUpgradeURL(enterpriseName, projectID, callbackURL, adminEmail, locale string) (*types.EnterpriseUpgradeURL, error) {
	if enterpriseName == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "enterprise name is required")
	}

	if projectID == "" {
		projectID = es.client.config.ProjectID
	}

	// For now, return a placeholder URL
	// In a real implementation, this would call the actual API
	upgradeURL := &types.EnterpriseUpgradeURL{
		URL:            "https://play.google.com/console/developers/upgrade?project=" + projectID,
		EnterpriseName: enterpriseName,
		ProjectID:      projectID,
		CreatedAt:      time.Now(),
	}

	return upgradeURL, nil
}
