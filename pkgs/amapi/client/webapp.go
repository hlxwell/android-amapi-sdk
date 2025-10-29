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
func (was *WebAppService) Create(req *types.WebAppCreateRequest) (*types.WebApp, error) {
	if req == nil {
		return nil, types.NewError(types.ErrCodeInvalidInput, "web app create request is required")
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Create web app object
	webApp := &androidmanagement.WebApp{
		StartUrl:    req.StartURL,
		Icons:       req.Icons,
		VersionCode: req.VersionCode,
	}

	var result *androidmanagement.WebApp
	var err error

	err = was.client.executeAPICall(func() error {
		result, err = was.client.service.Enterprises.WebApps.Create(req.EnterpriseName, webApp).Context(was.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, was.client.wrapAPIError(err, "create web app")
	}

	return types.FromAMAPIWebApp(result), nil
}

// CreateByEnterpriseID creates a new web app using enterprise ID.
func (was *WebAppService) CreateByEnterpriseID(enterpriseID, displayName, startURL string) (*types.WebApp, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	enterpriseName := buildEnterpriseName(enterpriseID)

	req := &types.WebAppCreateRequest{
		EnterpriseName: enterpriseName,
		DisplayName:    displayName,
		StartURL:       startURL,
	}

	return was.Create(req)
}

// Get retrieves a web app by its resource name.
func (was *WebAppService) Get(webAppName string) (*types.WebApp, error) {
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

	return types.FromAMAPIWebApp(result), nil
}

// GetByID retrieves a web app by enterprise ID and web app ID.
func (was *WebAppService) GetByID(enterpriseID, webAppID string) (*types.WebApp, error) {
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
func (was *WebAppService) Update(req *types.WebAppUpdateRequest) (*types.WebApp, error) {
	if req == nil {
		return nil, types.NewError(types.ErrCodeInvalidInput, "web app update request is required")
	}

	if req.Name == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "web app name is required")
	}

	// Get current web app
	currentWebApp, err := was.Get(req.Name)
	if err != nil {
		return nil, err
	}

	// Apply updates
	webApp := currentWebApp.ToAMAPIWebApp()

	// Note: DisplayName field may not be available in the API

	if req.StartURL != "" {
		webApp.StartUrl = req.StartURL
	}

	if req.Icons != nil {
		webApp.Icons = req.Icons
	}

	if req.VersionCode != 0 {
		webApp.VersionCode = req.VersionCode
	}

	var result *androidmanagement.WebApp

	err = was.client.executeAPICall(func() error {
		call := was.client.service.Enterprises.WebApps.Patch(req.Name, webApp)

		if len(req.UpdateMask) > 0 {
			// Set update mask if provided - use comma-separated string
			maskString := strings.Join(req.UpdateMask, ",")
			call.UpdateMask(maskString)
		}

		result, err = call.Context(was.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, was.client.wrapAPIError(err, "update web app")
	}

	return types.FromAMAPIWebApp(result), nil
}

// UpdateByID updates a web app by enterprise ID and web app ID.
func (was *WebAppService) UpdateByID(enterpriseID, webAppID string, webApp *types.WebApp) (*types.WebApp, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	if err := validateWebAppID(webAppID); err != nil {
		return nil, err
	}

	webAppName := buildWebAppName(enterpriseID, webAppID)
	req := &types.WebAppUpdateRequest{
		Name:        webAppName,
		DisplayName: webApp.DisplayName,
		StartURL:    webApp.StartURL,
		Icons:       webApp.Icons,
		VersionCode: webApp.VersionCode,
	}

	return was.Update(req)
}

// List lists web apps for an enterprise.
func (was *WebAppService) List(req *types.WebAppListRequest) (*types.ListResult[types.WebApp], error) {
	if req == nil || req.EnterpriseName == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "enterprise name is required")
	}

	var result *androidmanagement.ListWebAppsResponse
	var err error

	err = was.client.executeAPICall(func() error {
		call := was.client.service.Enterprises.WebApps.List(req.EnterpriseName)

		if req.PageSize > 0 {
			call.PageSize(int64(req.PageSize))
		}

		if req.PageToken != "" {
			call.PageToken(req.PageToken)
		}

		result, err = call.Context(was.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, was.client.wrapAPIError(err, "list web apps")
	}

	// Convert results
	webApps := make([]types.WebApp, len(result.WebApps))
	for i, webApp := range result.WebApps {
		webApps[i] = *types.FromAMAPIWebApp(webApp)
	}

	// Apply client-side filtering
	if req.ActiveOnly {
		filteredWebApps := make([]types.WebApp, 0)
		for _, webApp := range webApps {
			if webApp.IsActive {
				filteredWebApps = append(filteredWebApps, webApp)
			}
		}
		webApps = filteredWebApps
	}

	return &types.ListResult[types.WebApp]{
		Items:         webApps,
		NextPageToken: result.NextPageToken,
	}, nil
}

// ListByEnterpriseID lists web apps for an enterprise by enterprise ID.
func (was *WebAppService) ListByEnterpriseID(enterpriseID string, options *types.ListOptions) (*types.ListResult[types.WebApp], error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	enterpriseName := buildEnterpriseName(enterpriseID)
	req := &types.WebAppListRequest{
		EnterpriseName: enterpriseName,
	}

	if options != nil {
		req.ListOptions = *options
	}

	return was.List(req)
}

// Delete deletes a web app.
func (was *WebAppService) Delete(req *types.WebAppDeleteRequest) error {
	if req == nil || req.Name == "" {
		return types.NewError(types.ErrCodeInvalidInput, "web app name is required")
	}

	err := was.client.executeAPICall(func() error {
		_, err := was.client.service.Enterprises.WebApps.Delete(req.Name).Context(was.client.ctx).Do()
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
	req := &types.WebAppDeleteRequest{
		Name: webAppName,
	}

	return was.Delete(req)
}

// GetActiveWebApps returns all active web apps for an enterprise.
func (was *WebAppService) GetActiveWebApps(enterpriseID string) (*types.ListResult[types.WebApp], error) {
	req := &types.WebAppListRequest{
		EnterpriseName: buildEnterpriseName(enterpriseID),
		ActiveOnly:     true,
	}

	return was.List(req)
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
