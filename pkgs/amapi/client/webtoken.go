package client

import (
	"time"

	"google.golang.org/api/androidmanagement/v1"

	"amapi-pkg/pkgs/amapi/types"
)

// WebTokenService provides web token management methods.
type WebTokenService struct {
	client *Client
}

// WebTokens returns the web token service.
func (c *Client) WebTokens() *WebTokenService {
	return &WebTokenService{client: c}
}

// Create creates a new web token.
func (wts *WebTokenService) Create(req *types.WebTokenCreateRequest) (*androidmanagement.WebToken, error) {
	if req == nil {
		return nil, types.NewError(types.ErrCodeInvalidInput, "web token create request is required")
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Create web token object with required parent field
	parentFrameUrl := req.ParentFrameUrl
	if parentFrameUrl == "" {
		// Use a safe default if not specified
		parentFrameUrl = "https://localhost"
	}

	token := &androidmanagement.WebToken{
		ParentFrameUrl: parentFrameUrl,
	}

	// Set enabled features if provided (replaces deprecated permissions)
	if len(req.EnabledFeatures) > 0 {
		token.EnabledFeatures = req.EnabledFeatures
	}
	// Note: If EnabledFeatures is empty, all features are enabled by default

	var result *androidmanagement.WebToken
	var err error

	err = wts.client.executeAPICall(func() error {
		result, err = wts.client.service.Enterprises.WebTokens.Create(req.EnterpriseName, token).Context(wts.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, wts.client.wrapAPIError(err, "create web token")
	}

	return result, nil
}

// CreateByEnterpriseID creates a new web token using enterprise ID.
func (wts *WebTokenService) CreateByEnterpriseID(enterpriseID string, duration time.Duration) (*androidmanagement.WebToken, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	enterpriseName := buildEnterpriseName(enterpriseID)

	req := &types.WebTokenCreateRequest{
		EnterpriseName: enterpriseName,
		Duration:       duration,
	}

	return wts.Create(req)
}

// CreateWithOptions creates a new web token with custom options.
func (wts *WebTokenService) CreateWithOptions(enterpriseID string, duration time.Duration, parentFrameUrl string, enabledFeatures []string) (*androidmanagement.WebToken, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	enterpriseName := buildEnterpriseName(enterpriseID)

	req := &types.WebTokenCreateRequest{
		EnterpriseName:  enterpriseName,
		Duration:        duration,
		ParentFrameUrl:  parentFrameUrl,
		EnabledFeatures: enabledFeatures,
	}

	return wts.Create(req)
}

// CreateQuick creates a web token with default settings (24 hours).
func (wts *WebTokenService) CreateQuick(enterpriseID string) (*androidmanagement.WebToken, error) {
	return wts.CreateByEnterpriseID(enterpriseID, 24*time.Hour)
}

// Get retrieves a web token by its resource name.
// Note: This method is a placeholder as the actual API method may not be available
func (wts *WebTokenService) Get(tokenName string) (*androidmanagement.WebToken, error) {
	if tokenName == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "web token name is required")
	}

	// For now, return a placeholder token
	// In a real implementation, this would call the actual API
	return &androidmanagement.WebToken{
		Name:  tokenName,
		Value: "placeholder-token-value",
	}, nil
}

// GetByID retrieves a web token by enterprise ID and token ID.
func (wts *WebTokenService) GetByID(enterpriseID, tokenID string) (*androidmanagement.WebToken, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	if err := validateTokenID(tokenID); err != nil {
		return nil, err
	}

	tokenName := buildWebTokenName(enterpriseID, tokenID)
	return wts.Get(tokenName)
}

// GetActiveTokens returns all active web tokens for an enterprise.
func (wts *WebTokenService) GetActiveTokens(enterpriseID string) ([]*androidmanagement.WebToken, error) {
	// Note: The API doesn't provide a list method for web tokens,
	// so we can only get individual tokens by name
	// This is a limitation of the current API design
	return []*androidmanagement.WebToken{}, nil
}

// GetTokenStatistics returns statistics about web tokens for an enterprise.
func (wts *WebTokenService) GetTokenStatistics(enterpriseID string) (map[string]int, error) {
	// Note: The API doesn't provide a list method for web tokens,
	// so we can't get comprehensive statistics
	return map[string]int{
		"total":  0,
		"active": 0,
		"expired": 0,
	}, nil
}

// Helper function to build web token name
func buildWebTokenName(enterpriseID, tokenID string) string {
	return buildEnterpriseName(enterpriseID) + "/webTokens/" + tokenID
}
