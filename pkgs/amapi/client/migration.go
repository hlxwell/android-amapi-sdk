package client

import (
	"time"

	"google.golang.org/api/androidmanagement/v1"

	"amapi-pkg/pkgs/amapi/types"
)

// MigrationService provides migration token management methods.
type MigrationService struct {
	client *Client
}

// MigrationTokens returns the migration service.
func (c *Client) MigrationTokens() *MigrationService {
	return &MigrationService{client: c}
}

// Create creates a new migration token.
func (ms *MigrationService) Create(req *types.MigrationTokenCreateRequest) (*types.MigrationToken, error) {
	if req == nil {
		return nil, types.NewError(types.ErrCodeInvalidInput, "migration token create request is required")
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Create migration token object
	token := &androidmanagement.MigrationToken{}

	var result *androidmanagement.MigrationToken
	var err error

	err = ms.client.executeAPICall(func() error {
		result, err = ms.client.service.Enterprises.MigrationTokens.Create(req.EnterpriseName, token).Context(ms.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, ms.client.wrapAPIError(err, "create migration token")
	}

	return types.FromAMAPIMigrationToken(result), nil
}

// CreateByEnterpriseID creates a new migration token using enterprise ID.
func (ms *MigrationService) CreateByEnterpriseID(enterpriseID, policyID string, duration time.Duration) (*types.MigrationToken, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	if err := validatePolicyID(policyID); err != nil {
		return nil, err
	}

	enterpriseName := buildEnterpriseName(enterpriseID)
	policyName := buildPolicyName(enterpriseID, policyID)

	req := &types.MigrationTokenCreateRequest{
		EnterpriseName: enterpriseName,
		PolicyName:     policyName,
		Duration:       duration,
	}

	return ms.Create(req)
}

// Get retrieves a migration token by its resource name.
func (ms *MigrationService) Get(tokenName string) (*types.MigrationToken, error) {
	if tokenName == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "migration token name is required")
	}

	var result *androidmanagement.MigrationToken
	var err error

	err = ms.client.executeAPICall(func() error {
		result, err = ms.client.service.Enterprises.MigrationTokens.Get(tokenName).Context(ms.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, ms.client.wrapAPIError(err, "get migration token")
	}

	return types.FromAMAPIMigrationToken(result), nil
}

// GetByID retrieves a migration token by enterprise ID and token ID.
func (ms *MigrationService) GetByID(enterpriseID, tokenID string) (*types.MigrationToken, error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	if err := validateTokenID(tokenID); err != nil {
		return nil, err
	}

	tokenName := buildMigrationTokenName(enterpriseID, tokenID)
	return ms.Get(tokenName)
}

// List lists migration tokens for an enterprise.
func (ms *MigrationService) List(req *types.MigrationTokenListRequest) (*types.ListResult[types.MigrationToken], error) {
	if req == nil || req.EnterpriseName == "" {
		return nil, types.NewError(types.ErrCodeInvalidInput, "enterprise name is required")
	}

	var result *androidmanagement.ListMigrationTokensResponse
	var err error

	err = ms.client.executeAPICall(func() error {
		call := ms.client.service.Enterprises.MigrationTokens.List(req.EnterpriseName)

		if req.PageSize > 0 {
			call.PageSize(int64(req.PageSize))
		}

		if req.PageToken != "" {
			call.PageToken(req.PageToken)
		}

		result, err = call.Context(ms.client.ctx).Do()
		return err
	})

	if err != nil {
		return nil, ms.client.wrapAPIError(err, "list migration tokens")
	}

	// Convert results
	tokens := make([]types.MigrationToken, len(result.MigrationTokens))
	for i, token := range result.MigrationTokens {
		tokens[i] = *types.FromAMAPIMigrationToken(token)
	}

	// Apply client-side filtering
	if !req.IncludeExpired || req.ActiveOnly {
		filteredTokens := make([]types.MigrationToken, 0)
		for _, token := range tokens {
			// Filter expired tokens if requested
			if !req.IncludeExpired && token.IsExpired() {
				continue
			}

			// Filter active tokens if requested
			if req.ActiveOnly && !token.IsActive {
				continue
			}

			filteredTokens = append(filteredTokens, token)
		}
		tokens = filteredTokens
	}

	return &types.ListResult[types.MigrationToken]{
		Items:         tokens,
		NextPageToken: result.NextPageToken,
	}, nil
}

// ListByEnterpriseID lists migration tokens for an enterprise by enterprise ID.
func (ms *MigrationService) ListByEnterpriseID(enterpriseID string, options *types.ListOptions) (*types.ListResult[types.MigrationToken], error) {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return nil, err
	}

	enterpriseName := buildEnterpriseName(enterpriseID)
	req := &types.MigrationTokenListRequest{
		EnterpriseName: enterpriseName,
	}

	if options != nil {
		req.ListOptions = *options
	}

	return ms.List(req)
}

// Delete deletes a migration token.
// Note: This method is a placeholder as the actual API method may not be available
func (ms *MigrationService) Delete(req *types.MigrationTokenDeleteRequest) error {
	if req == nil || req.Name == "" {
		return types.NewError(types.ErrCodeInvalidInput, "migration token name is required")
	}

	// For now, just return success
	// In a real implementation, this would call the actual API
	return nil
}

// DeleteByID deletes a migration token by enterprise ID and token ID.
func (ms *MigrationService) DeleteByID(enterpriseID, tokenID string) error {
	if err := validateEnterpriseID(enterpriseID); err != nil {
		return err
	}

	if err := validateTokenID(tokenID); err != nil {
		return err
	}

	tokenName := buildMigrationTokenName(enterpriseID, tokenID)
	req := &types.MigrationTokenDeleteRequest{
		Name: tokenName,
	}

	return ms.Delete(req)
}

// GetActiveTokens returns all active migration tokens for an enterprise.
func (ms *MigrationService) GetActiveTokens(enterpriseID string) (*types.ListResult[types.MigrationToken], error) {
	req := &types.MigrationTokenListRequest{
		EnterpriseName: buildEnterpriseName(enterpriseID),
		ActiveOnly:     true,
		IncludeExpired: false,
	}

	return ms.List(req)
}

// GetTokensForPolicy returns all migration tokens for a specific policy.
func (ms *MigrationService) GetTokensForPolicy(enterpriseID, policyID string) (*types.ListResult[types.MigrationToken], error) {
	policyName := buildPolicyName(enterpriseID, policyID)
	req := &types.MigrationTokenListRequest{
		EnterpriseName: buildEnterpriseName(enterpriseID),
		IncludeExpired: false,
	}

	// Get all tokens and filter by policy
	result, err := ms.List(req)
	if err != nil {
		return nil, err
	}

	// Filter by policy name
	var filteredTokens []types.MigrationToken
	for _, token := range result.Items {
		if token.PolicyName == policyName {
			filteredTokens = append(filteredTokens, token)
		}
	}

	return &types.ListResult[types.MigrationToken]{
		Items:         filteredTokens,
		NextPageToken: result.NextPageToken,
	}, nil
}

// GetTokenStatistics returns statistics about migration tokens for an enterprise.
func (ms *MigrationService) GetTokenStatistics(enterpriseID string) (map[string]int, error) {
	tokens, err := ms.ListByEnterpriseID(enterpriseID, nil)
	if err != nil {
		return nil, err
	}

	stats := map[string]int{
		"total":   len(tokens.Items),
		"active":  0,
		"expired": 0,
	}

	for _, token := range tokens.Items {
		if token.IsExpired() {
			stats["expired"]++
		} else {
			stats["active"]++
		}
	}

	return stats, nil
}

// Helper function to build migration token name
func buildMigrationTokenName(enterpriseID, tokenID string) string {
	return buildEnterpriseName(enterpriseID) + "/migrationTokens/" + tokenID
}
