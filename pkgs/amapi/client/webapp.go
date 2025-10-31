package client

import (
	"strings"

	"google.golang.org/api/androidmanagement/v1"

	"amapi-pkg/pkgs/amapi/types"
)

// WebAppService provides web app management methods.
type WebAppService struct {
	client *Client
}

// WebApps returns the web app service.
func (c *Client) WebApps() *WebAppService {
	return &WebAppService{client: c}
}

// Create creates a new web app.
func (was *WebAppService) Create(enterpriseName, startURL string, icons []*androidmanagement.WebAppIcon, versionCode int64) (*androidmanagement.WebApp, error) {
	if enterpriseName == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "enterprise name is required")
	}

	if startURL == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "start URL is required")
	}

	// Create web app object
	webApp := &androidmanagement.WebApp{
		StartUrl:    startURL,
		Icons:       icons,
		VersionCode: versionCode,
	}

	var result *androidmanagement.WebApp
	var err error

	err = was.client.executeAPICall(func() error {
		result, err = was.client.service.Enterprises.WebApps.Create(enterpriseName, webApp).Context(was.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, was.client.wrapAPIError(err, "create web app")
	}

	return result, nil
}

// CreateByEnterpriseID creates a new web app using enterprise ID.
func (was *WebAppService) CreateByEnterpriseID(enterpriseID, displayName, startURL string) (*androidmanagement.WebApp, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	enterpriseName := buildEnterpriseName(enterpriseID)
	return was.Create(enterpriseName, startURL, nil, 0)
}

// Get retrieves a web app by its resource name.
func (was *WebAppService) Get(webAppName string) (*androidmanagement.WebApp, error) {
	if webAppName == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "web app name is required")
	}

	var result *androidmanagement.WebApp
	var err error

	err = was.client.executeAPICall(func() error {
		result, err = was.client.service.Enterprises.WebApps.Get(webAppName).Context(was.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, was.client.wrapAPIError(err, "get web app")
	}

	return result, nil
}

// GetByID retrieves a web app by enterprise ID and web app ID.
func (was *WebAppService) GetByID(enterpriseID, webAppID string) (*androidmanagement.WebApp, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	if err := validateWebAppID(webAppID); err != nil {
		return nil, err
	}

	webAppName := buildWebAppName(enterpriseID, webAppID)
	return was.Get(webAppName)
}

// Update updates an existing web app.
func (was *WebAppService) Update(webAppName string, webApp *androidmanagement.WebApp, updateMask []string) (*androidmanagement.WebApp, error) {
	if webAppName == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "web app name is required")
	}

	if webApp == nil {
		return nil, types.NewError(types.ErrCodeInvalidInput, "web app is required")
	}

	var result *androidmanagement.WebApp
	var err error

	err = was.client.executeAPICall(func() error {
		call := was.client.service.Enterprises.WebApps.Patch(webAppName, webApp)

		if len(updateMask) > 0 {
			// Set update mask if provided - use comma-separated string
			maskString := strings.Join(updateMask, ",")
			call.UpdateMask(maskString)
		}

		result, err = call.Context(was.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, was.client.wrapAPIError(err, "update web app")
	}

	return result, nil
}

// UpdateByID updates a web app by enterprise ID and web app ID.
func (was *WebAppService) UpdateByID(enterpriseID, webAppID string, webApp *androidmanagement.WebApp, updateMask []string) (*androidmanagement.WebApp, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	if err := validateWebAppID(webAppID); err != nil {
		return nil, err
	}

	webAppName := buildWebAppName(enterpriseID, webAppID)
	return was.Update(webAppName, webApp, updateMask)
}

// List lists web apps for an enterprise.
func (was *WebAppService) List(enterpriseName string, pageSize int, pageToken string) (*types.ListResult[*androidmanagement.WebApp], error) {
	if enterpriseName == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "enterprise name is required")
	}

	var result *androidmanagement.ListWebAppsResponse
	var err error

	err = was.client.executeAPICall(func() error {
		call := was.client.service.Enterprises.WebApps.List(enterpriseName)

		if pageSize > 0 {
			call.PageSize(int64(pageSize))
		}

		if pageToken != "" {
			call.PageToken(pageToken)
		}

		result, err = call.Context(was.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, was.client.wrapAPIError(err, "list web apps")
	}

	// Return results directly
	webApps := make([]*androidmanagement.WebApp, len(result.WebApps))
	copy(webApps, result.WebApps)

	return &types.ListResult[*androidmanagement.WebApp]{
		Items:         webApps,
		NextPageToken: result.NextPageToken,
	}, nil
}

// ListByEnterpriseID lists web apps for an enterprise by enterprise ID.
func (was *WebAppService) ListByEnterpriseID(enterpriseID string, pageSize int, pageToken string) (*types.ListResult[*androidmanagement.WebApp], error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	enterpriseName := buildEnterpriseName(enterpriseID)
	return was.List(enterpriseName, pageSize, pageToken)
}

// Delete deletes a web app.
func (was *WebAppService) Delete(webAppName string) error {
	if webAppName == "" {
		return types.NewError(types.ErrCodeInvalidInput, "web app name is required")
	}

	err := was.client.executeAPICall(func() error {
		_, err := was.client.service.Enterprises.WebApps.Delete(webAppName).Context(was.client.ctx).Do()
		return err
	})

	if err != nil {
		return was.client.wrapAPIError(err, "delete web app")
	}

	return nil
}

// DeleteByID deletes a web app by enterprise ID and web app ID.
func (was *WebAppService) DeleteByID(enterpriseID, webAppID string) error {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return err
	}

	if err := validateWebAppID(webAppID); err != nil {
		return err
	}

	webAppName := buildWebAppName(enterpriseID, webAppID)
	return was.Delete(webAppName)
}

// GetActiveWebApps returns all active web apps for an enterprise.
func (was *WebAppService) GetActiveWebApps(enterpriseID string) (*types.ListResult[*androidmanagement.WebApp], error) {
	enterpriseName := buildEnterpriseName(enterpriseID)
	return was.List(enterpriseName, 0, "")
}

// Helper function to build web app name
func buildWebAppName(enterpriseID, webAppID string) string {
	return buildEnterpriseName(enterpriseID) + "/webApps/" + webAppID
}

// Helper function to validate web app ID
func validateWebAppID(webAppID string) error {
	if webAppID == "" {
		return types.NewError(types.ErrCodeInvalidInput, "web app ID is required")
	}
	return nil
}
